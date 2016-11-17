package mobile

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"strings"

	"github.com/skycoin/skycoin-exchange/src/coin"
	skycoin "github.com/skycoin/skycoin-exchange/src/coin/skycoin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/sknet"
	"github.com/skycoin/skycoin-exchange/src/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

// var privateKey cipher.SecKey

type skyNode struct {
	NodeAddr string
}

type skySendParams struct {
	WalletID string
	ToAddr   string
	Amount   uint64
}

func (sn skyNode) getOutputs(addrs []string) ([]*pp.SkyUtxo, error) {
	req := pp.GetUtxoReq{
		CoinType:  pp.PtrString("skycoin"),
		Addresses: addrs,
	}
	res := pp.GetUtxoRes{}
	if err := sknet.EncryGet(sn.NodeAddr, "/auth/get/utxos", req, &res); err != nil {
		return nil, err
	}

	if !res.Result.GetSuccess() {
		return nil, fmt.Errorf("get utxos failed: %v", res.Result.GetReason())
	}

	return res.SkyUtxos, nil
}

func (sn skyNode) GetBalance(addrs []string) (uint64, error) {
	utxos, err := sn.getOutputs(addrs)
	if err != nil {
		return 0, err
	}

	var bal uint64
	for _, u := range utxos {
		bal += u.GetCoins()
	}

	return bal, nil
}

func (sn skyNode) ValidateAddr(address string) error {
	_, err := cipher.DecodeBase58Address(address)
	return err
}

func (sn skyNode) CreateRawTx(txIns []coin.TxIn, getKey coin.GetPrivKey, txOuts interface{}) (string, error) {
	tx := skycoin.Transaction{}
	for _, in := range txIns {
		tx.PushInput(cipher.MustSHA256FromHex(in.Txid))
	}

	s := reflect.ValueOf(txOuts)
	if s.Kind() != reflect.Slice {
		return "", errors.New("error tx out type")
	}
	outs := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		outs[i] = s.Index(i).Interface()
	}

	if len(outs) > 2 {
		return "", errors.New("out address more than 2")
	}

	for _, o := range outs {
		out := o.(skycoin.TxOut)
		if (out.Coins % 1e6) != 0 {
			return "", errors.New("skycoin coins must be multiple of 1e6")
		}
		tx.PushOutput(out.Address, out.Coins, out.Hours)
	}

	keys := make([]cipher.SecKey, len(txIns))
	for i, in := range txIns {
		s, err := getKey(in.Address)
		if err != nil {
			return "", fmt.Errorf("get private key failed:%v", err)
		}
		k, err := cipher.SecKeyFromHex(s)
		if err != nil {
			return "", fmt.Errorf("invalid private key:%v", err)
		}
		keys[i] = k
	}

	tx.SignInputs(keys)
	tx.UpdateHeader()

	d, err := tx.Serialize()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(d), nil
}

func (sn skyNode) BroadcastTx(rawtx string) (string, error) {
	req := pp.InjectTxnReq{
		CoinType: pp.PtrString("skycoin"),
		Tx:       pp.PtrString(rawtx),
	}
	res := pp.InjectTxnRes{}
	if err := sknet.EncryGet(sn.NodeAddr, "/auth/inject/tx", req, &res); err != nil {
		return "", err
	}

	if !res.Result.GetSuccess() {
		return "", fmt.Errorf("broadcast tx failed: %v", res.Result.GetReason())
	}

	return res.GetTxid(), nil
}

func (sn skyNode) GetTransactionByID(txid string) (string, error) {
	req := pp.GetTxReq{
		CoinType: pp.PtrString("skycoin"),
		Txid:     pp.PtrString(txid),
	}
	res := pp.GetTxRes{}
	if err := sknet.EncryGet(sn.NodeAddr, "/auth/get/tx", req, &res); err != nil {
		return "", err
	}

	if !res.Result.GetSuccess() {
		return "", fmt.Errorf("get skycoin transaction by id failed: %v", res.Result.GetReason())
	}
	d, err := json.Marshal(res.GetTx())
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func (sn skyNode) PrepareTx(params interface{}) ([]coin.TxIn, interface{}, error) {
	p := params.(skySendParams)

	tp := strings.Split(p.WalletID, "_")[0]
	if tp != "skycoin" {
		return nil, nil, fmt.Errorf("invalid wallet %v", tp)
	}

	// validate address
	if err := sn.ValidateAddr(p.ToAddr); err != nil {
		return nil, nil, err
	}

	addrs, err := wallet.GetAddresses(p.WalletID)
	if err != nil {
		return nil, nil, err
	}

	// outMap := make(map[string][]*pp.SkyUtxo)
	totalUtxos, err := sn.getOutputs(addrs)
	if err != nil {
		return nil, nil, err
	}

	utxos, err := sn.getSufficientOutputs(totalUtxos, p.Amount)
	if err != nil {
		return nil, nil, err
	}

	bal, hours := func(utxos []*pp.SkyUtxo) (uint64, uint64) {
		var c, h uint64
		for _, u := range utxos {
			c += u.GetCoins()
			h += u.GetHours()
		}
		return c, h
	}(utxos)

	txIns := make([]coin.TxIn, len(utxos))
	for i, u := range utxos {
		txIns[i] = coin.TxIn{
			Txid:    u.GetHash(),
			Address: u.GetAddress(),
		}
	}

	var txOut []skycoin.TxOut
	chgAmt := bal - p.Amount
	chgHours := hours / 4
	chgAddr := addrs[0]
	if chgAmt > 0 {
		txOut = append(txOut,
			sn.makeTxOut(p.ToAddr, p.Amount, chgHours/2),
			sn.makeTxOut(chgAddr, chgAmt, chgHours/2))
	} else {
		txOut = append(txOut, sn.makeTxOut(p.ToAddr, p.Amount, chgHours/2))
	}
	return txIns, txOut, nil
}

func (sn skyNode) getSufficientOutputs(utxos []*pp.SkyUtxo, amt uint64) ([]*pp.SkyUtxo, error) {
	outMap := make(map[string][]*pp.SkyUtxo)
	for _, u := range utxos {
		outMap[u.GetAddress()] = append(outMap[u.GetAddress()], u)
	}

	allUtxos := []*pp.SkyUtxo{}
	var allBal uint64
	for _, utxos := range outMap {
		allBal += func(utxos []*pp.SkyUtxo) uint64 {
			var bal uint64
			for _, u := range utxos {
				if u.GetCoins() == 0 {
					continue
				}
				bal += u.GetCoins()
			}
			return bal
		}(utxos)

		allUtxos = append(allUtxos, utxos...)
		if allBal >= amt {
			return allUtxos, nil
		}
	}

	return nil, errors.New("insufficient balance")
}

func (sn skyNode) makeTxOut(addr string, coins uint64, hours uint64) skycoin.TxOut {
	out := skycoin.TxOut{}
	out.Address = cipher.MustDecodeBase58Address(addr)
	out.Coins = coins
	out.Hours = hours
	return out
}
