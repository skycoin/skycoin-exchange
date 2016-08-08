package server

import (
	"github.com/skycoin/skycoin-exchange/src/server/api"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/net"
)

func NewNet(ee engine.Exchange, quit chan bool) *net.Engine {
	nt := net.New(quit)
	nt.Use(api.Authorize(ee))
	nt.Register("/get/coins", api.GetCoins(ee))
	return nt
}
