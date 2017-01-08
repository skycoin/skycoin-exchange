package mobile

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"strings"

	"github.com/skycoin/skycoin-exchange/src/coin"
	skycoin "github.com/skycoin/skycoin-exchange/src/coin/skycoin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/sknet"
	"github.com/skycoin/skycoin-exchange/src/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

// Option used as option argument in coin.Send method.
type Option func(c interface{})

// Coiner coin client interface
type Coiner interface {
	Name() string
	GetBalance(addrs []string) (uint64, error)
	ValidateAddr(addr string) error
	PrepareTx(params interface{}) ([]coin.TxIn, interface{}, error)
	CreateRawTx(txIns []coin.TxIn, getKey coin.GetPrivKey, txOuts interface{}) (string, error)
	BroadcastTx(rawtx string) (string, error)
	GetTransactionByID(txid string) (string, error)
	GetOutputByID(outid string) (string, error)
	GetNodeAddr() string
	Send(walletID string, toAddr string, amount string, ops ...Option) (string, error)
}

// CoinEx implements the Coin interface.
type coinEx struct {
	name     string
	nodeAddr string
}

type sendParams struct {
	WalletID string
	ToAddr   string
	Amount   uint64
}

func newCoin(name, nodeAddr string) *coinEx {
	return &coinEx{name: name, nodeAddr: nodeAddr}
}

func (cn coinEx) Name() string {
	return cn.name
}

// GetNodeAddr returns the coin's node address
func (cn coinEx) GetNodeAddr() string {
	return cn.nodeAddr
}

func (cn coinEx) getOutputs(addrs []string) ([]*pp.SkyUtxo, error) {
	req := pp.GetUtxoReq{
		CoinType:  pp.PtrString(cn.name),
		Addresses: addrs,
	}
	res := pp.GetUtxoRes{}
	if err := sknet.EncryGet(cn.nodeAddr, "/get/utxos", req, &res); err != nil {
		return nil, err
	}

	if !res.Result.GetSuccess() {
		return nil, fmt.Errorf("get utxos failed: %v", res.Result.GetReason())
	}

	return res.SkyUtxos, nil
}

// GetBalance gets balance of specific addresses
func (cn coinEx) GetBalance(addrs []string) (uint64, error) {
	utxos, err := cn.getOutputs(addrs)
	if err != nil {
		return 0, err
	}

	var bal uint64
	for _, u := range utxos {
		bal += u.GetCoins()
	}

	return bal, nil
}

// ValidateAddr check if the address is validated
func (cn coinEx) ValidateAddr(address string) error {
	_, err := cipher.DecodeBase58Address(address)
	return err
}

// CreateRawTx creates raw transaction
func (cn coinEx) CreateRawTx(txIns []coin.TxIn, getKey coin.GetPrivKey, txOuts interface{}) (string, error) {
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
			return "", fmt.Errorf("%s coins must be multiple of 1e6", cn.Name())
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

// BroadcastTx injects transaction
func (cn coinEx) BroadcastTx(rawtx string) (string, error) {
	req := pp.InjectTxnReq{
		CoinType: pp.PtrString(cn.name),
		Tx:       pp.PtrString(rawtx),
	}
	res := pp.InjectTxnRes{}
	if err := sknet.EncryGet(cn.nodeAddr, "/inject/tx", req, &res); err != nil {
		return "", err
	}

	if !res.Result.GetSuccess() {
		return "", fmt.Errorf("broadcast tx failed: %v", res.Result.GetReason())
	}

	return res.GetTxid(), nil
}

// GetTransactionByID gets transaction verbose info by id
func (cn coinEx) GetTransactionByID(txid string) (string, error) {
	req := pp.GetTxReq{
		CoinType: pp.PtrString(cn.name),
		Txid:     pp.PtrString(txid),
	}
	res := pp.GetTxRes{}
	if err := sknet.EncryGet(cn.nodeAddr, "/get/tx", req, &res); err != nil {
		return "", err
	}

	if !res.Result.GetSuccess() {
		return "", fmt.Errorf("get %s transaction by id failed: %v", cn.Name(), res.Result.GetReason())
	}
	d, err := json.Marshal(res.GetTx())
	if err != nil {
		return "", err
	}
	return string(d), nil
}

// PrepareTx prepares the transaction info
func (cn coinEx) PrepareTx(params interface{}) ([]coin.TxIn, interface{}, error) {
	p := params.(sendParams)

	tp := strings.Split(p.WalletID, "_")[0]
	if tp != cn.name {
		return nil, nil, fmt.Errorf("invalid wallet %v", tp)
	}

	// validate address
	if err := cn.ValidateAddr(p.ToAddr); err != nil {
		return nil, nil, err
	}

	addrs, err := wallet.GetAddresses(p.WalletID)
	if err != nil {
		return nil, nil, err
	}

	// outMap := make(map[string][]*pp.SkyUtxo)
	totalUtxos, err := cn.getOutputs(addrs)
	if err != nil {
		return nil, nil, err
	}

	utxos, err := cn.getSufficientOutputs(totalUtxos, p.Amount)
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
			cn.makeTxOut(p.ToAddr, p.Amount, chgHours/2),
			cn.makeTxOut(chgAddr, chgAmt, chgHours/2))
	} else {
		txOut = append(txOut, cn.makeTxOut(p.ToAddr, p.Amount, chgHours/2))
	}
	return txIns, txOut, nil
}

// Send sends numbers of coins to toAddr from specific wallet
func (cn *coinEx) Send(walletID, toAddr, amount string, ops ...Option) (string, error) {
	for _, op := range ops {
		op(cn)
	}

	// validate amount
	amt, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return "", fmt.Errorf("parse amount string to uint64 failed: %v", err)
	}

	params := sendParams{WalletID: walletID, ToAddr: toAddr, Amount: amt}

	txIns, txOut, err := cn.PrepareTx(params)
	if err != nil {
		return "", err
	}

	// prepare keys
	rawtx, err := cn.CreateRawTx(txIns, getPrivateKey(walletID), txOut)
	if err != nil {
		return "", fmt.Errorf("create raw transaction failed:%v", err)
	}

	txid, err := cn.BroadcastTx(rawtx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`{"txid":"%s"}`, txid), nil
}

func (cn coinEx) GetOutputByID(outid string) (string, error) {
	req := pp.GetOutputReq{
		CoinType: pp.PtrString(cn.Name()),
		Hash:     pp.PtrString(outid),
	}

	res := pp.GetOutputRes{}
	if err := sknet.EncryGet(cn.GetNodeAddr(), "/get/output", req, &res); err != nil {
		return "", err
	}

	if !res.Result.GetSuccess() {
		return "", fmt.Errorf("get output failed: %v", res.Result.GetReason())
	}

	d, err := json.Marshal(res.GetOutput())
	if err != nil {
		return "", fmt.Errorf("unmarshal result failed, %v", err)
	}

	return string(d), nil
}

func (cn coinEx) getSufficientOutputs(utxos []*pp.SkyUtxo, amt uint64) ([]*pp.SkyUtxo, error) {
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

func (cn coinEx) makeTxOut(addr string, coins uint64, hours uint64) skycoin.TxOut {
	out := skycoin.TxOut{}
	out.Address = cipher.MustDecodeBase58Address(addr)
	out.Coins = coins
	out.Hours = hours
	return out
}
