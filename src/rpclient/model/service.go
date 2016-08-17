package model

import (
	"log"
	"net/http"

	"gopkg.in/op/go-logging.v1"

	"github.com/skycoin/skycoin-exchange/src/rpclient/router"
	"github.com/skycoin/skycoin/src/cipher"
)

var logger = logging.MustGetLogger("client.model")

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

func (se *Service) Run(addr string) {
	r := router.New(se)
	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Critical("%s", r)
				}
			}()
			logger.Info("client started, listen on port%s", addr)
			log.Fatal(http.ListenAndServe(addr, r))
		}()
	}
}
