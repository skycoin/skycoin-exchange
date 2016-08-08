package server

import (
	"github.com/skycoin/skycoin-exchange/src/server/api"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/net"
)

func NewNet(ee engine.Exchange) *net.Engine {
	engine := net.New()
	engine.Register("/getcoins", api.GetCoins(ee))
	return engine
}
