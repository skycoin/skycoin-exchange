package rpclient

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/skycoin/skycoin-exchange/src/rpclient/router"
	"github.com/skycoin/skycoin-exchange/src/wallet"
	gui "github.com/skycoin/skycoin-exchange/src/web-app"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
	"gopkg.in/op/go-logging.v1"
)

var logger = logging.MustGetLogger("client.rpclient")

// Config client coinfig
type Config struct {
	APIRoot    string
	ServPubkey cipher.PubKey
}

// New create client service
func New(cfg Config) *Service {
	return &Service{
		ServAddr:   cfg.APIRoot,
		ServPubkey: cfg.ServPubkey,
	}
}

// Service rpc client service.
type Service struct {
	ServAddr   string        // exchange server addr.
	ServPubkey cipher.PubKey // exchagne server pubkey.
}

// GetServKey get server pubkey.
func (se Service) GetServKey() cipher.PubKey {
	return se.ServPubkey
}

// GetServAddr get exchange server addresse.
func (se Service) GetServAddr() string {
	return se.ServAddr
}

// Run start the client service.
func (se *Service) Run(addr string, guiDir string) {
	// init wallet
	wallet.InitDir(filepath.Join(util.UserHome(), ".exchange-client/wallet"))

	r := router.New(se)
	if err := gui.LaunchWebInterface(addr, guiDir, r); err != nil {
		panic(err)
	}
	go func() {
		// Wait a moment just to make sure the http interface is up
		time.Sleep(time.Millisecond * 100)
		fulladdress := fmt.Sprintf("http://%s", addr)
		logger.Info("Launching System Browser with %s", fulladdress)
		if err := util.OpenBrowser(fulladdress); err != nil {
			logger.Error(err.Error())
		}
	}()
}
