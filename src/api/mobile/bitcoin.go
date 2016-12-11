package mobile

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"strings"

	"github.com/skycoin/skycoin-exchange/src/coin"
	bitcoin "github.com/skycoin/skycoin-exchange/src/coin/bitcoin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/sknet"
	"github.com/skycoin/skycoin-exchange/src/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

type btcNode struct {
	NodeAddr string
	fee      string // bitcoin fee
}

type btcSendParams struct {
	WalletID string
	ToAddr   string
	Amount   uint64
	Fee      uint64
}

func newBitcoin(nodeAddr string) *btcNode {
	return &btcNode{NodeAddr: nodeAddr, fee: "2000"} // default transaction fee is 2000
}

func (bn btcNode) GetNodeAddr() string {
	return bn.NodeAddr
}

func (bn btcNode) ValidateAddr(address string) error {
	_, err := cipher.BitcoinDecodeBase58Address(address)
	return err
}

func (bn btcNode) GetBalance(addrs []string) (uint64, error) {
	req := pp.GetUtxoReq{
		CoinType:  pp.PtrString("bitcoin"),
		Addresses: addrs,
	}
	res := pp.GetUtxoRes{}
	if err := sknet.EncryGet(bn.NodeAddr, "/get/utxos", req, &res); err != nil {
		return 0, err
	}
	var bal uint64
	for _, u := range res.BtcUtxos {
		bal += u.GetAmount()
	}

	return bal, nil
}

func (bn btcNode) CreateRawTx(txIns []coin.TxIn, getKey coin.GetPrivKey, txOuts interface{}) (string, error) {
	coin := bitcoin.Bitcoin{}
	rawtx, err := coin.CreateRawTx(txIns, txOuts)
	if err != nil {
		return "", fmt.Errorf("create raw tx failed:%v", err)
	}

	return coin.SignRawTx(rawtx, getKey)
}

func (bn btcNode) BroadcastTx(rawtx string) (string, error) {
	req := pp.InjectTxnReq{
		CoinType: pp.PtrString("bitcoin"),
		Tx:       pp.PtrString(rawtx),
	}
	res := pp.InjectTxnRes{}
	if err := sknet.EncryGet(bn.NodeAddr, "/inject/tx", req, &res); err != nil {
		return "", err
	}

	if !res.Result.GetSuccess() {
		return "", fmt.Errorf("broadcast tx failed: %v", res.Result.GetReason())
	}

	return res.GetTxid(), nil
}

func (bn btcNode) GetTransactionByID(txid string) (string, error) {
	req := pp.GetTxReq{
		CoinType: pp.PtrString("bitcoin"),
		Txid:     pp.PtrString(txid),
	}
	res := pp.GetTxRes{}
	if err := sknet.EncryGet(bn.NodeAddr, "/get/tx", req, &res); err != nil {
		return "", err
	}

	if !res.Result.GetSuccess() {
		return "", fmt.Errorf("get bitcoin transaction by id failed: %v", res.Result.GetReason())
	}

	d, err := json.Marshal(res.GetTx())
	if err != nil {
		return "", err
	}
	return string(d), nil
}

// Fee option for setting transaction fee.
func Fee(n string) Option {
	return func(v interface{}) {
		btc := v.(*btcNode)
		btc.fee = n
	}
}

// Send amount bitcoins to address from specific wallet, Note: should not be used concurrently.
func (bc *btcNode) Send(walletID, toAddr, amount string, ops ...Option) (string, error) {
	for _, op := range ops {
		op(bc)
	}

	// validate amount
	amt, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return "", fmt.Errorf("parse amount string to uint64 failed: %v", err)
	}

	// validate fee
	fe, err := strconv.ParseUint(bc.fee, 10, 64)
	if err != nil {
		return "", fmt.Errorf("parse fee string to uint64 failed: %v", err)
	}

	if fe < 1000 {
		return "", fmt.Errorf("insufficient fee")
	}

	params := btcSendParams{WalletID: walletID, ToAddr: toAddr, Amount: amt, Fee: fe}

	txIns, txOut, err := bc.PrepareTx(params)
	if err != nil {
		return "", err
	}

	rawtx, err := bc.CreateRawTx(txIns, getPrivateKey(walletID), txOut)
	if err != nil {
		return "", fmt.Errorf("create raw transaction failed:%v", err)
	}

	txid, err := bc.BroadcastTx(rawtx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`{"txid":"%s"}`, txid), nil
}

func (bn btcNode) PrepareTx(params interface{}) ([]coin.TxIn, interface{}, error) {
	p := params.(btcSendParams)

	tp := strings.Split(p.WalletID, "_")[0]
	if tp != "bitcoin" {
		return nil, nil, fmt.Errorf("invalid wallet %v", tp)
	}

	// valid address
	if err := bn.ValidateAddr(p.ToAddr); err != nil {
		return nil, nil, err
	}

	addrs, err := wallet.GetAddresses(p.WalletID)
	if err != nil {
		return nil, nil, err
	}

	totalUtxos, err := bn.getOutputs(addrs)
	if err != nil {
		return nil, nil, err
	}

	utxos, bal, err := bn.getSufficientOutputs(totalUtxos, p.Amount+p.Fee)
	if err != nil {
		return nil, nil, err
	}

	txIns := make([]coin.TxIn, len(utxos))
	for i, u := range utxos {
		txIns[i] = coin.TxIn{
			Txid:    u.GetTxid(),
			Vout:    u.GetVout(),
			Address: u.GetAddress(),
		}
	}

	var txOut []bitcoin.TxOut
	chgAmt := bal - p.Amount - p.Fee
	chgAddr := addrs[0]
	if chgAmt > 0 {
		txOut = append(txOut,
			bn.makeTxOut(p.ToAddr, p.Amount),
			bn.makeTxOut(chgAddr, chgAmt))
	} else {
		txOut = append(txOut, bn.makeTxOut(p.ToAddr, p.Amount))
	}

	return txIns, txOut, nil
}

func (bn btcNode) makeTxOut(addr string, value uint64) bitcoin.TxOut {
	return bitcoin.TxOut{
		Addr:  addr,
		Value: value,
	}
}

func (bn btcNode) getSufficientOutputs(utxos []*pp.BtcUtxo, amt uint64) ([]*pp.BtcUtxo, uint64, error) {
	outMap := make(map[string][]*pp.BtcUtxo)
	for _, u := range utxos {
		outMap[u.GetAddress()] = append(outMap[u.GetAddress()], u)
	}

	allUtxos := []*pp.BtcUtxo{}
	var allBal uint64
	for _, utxos := range outMap {
		allBal += func(utxos []*pp.BtcUtxo) uint64 {
			var bal uint64
			for _, u := range utxos {
				if u.GetAmount() == 0 {
					continue
				}
				bal += u.GetAmount()
			}
			return bal
		}(utxos)

		allUtxos = append(allUtxos, utxos...)
		if allBal >= amt {
			return allUtxos, allBal, nil
		}
	}
	return nil, 0, errors.New("insufficient balance")
}

func (bn btcNode) getOutputs(addrs []string) ([]*pp.BtcUtxo, error) {
	req := pp.GetUtxoReq{
		CoinType:  pp.PtrString("bitcoin"),
		Addresses: addrs,
	}
	res := pp.GetUtxoRes{}
	if err := sknet.EncryGet(bn.NodeAddr, "/get/utxos", req, &res); err != nil {
		return nil, err
	}

	if !res.Result.GetSuccess() {
		return nil, fmt.Errorf("get utxos failed: %v", res.Result.GetReason())
	}

	return res.BtcUtxos, nil
}
