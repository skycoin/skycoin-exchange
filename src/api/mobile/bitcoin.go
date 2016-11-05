package mobile

import (
	"errors"
	"fmt"

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
}

type btcSendParams struct {
	WalletID string
	ToAddr   string
	Amount   uint64
	Fee      uint64
}

func (bn btcNode) ValidateAddr(address string) error {
	_, err := cipher.BitcoinDecodeBase58Address(address)
	return err
}

func (bn btcNode) GetBalance(addrs []string) (uint64, error) {
	// get uxout of the address
	_, s := cipher.GenerateKeyPair()
	sknet.SetKey(s.Hex())

	req := pp.GetUtxoReq{
		CoinType:  pp.PtrString("bitcoin"),
		Addresses: addrs,
	}
	res := pp.GetUtxoRes{}
	if err := sknet.EncryGet(bn.NodeAddr, "/auth/get/utxos", req, &res); err != nil {
		return 0, err
	}
	var bal uint64
	for _, u := range res.BtcUtxos {
		bal += u.GetAmount()
	}

	return bal, nil
}

func (bn btcNode) CreateRawTx(txIns []coin.TxIn, getKey coin.GetPrivKey, txOuts interface{}) (string, error) {
	gw := bitcoin.Gateway{}
	rawtx, err := gw.CreateRawTx(txIns, txOuts)
	if err != nil {
		return "", fmt.Errorf("create raw tx failed:%v", err)
	}

	return gw.SignRawTx(rawtx, getKey)
}

func (bn btcNode) BroadcastTx(rawtx string) (string, error) {
	gw := bitcoin.Gateway{}
	return gw.InjectTx(rawtx)
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
	// get uxout of the address
	_, s := cipher.GenerateKeyPair()
	sknet.SetKey(s.Hex())

	req := pp.GetUtxoReq{
		CoinType:  pp.PtrString("bitcoin"),
		Addresses: addrs,
	}
	res := pp.GetUtxoRes{}
	if err := sknet.EncryGet(bn.NodeAddr, "/auth/get/utxos", req, &res); err != nil {
		return nil, err
	}

	if !res.Result.GetSuccess() {
		return nil, fmt.Errorf("get utxos failed: %v", res.Result.GetReason())
	}

	return res.BtcUtxos, nil
}
