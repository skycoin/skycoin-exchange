package router

import (
	"github.com/skycoin/skycoin-exchange/src/server/api"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// NewNet create sknet engine and register handlers.
func New(ee engine.Exchange, quit chan bool) *sknet.Engine {
	nt := sknet.New(quit)
	nt.Use(sknet.Recovery())
	nt.Use(sknet.Logger())

	auth := nt.Group("/auth", api.Authorize(ee))
	{
		auth.Register("/create/account", api.CreateAccount(ee))
		auth.Register("/create/deposit_address", api.GetNewAddress(ee))
		auth.Register("/get/balance", api.GetBalance(ee))
		auth.Register("/withdrawl", api.Withdraw(ee))
		auth.Register("/create/order", api.CreateOrder(ee))
		auth.Register("/get/coins", api.GetCoins(ee))
		auth.Register("/get/orders", api.GetOrders(ee))

		// utxos handler
		auth.Register("/get/utxos", api.GetUtxos(ee))

		// transaction handler
		auth.Register("/inject/tx", api.InjectTx(ee))
		auth.Register("/get/tx", api.GetTx(ee))
	}

	return nt
}
