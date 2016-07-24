package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/skycoin/skycoin-exchange/src/server/account"
	bitcoin "github.com/skycoin/skycoin-exchange/src/server/coin_interface/bitcoin"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
)

var CheckTick = 5 * time.Second

// Config store server's configuration.
type Config struct {
	Port         int           // api port
	Fee          int           // transaction fee
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
	cfg    Config
	wallet wallet.Wallet
	wltMtx sync.RWMutex // mutex for protecting the wallet.
}

// New create new server
func New(cfg Config) engine.Exchange {
	// init the data dir
	path := util.InitDataDir(cfg.DataDir)

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

	btcum := bitcoin.NewUtxoManager(wlt, cfg.UtxoPoolSize)
	s := &ExchangeServer{
		cfg:            cfg,
		wallet:         wlt,
		AccountManager: acntMgr,
		btcum:          btcum,
	}
	return s
}

func (self *ExchangeServer) Run() {
	log.Println("skycoin-exchange server started, port:", self.cfg.Port)

	// start the utxo manager
	c := make(chan bool)
	go func() { self.btcum.Start(c) }()

	// start the api server.
	r := NewRouter(self)
	r.Run(fmt.Sprintf(":%d", self.cfg.Port))
}

func (self ExchangeServer) GetFee() uint64 {
	return uint64(self.cfg.Fee)
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
func (self *ExchangeServer) BtcChooseUtxos(amount uint64) ([]bitcoin.Utxo, error) {
	return self.btcum.ChooseUtxos(amount)
}

// AddWatchAddress add watch address for utxo manager.
func (self *ExchangeServer) AddWatchAddress(ct wallet.CoinType, addr string) {
	switch ct {
	case wallet.Bitcoin:
		self.btcum.AddWatchAddress(addr)
	}
}

func (self *ExchangeServer) BtcPutUtxos(ct wallet.CoinType, utxos []bitcoin.Utxo) {
	for _, u := range utxos {
		self.btcum.PutUtxo(ct, u)
	}
}

func (self *ExchangeServer) SaveAccount() error {
	return self.Save()
}
