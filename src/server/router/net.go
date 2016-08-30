package router

import (
	"github.com/skycoin/skycoin-exchange/src/server/api"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// New create sknet engine and register handlers.
func New(ee engine.Exchange, quit chan bool) *sknet.Engine {
	nt := sknet.New(quit)
	nt.Use(sknet.Recovery())
	nt.Use(sknet.Logger())

	auth := nt.Group("/auth", api.Authorize(ee))
	{
		// baseic handlers.
		auth.Register("/create/account", api.CreateAccount(ee))
		auth.Register("/create/deposit_address", api.GetNewAddress(ee))
		auth.Register("/get/account/balance", api.GetAccountBalance(ee))
		auth.Register("/get/address/balance", api.GetAddrBalance(ee))
		auth.Register("/withdrawl", api.Withdraw(ee))
		auth.Register("/create/order", api.CreateOrder(ee))
		auth.Register("/get/coins", api.GetCoins(ee))
		auth.Register("/get/orders", api.GetOrders(ee))

		// utxos handler
		auth.Register("/get/utxos", api.GetUtxos(ee))

		// transaction handler
		auth.Register("/inject/tx", api.InjectTx(ee))
		auth.Register("/get/tx", api.GetTx(ee))
		auth.Register("/get/rawtx", api.GetRawTx(ee))
	}

	admin := nt.Group("/admin", api.Authorize(ee), api.IsAdmin(ee))
	{
		admin.Register("/update/credit", api.UpdateCredit(ee))
	}

	return nt
}
