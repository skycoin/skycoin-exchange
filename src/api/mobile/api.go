package mobile

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/wallet"
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

// func SendSky(walletID string, toAddr string, amount string) (string, error) {
func SendSky(walletID string, toAddr string, amount string) (string, error) {
	// validate amount
	amt, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return "", fmt.Errorf("parse amount string to uint64 failed: %v", err)
	}

	params := skySendParams{WalletID: walletID, ToAddr: toAddr, Amount: amt}
	node, ok := nodeMap["skycoin"]
	if !ok {
		return "", errors.New("skycoin is not supported")
	}

	txIns, txOut, err := node.PrepareTx(params)
	if err != nil {
		return "", err
	}

	// prepare keys
	rawtx, err := node.CreateRawTx(txIns, getPrivateKey(walletID), txOut)
	if err != nil {
		return "", fmt.Errorf("create raw transaction failed:%v", err)
	}

	return node.BroadcastTx(rawtx)
}

func SendBtc(walletID string, toAddr string, amount string, fee string) (string, error) {
	// validate amount
	amt, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return "", fmt.Errorf("parse amount string to uint64 failed: %v", err)
	}

	// validate fee
	fe, err := strconv.ParseUint(fee, 10, 64)
	if err != nil {
		return "", fmt.Errorf("parse fee string to uint64 failed: %v", err)
	}

	if fe < 1000 {
		return "", fmt.Errorf("insufficient fee")
	}

	params := btcSendParams{WalletID: walletID, ToAddr: toAddr, Amount: amt, Fee: fe}
	node, ok := nodeMap["bitcoin"]
	if !ok {
		return "", errors.New("bitcoin is not supported")
	}

	txIns, txOut, err := node.PrepareTx(params)
	if err != nil {
		return "", err
	}

	rawtx, err := node.CreateRawTx(txIns, getPrivateKey(walletID), txOut)
	if err != nil {
		return "", fmt.Errorf("create raw transaction failed:%v", err)
	}

	return node.BroadcastTx(rawtx)
}

func getPrivateKey(walletID string) coin.GetPrivKey {
	return func(addr string) (string, error) {
		_, s, err := wallet.GetKeypair(walletID, addr)
		return s, err
	}
}
