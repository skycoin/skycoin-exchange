package client

import (
	"fmt"
	"time"

	"github.com/skycoin/skycoin-exchange/src/client/account"
	"github.com/skycoin/skycoin-exchange/src/client/router"
	"github.com/skycoin/skycoin-exchange/src/wallet"
	gui "github.com/skycoin/skycoin-exchange/src/web-app"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
	"gopkg.in/op/go-logging.v1"
)

var logger = logging.MustGetLogger("client.rpclient")

// Config client coinfig
type Config struct {
	ServAddr   string
	ServPubkey cipher.PubKey
	Port       int
	GuiDir     string
	AccountDir string
	WalletDir  string
}

// New create client service
func New(cfg Config) *Service {
	return &Service{
		cfg,
	}
}

// Service rpc client service.
type Service struct {
	cfg Config
}

// GetServKey get server pubkey.
func (se Service) GetServKey() cipher.PubKey {
	return se.cfg.ServPubkey
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
		if err := util.OpenBrowser(fulladdress); err != nil {
			logger.Error(err.Error())
		}
	}()
}
