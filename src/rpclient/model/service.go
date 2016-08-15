package model

import (
	"log"
	"net/http"

	"github.com/skycoin/skycoin-exchange/src/rpclient/router"
	"github.com/skycoin/skycoin/src/cipher"
)

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
					log.Println(r)
				}
			}()
			log.Println("client started ", addr)
			log.Fatal(http.ListenAndServe(addr, r))
		}()
	}
}
