package server

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gopkg.in/op/go-logging.v1"

	"github.com/skycoin/skycoin-exchange/src/server/account"
	"github.com/skycoin/skycoin-exchange/src/server/coin"
	bitcoin "github.com/skycoin/skycoin-exchange/src/server/coin/bitcoin"
	skycoin "github.com/skycoin/skycoin-exchange/src/server/coin/skycoin"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/order"
	"github.com/skycoin/skycoin-exchange/src/server/router"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
)

var logger = logging.MustGetLogger("exchange.server")

// Config store server's configuration.
type Config struct {
	Port         int           // api port
	BtcFee       int           // btc transaction fee
	DataDir      string        // data directory
	WalletName   string        // wallet name
	Seed         string        // seed
	Seckey       cipher.SecKey // private key
	UtxoPoolSize int           // utxo pool size.
}

type ExchangeServer struct {
	account.Manager
	btcum         bitcoin.UtxoManager
	skyum         skycoin.UtxoManager
	orderManager  *order.Manager
	cfg           Config
	wallet        wallet.Wallet
	wltMtx        sync.RWMutex                // mutex for protecting the wallet.
	orderHandlers map[string]chan order.Order // order handlers, for handleing bid and ask.
	coins         []string                    // supported coins
}

// New create new server
func New(cfg Config) engine.Exchange {
	// init the data dir
	path := initDataDir(cfg.DataDir)

	// init the wallet dir.
	wallet.InitDir(filepath.Join(path, "wallet"))

	// init the account dir
	account.InitDir(filepath.Join(path, "account"))

	// init the order book dir.
	order.InitDir(filepath.Join(path, "orderbook"))

	var (
		acntMgr account.Manager
		err     error
	)

	// load account manager if exist.
	acntMgr, err = account.LoadManager()
	if err != nil {
		if os.IsNotExist(err) {
			// create new account manager.
			acntMgr = account.NewManager()
		} else {
			panic(err)
		}
	}

	// get wallet
	var wlt wallet.Wallet
	wlt, err = wallet.Load(cfg.WalletName)
	if err != nil {
		if os.IsNotExist(err) {
			wlt, err = wallet.New(cfg.WalletName, wallet.Deterministic, cfg.Seed)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	// load or create order books.
	var orderManager *order.Manager
	orderManager, err = order.LoadManager()
	if err != nil {
		if os.IsNotExist(err) {
			orderManager = order.NewManager()
			orderManager.AddBook("bitcoin/skycoin", &order.Book{})
		} else {
			panic(err)
		}
	}

	s := &ExchangeServer{
		cfg:          cfg,
		wallet:       wlt,
		Manager:      acntMgr,
		btcum:        bitcoin.NewUtxoManager(cfg.UtxoPoolSize, wlt.GetAddresses(coin.Bitcoin)),
		skyum:        skycoin.NewUtxoManager(cfg.UtxoPoolSize, wlt.GetAddresses(coin.Skycoin)),
		orderManager: orderManager,
		coins:        []string{"BTC", "SKY"},
		orderHandlers: map[string]chan order.Order{
			"bitcoin/skycoin": make(chan order.Order, 100),
		},
	}

	return s
}

// Run start the exchange server.
func (self *ExchangeServer) Run() {
	logger.Info("server started, port:%d", self.cfg.Port)
	// register coins
	coin.RegisterGateway(coin.Bitcoin, &bitcoin.GatewayIns)
	coin.RegisterGateway(coin.Skycoin, &skycoin.GatewayIns)

	// register the order handlers
	for cp, c := range self.orderHandlers {
		self.orderManager.RegisterOrderChan(cp, c)
	}

	// start the utxo manager
	c := make(chan bool)
	go self.btcum.Start(c)
	go self.skyum.Start(c)

	go self.orderManager.Start(1*time.Second, c)
	self.handleOrders(c)

	// start the api server.
	// r := NewRouter(self)
	r := router.New(self, c)
	r.Run(self.cfg.Port)
}

// GetBtcFee get transaction fee of bitcoin.
func (self ExchangeServer) GetBtcFee() uint64 {
	return uint64(self.cfg.BtcFee)
}

// GetServPrivKey returnt he sever's private key.
func (self ExchangeServer) GetServPrivKey() cipher.SecKey {
	return self.cfg.Seckey
}

// GetPrivKey get the private key of specific address.
func (self ExchangeServer) GetAddrPrivKey(ct coin.Type, addr string) (string, error) {
	entry, err := self.wallet.GetAddressEntry(ct, addr)
	if err != nil {
		return "", err
	}

	return entry.Secret, nil
}

// GetNewAddress create new address of specific coin type.
func (self *ExchangeServer) GetNewAddress(coinType coin.Type) string {
	self.wltMtx.Lock()
	defer self.wltMtx.Unlock()
	addrEntry, err := self.wallet.NewAddresses(coinType, 1)
	if err != nil {
		panic("server get new address failed")
	}
	return addrEntry[0].Address
}

// ChooseUtxos choose appropriate bitcoin utxos,
func (self *ExchangeServer) ChooseUtxos(ct coin.Type, amount uint64, tm time.Duration) (interface{}, error) {
	switch ct {
	case coin.Bitcoin:
		return self.btcum.ChooseUtxos(amount, tm)
	case coin.Skycoin:
		return self.skyum.ChooseUtxos(amount, tm)
	default:
		return nil, errors.New("unknow coin type")
	}
}

// PutUtxos set back the utxos of specific coin type.
func (self *ExchangeServer) PutUtxos(ct coin.Type, utxos interface{}) {
	switch ct {
	case coin.Bitcoin:
		btcUtxos := utxos.([]bitcoin.Utxo)
		for _, u := range btcUtxos {
			self.btcum.PutUtxo(u)
		}
	case coin.Skycoin:
		skyUtxos := utxos.([]skycoin.Utxo)
		for _, u := range skyUtxos {
			self.skyum.PutUtxo(u)
		}
	}
}

// AddWatchAddress add watch address to utxo manager.
func (self *ExchangeServer) WatchAddress(ct coin.Type, addr string) {
	switch ct {
	case coin.Bitcoin:
		self.btcum.WatchAddresses([]string{addr})
	case coin.Skycoin:
		self.skyum.WatchAddresses([]string{addr})
	}
}

func (self *ExchangeServer) SaveAccount() error {
	return self.Save()
}

func (self *ExchangeServer) AddOrder(cp string, odr order.Order) (uint64, error) {
	return self.orderManager.AddOrder(cp, odr)
}

// initDataDir init the data dir of skycoin exchange.
func initDataDir(dir string) string {
	if dir == "" {
		logger.Error("data directory is nil")
	}

	home := util.UserHome()
	if home == "" {
		logger.Warning("Failed to get home directory")
		dir = filepath.Join("./", dir)
	} else {
		dir = filepath.Join(home, dir)
	}

	if err := os.MkdirAll(dir, os.FileMode(0700)); err != nil {
		logger.Error("Failed to create directory %s: %v", dir, err)
	}
	return dir
}

func (self *ExchangeServer) handleOrders(c chan bool) {
	for cp, ch := range self.orderHandlers {
		go func(cp string, ch chan order.Order, closing chan bool) {
			for {
				select {
				case <-closing:
					return
				case order := <-ch:
					// handle the order
					self.settleOrder(cp, order)
				}
			}
		}(cp, ch, c)
	}
}

func (self *ExchangeServer) settleOrder(cp string, od order.Order) {
	logger.Info("match order=== type:%s, price:%d, amount:%d", od.Type, od.Price, od.Amount)
	acnt, err := self.GetAccount(od.AccountID)
	if err != nil {
		panic("error account id")
	}

	pair := strings.Split(cp, "/")
	if len(pair) != 2 {
		panic("error coin pair")
	}
	mainCt, err := coin.TypeFromStr(pair[0])
	if err != nil {
		panic(err)
	}
	subCt, err := coin.TypeFromStr(pair[1])
	if err != nil {
		panic(err)
	}

	switch od.Type {
	case order.Bid:
		// increase main coin balance
		logger.Info("account:%s increase %s:%d", od.AccountID, mainCt, od.Amount)
		if err := acnt.IncreaseBalance(mainCt, od.Amount); err != nil {
			panic(err)
		}

		self.SaveAccount()
	case order.Ask:
		// increase sub coin balance.
		logger.Info("account:%s increase %s:%d", od.AccountID, subCt, od.Price*od.Amount)
		if err := acnt.IncreaseBalance(subCt, od.Price*od.Amount); err != nil {
			panic(err)
		}
		// decrease main coin balance.
		logger.Info("account:%s decrease %s:%d", od.AccountID, mainCt, od.Amount)
		if err := acnt.DecreaseBalance(mainCt, od.Amount); err != nil {
			panic(err)
		}
		self.SaveAccount()
	}
}

func (self *ExchangeServer) GetOrders(cp string, tp order.Type, start, end int64) ([]order.Order, error) {
	return self.orderManager.GetOrders(cp, tp, start, end)
}

func (self ExchangeServer) GetSupportCoins() []string {
	return self.coins
}
