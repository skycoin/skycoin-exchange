package server

import (
	"github.com/skycoin/skycoin-exchange/src/server/api"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/net"
)

func NewNet(ee engine.Exchange, quit chan bool) *net.Engine {
	nt := net.New(quit)
	auth := nt.Group("/auth", api.Authorize(ee))
	{
		auth.Register("/create/account", api.CreateAccount(ee))
		auth.Register("/create/deposit_address", api.GetNewAddress(ee))
		auth.Register("/get/balance", api.GetBalance(ee))
		auth.Register("/withdrawl", api.Withdraw(ee))
		auth.Register("/create/order", api.CreateOrder(ee))
	}
	nt.Register("/get/coins", api.GetCoins(ee))
	nt.Register("/get/orders", api.GetOrders(ee))
	return nt
}
