package server

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/skycoin/skycoin-exchange/src/server/account"
	bitcoin "github.com/skycoin/skycoin-exchange/src/server/coin_interface/bitcoin"
	skycoin "github.com/skycoin/skycoin-exchange/src/server/coin_interface/skycoin"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
)

// Config store server's configuration.
type Config struct {
	Port         int           // api port
	BtcFee       int           // btc transaction fee
	DataDir      string        // data directory
	WalletName   string        // wallet name
	AcntName     string        // accounts file name
	Seed         string        // seed
	Seckey       cipher.SecKey // private key
	UtxoPoolSize int           // utxo pool size.
}

/*
	The server gets events from the client and processes them
	- get balance/status
	- get deposit addresses
	- withdrawl bitcoin
	- withdrawl skycoin
	- add bid
	- add ask
	- get order book
*/

type ExchangeServer struct {
	account.AccountManager
	btcum  bitcoin.UtxoManager
	skyum  skycoin.UtxoManager
	cfg    Config
	wallet wallet.Wallet
	wltMtx sync.RWMutex // mutex for protecting the wallet.
}

// New create new server
func New(cfg Config) engine.Exchange {
	// init the data dir
	path := initDataDir(cfg.DataDir)

	// init the wallet dir.
	wallet.InitDir(filepath.Join(path, "wallets"))

	// init the account dir
	account.InitDir(filepath.Join(path, "account/server"))

	// load account manager if exist.
	var (
		acntMgr account.AccountManager
		err     error
	)

	acntMgr, err = account.LoadAccountManager(cfg.AcntName)
	if err != nil {
		glog.Error(err)
		if os.IsNotExist(err) {
			acntMgr = account.NewAccountManager(cfg.AcntName)
		} else {
			panic(err)
		}
	}

	// get wallet
	var wlt wallet.Wallet
	wlt, err = wallet.Load(cfg.WalletName)
	if err != nil {
		if os.IsNotExist(err) {
			glog.Info("wallet file not exist")
			wlt, err = wallet.New(cfg.WalletName, wallet.Deterministic, cfg.Seed)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	s := &ExchangeServer{
		cfg:            cfg,
		wallet:         wlt,
		AccountManager: acntMgr,
		btcum:          bitcoin.NewUtxoManager(cfg.UtxoPoolSize, wlt.GetAddresses(wallet.Bitcoin)),
		skyum:          skycoin.NewUtxoManager(cfg.UtxoPoolSize, wlt.GetAddresses(wallet.Skycoin)),
	}
	return s
}

func (self *ExchangeServer) Run() {
	log.Println("skycoin-exchange server started, port:", self.cfg.Port)

	// start the utxo manager
	c := make(chan bool)
	go self.btcum.Start(c)
	go self.skyum.Start(c)

	// start the api server.
	r := NewRouter(self)
	r.Run(fmt.Sprintf(":%d", self.cfg.Port))
}

func (self ExchangeServer) GetBtcFee() uint64 {
	return uint64(self.cfg.BtcFee)
}

// GetServPrivKey returnt he sever's private key.
func (self ExchangeServer) GetServPrivKey() cipher.SecKey {
	return self.cfg.Seckey
}

// GetPrivKey return the private key of specific address.
func (self ExchangeServer) GetPrivKey(ct wallet.CoinType, addr string) (string, error) {
	entry, err := self.wallet.GetAddressEntry(ct, addr)
	if err != nil {
		return "", err
	}

	return entry.Secret, nil
}

func (self *ExchangeServer) GetNewAddress(coinType wallet.CoinType) string {
	self.wltMtx.Lock()
	defer self.wltMtx.Unlock()
	addrEntry, err := self.wallet.NewAddresses(coinType, 1)
	if err != nil {
		panic("server get new address failed")
	}
	return addrEntry[0].Address
}

// BtcChooseUtxos choose appropriate bitcoin utxos,
func (self *ExchangeServer) ChooseUtxos(ct wallet.CoinType, amount uint64, tm time.Duration) (interface{}, error) {
	switch ct {
	case wallet.Bitcoin:
		return self.btcum.ChooseUtxos(amount, tm)
	case wallet.Skycoin:
		return self.skyum.ChooseUtxos(amount, tm)
	default:
		return nil, errors.New("unknow coin type")
	}
}

func (self *ExchangeServer) PutUtxos(ct wallet.CoinType, utxos interface{}) {
	switch ct {
	case wallet.Bitcoin:
		btcUtxos := utxos.([]bitcoin.Utxo)
		for _, u := range btcUtxos {
			self.btcum.PutUtxo(u)
		}
	case wallet.Skycoin:
		skyUtxos := utxos.([]skycoin.Utxo)
		for _, u := range skyUtxos {
			self.skyum.PutUtxo(u)
		}
	}
}

// AddWatchAddress add watch address for utxo manager.
func (self *ExchangeServer) WatchAddress(ct wallet.CoinType, addr string) {
	switch ct {
	case wallet.Bitcoin:
		self.btcum.WatchAddresses([]string{addr})
	}
}

func (self *ExchangeServer) BtcPutUtxos(utxos []bitcoin.Utxo) {
	for _, u := range utxos {
		self.btcum.PutUtxo(u)
	}
}

func (self *ExchangeServer) SaveAccount() error {
	return self.Save()
}

func initDataDir(dir string) string {
	//DataDir = dir
	if dir == "" {
		glog.Error("data directory is nil")
	}

	home := util.UserHome()
	if home == "" {
		glog.Warning("Failed to get home directory")
		dir = filepath.Join("./", dir)
	} else {
		dir = filepath.Join(home, dir)
	}

	if err := os.MkdirAll(dir, os.FileMode(0700)); err != nil {
		glog.Error("Failed to create directory %s: %v", dir, err)
	}
	return dir
}
