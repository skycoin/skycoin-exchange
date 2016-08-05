package rpclient

import (
	"log"
	"net/http"

	"github.com/skycoin/skycoin/src/cipher"
)

type Client interface {
	Run(addr string)
	GetServApiRoot() string
	GetServPubkey() cipher.PubKey
}

type RpcClient struct {
	Cfg Config
}

type Config struct {
	ApiRoot    string
	ServPubkey cipher.PubKey
}

func New(cfg Config) Client {
	return &RpcClient{
		Cfg: cfg,
	}
}

func (rc RpcClient) GetServApiRoot() string {
	return rc.Cfg.ApiRoot
}

func (rc RpcClient) GetServPubkey() cipher.PubKey {
	return rc.Cfg.ServPubkey
}

func (rc *RpcClient) Run(addr string) {
	r := NewRouter(rc)
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
