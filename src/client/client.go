package client

import (
	"fmt"
	"time"

	logging "github.com/op/go-logging"
	"github.com/skycoin/skycoin-exchange/src/client/account"
	"github.com/skycoin/skycoin-exchange/src/client/router"
	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/wallet"
	gui "github.com/skycoin/skycoin-exchange/src/web-app"
	"github.com/skycoin/skycoin/src/util/browser"
)

var logger = logging.MustGetLogger("client.rpclient")

// Config client coinfig
type Config struct {
	ServAddr   string
	Port       int
	GuiDir     string
	AccountDir string
	WalletDir  string
}

// Service rpc client service.
type Service struct {
	cfg   Config
	coins map[string]coin.Gateway
}

// New create client service
func New(cfg Config) *Service {
	return &Service{
		cfg:   cfg,
		coins: make(map[string]coin.Gateway),
	}
}

// BindCoins register coins.
func (se *Service) BindCoins(cs ...coin.Gateway) error {
	for _, c := range cs {
		// check if the coin already registered
		if _, exist := se.coins[c.Type()]; exist {
			return fmt.Errorf("%s coin already registered", c.Type())
		}
		se.coins[c.Type()] = c
	}
	return nil
}

// GetCoin gets coin gatway of specific type.
func (se Service) GetCoin(coinType string) (coin.Gateway, error) {
	c, ok := se.coins[coinType]
	if !ok {
		return nil, fmt.Errorf("%s coin is not supported", coinType)
	}

	return c, nil
}

// GetServAddr get exchange server addresse.
func (se Service) GetServAddr() string {
	return se.cfg.ServAddr
}

// Run start the client service.
func (se *Service) Run() {
	// init wallet dir
	wallet.InitDir(se.cfg.WalletDir)

	// init account dir
	account.InitDir(se.cfg.AccountDir)

	// register coins
	r := router.New(se)
	addr := fmt.Sprintf("localhost:%d", se.cfg.Port)
	if err := gui.LaunchWebInterface(addr, se.cfg.GuiDir, r); err != nil {
		panic(err)
	}

	go func() {
		// Wait a moment just to make sure the http service is up
		time.Sleep(time.Millisecond * 100)
		fulladdress := fmt.Sprintf("http://%s", addr)
		logger.Info("Launching System Browser with %s", fulladdress)
		if err := browser.Open(fulladdress); err != nil {
			logger.Error(err.Error())
		}
	}()
}
