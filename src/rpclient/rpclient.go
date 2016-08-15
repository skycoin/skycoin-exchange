package rpclient

import (
	"github.com/skycoin/skycoin-exchange/src/rpclient/model"
	"github.com/skycoin/skycoin/src/cipher"
)

type Config struct {
	ApiRoot    string
	ServPubkey cipher.PubKey
}

func New(cfg Config) *model.Service {
	return &model.Service{
		ServAddr:   cfg.ApiRoot,
		ServPubkey: cfg.ServPubkey,
	}
}
