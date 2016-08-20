package rpclient

import (
	"github.com/skycoin/skycoin-exchange/src/rpclient/router"
	gui "github.com/skycoin/skycoin-exchange/src/web-app"
	"github.com/skycoin/skycoin/src/cipher"
	"gopkg.in/op/go-logging.v1"
)

var logger = logging.MustGetLogger("client.rpclient")

type Config struct {
	ApiRoot    string
	ServPubkey cipher.PubKey
}

func New(cfg Config) *Service {
	return &Service{
		ServAddr:   cfg.ApiRoot,
		ServPubkey: cfg.ServPubkey,
	}
}

type Service struct {
	ServAddr   string        // exchange server addr.
	ServPubkey cipher.PubKey // exchagne server pubkey.
}

func (se Service) GetServKey() cipher.PubKey {
	return se.ServPubkey
}

func (se Service) GetServAddr() string {
	return se.ServAddr
}

func (se *Service) Run(addr string, guiDir string) {
	r := router.New(se)
	if err := gui.LaunchWebInterface(addr, guiDir, r); err != nil {
		panic(err)
	}
}
