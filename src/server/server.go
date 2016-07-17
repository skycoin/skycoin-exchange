package server

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"time"

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
	um     UtxoManager
	cfg    Config
	wallet wallet.Wallet
	wltMtx sync.RWMutex // mutex for protecting the wallet.
}

// New create new server
func New(cfg Config) engine.Exchange {
	// init the data dir
	path := util.InitDataDir(cfg.DataDir)

	// set the wallet dir.
	wallet.WltDir = filepath.Join(path, "wallets")

	// get wallet
	var wlt wallet.Wallet
	var err error
	if wallet.IsExist(cfg.WalletName) {
		wlt, err = wallet.Load(cfg.WalletName)
		if err != nil {
			panic("server load walle failed")
		}
	} else {
		wlt, err = wallet.New(cfg.WalletName, wallet.Deterministic, cfg.Seed)
		if err != nil {
			panic("server create wallet failed")
		}
	}
	s := &ExchangeServer{
		cfg:            cfg,
		wallet:         wlt,
		AccountManager: account.NewExchangeAccountManager(),
		um: &ExUtxoManager{
			UtxosCh: map[wallet.CoinType]chan bitcoin.Utxo{
				wallet.Bitcoin: make(chan bitcoin.Utxo),
				wallet.Skycoin: make(chan bitcoin.Utxo),
			},
			UtxoStateMap: map[wallet.CoinType]map[string]bitcoin.Utxo{
				wallet.Bitcoin: make(map[string]bitcoin.Utxo),
				wallet.Skycoin: make(map[string]bitcoin.Utxo)},
		},
	}
	return s
}

func (self *ExchangeServer) Run() {
	log.Println("skycoin-exchange server started, port:", self.cfg.Port)

	// start the utxo manager
	c := make(chan bool)
	go func() { self.um.Start(c) }()

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

// ChooseUtxos choose appropriate utxos, if time out, and not found enough utxos,
// the utxos got before will put back to the utxos pool, and return error.
// the tm is millisecond
func (self *ExchangeServer) ChooseUtxos(ct wallet.CoinType, amount uint64, tm time.Duration) ([]bitcoin.UtxoWithkey, error) {
	var totalAmount uint64
	utxos := []bitcoin.UtxoWithkey{}
	for {
		select {
		case utxo := <-self.um.GetUtxo(ct):
			// get private key
			key, err := self.GetPrivKey(ct, utxo.GetAddress())
			if err != nil {
				self.um.PutUtxo(ct, utxo)
				return []bitcoin.UtxoWithkey{}, err
			}

			utxok := bitcoin.NewUtxoWithKey(utxo, key)
			utxos = append(utxos, utxok)

			totalAmount += utxo.GetAmount()
			if totalAmount >= (amount + self.GetFee()) {
				return utxos, nil
			}

		case <-time.After(tm):
			// put utxos back
			for _, u := range utxos {
				self.um.PutUtxo(ct, u)
			}
			return []bitcoin.UtxoWithkey{}, errors.New("choose utxo time out")
		}
	}
}

// AddWatchAddress add watch address for utxo manager.
func (self *ExchangeServer) AddWatchAddress(ct wallet.CoinType, addr string) {
	self.um.AddWatchAddress(ct, addr)
}

func (self *ExchangeServer) PutUtxos(ct wallet.CoinType, utxos []bitcoin.UtxoWithkey) {
	for _, u := range utxos {
		self.um.PutUtxo(ct, u)
	}
}
