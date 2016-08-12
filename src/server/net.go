package server

import (
	"github.com/skycoin/skycoin-exchange/src/server/api"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/net"
)

// NewNet create net engine and register handlers.
func NewNet(ee engine.Exchange, quit chan bool) *net.Engine {
	nt := net.New(quit)
	nt.Use(net.Recovery())
	nt.Use(net.Logger())

	auth := nt.Group("/auth", api.Authorize(ee))
	{
		auth.Register("/create/account", api.CreateAccount(ee))
		auth.Register("/create/deposit_address", api.GetNewAddress(ee))
		auth.Register("/get/balance", api.GetBalance(ee))
		auth.Register("/withdrawl", api.Withdraw(ee))
		auth.Register("/create/order", api.CreateOrder(ee))
		auth.Register("/get/coins", api.GetCoins(ee))
		auth.Register("/get/orders", api.GetOrders(ee))
	}

	return nt
}
