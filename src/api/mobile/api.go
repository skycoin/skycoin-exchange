package mobile

import (
	"encoding/json"
	"fmt"

	"strings"

	"strconv"

	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

var config Config
var nodeMap map[string]noder

// Config used for init the api env, includes wallet dir path, skycoin node and bitcoin node address.
// the node address is consisted of ip and port, eg: 127.0.0.1:6420
type Config struct {
	WalletDirPath string `json:"wallet_dir_path"`
	ServerAddr    string `json:"server_addr"`
}

// NewConfig create config instance.
func NewConfig() *Config {
	return &Config{}
}

// Init initialize wallet dir and node instance.
func Init(cfg *Config) {
	wallet.InitDir(cfg.WalletDirPath)
	config = *cfg
	nodeMap = map[string]noder{
		"skycoin": &skyNode{NodeAddr: config.ServerAddr},
		"bitcoin": &btcNode{NodeAddr: config.ServerAddr},
	}
}

// NewWallet create a new wallet base on the wallet type and seed
func NewWallet(coinType string, seed string) (string, error) {
	tp, err := coin.TypeFromStr(coinType)
	if err != nil {
		return "", err
	}
	wlt, err := wallet.New(tp, seed)
	if err != nil {
		return "", err
	}
	return wlt.GetID(), nil
}

// NewAddress generate address in specific wallet.
func NewAddress(walletID string, num int) (string, error) {
	es, err := wallet.NewAddresses(walletID, num)
	if err != nil {
		return "", err
	}
	var res = struct {
		Entries []coin.AddressEntry `json:"addresses"`
	}{
		es,
	}
	d, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(d), nil
}

// GetAddresses return all addresses in the wallet.
func GetAddresses(walletID string) (string, error) {
	addrs, err := wallet.GetAddresses(walletID)
	if err != nil {
		return "", err
	}
	var res = struct {
		Addresses []string `json:"addresses"`
	}{
		addrs,
	}

	d, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(d), nil
}

// GetKeyPairOfAddr get pubkey and seckey pair of address in specific wallet.
func GetKeyPairOfAddr(walletID string, addr string) (string, error) {
	p, s, err := wallet.GetKeypair(walletID, addr)
	if err != nil {
		return "", err
	}
	var res = struct {
		Pubkey string `json:"pubkey"`
		Seckey string `json:"seckey"`
	}{
		p,
		s,
	}

	d, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

// GetBalance return balance of specific address.
func GetBalance(coinType string, address string) (string, error) {
	node, ok := nodeMap[coinType]
	if !ok {
		return "", fmt.Errorf("%s coin does not support", coinType)
	}

	if err := node.ValidateAddr(address); err != nil {
		return "", err
	}

	bal, err := node.GetBalance([]string{address})
	if err != nil {
		return "", err
	}

	var res = struct {
		Balance uint64 `json:"balance"`
	}{
		bal,
	}

	d, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func Send(walletID string, toAddr string, amount string) (string, error) {
	tp := strings.Split(walletID, "_")[0]
	node, ok := nodeMap[tp]
	if !ok {
		return "", fmt.Errorf("%s coin does not support", tp)
	}

	// validate address
	if err := node.ValidateAddr(toAddr); err != nil {
		return "", err
	}

	// validate amount
	amt, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return "", fmt.Errorf("parse amount string to uint64 failed: %v", err)
	}

	addrs, err := wallet.GetAddresses(walletID)
	if err != nil {
		return "", err
	}

	txIns, inAddrs, txOut, err := node.PrepareTx(addrs, toAddr, amt)
	if err != nil {
		return "", err
	}

	// prepare keys
	keys := make([]cipher.SecKey, len(inAddrs))
	for i, a := range inAddrs {
		_, s, err := wallet.GetKeypair(walletID, a)
		if err != nil {
			return "", fmt.Errorf("get key failed:%v", err)
		}

		k, err := cipher.SecKeyFromHex(s)
		if err != nil {
			return "", fmt.Errorf("error private key:%v", err)
		}
		keys[i] = k
	}

	rawtx, err := node.CreateRawTx(txIns, keys, txOut)
	if err != nil {
		return "", fmt.Errorf("create raw transaction failed:%v", err)
	}

	return node.BroadcastTx(rawtx)
}
