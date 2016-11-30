package router

import (
	"github.com/skycoin/skycoin-exchange/src/server/api"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// New create sknet engine and register handlers.
func New(ee engine.Exchange, quit chan bool) *sknet.Engine {
	engine := sknet.New(ee.GetSecKey(), quit)
	engine.Use(sknet.Logger())

	engine.Register("/create/account", api.CreateAccount(ee))
	engine.Register("/create/deposit_address", api.GetNewAddress(ee))
	engine.Register("/get/account/balance", api.GetAccountBalance(ee))
	engine.Register("/get/address/balance", api.GetAddrBalance(ee))
	engine.Register("/withdrawl", api.Withdraw(ee))
	engine.Register("/create/order", api.CreateOrder(ee))
	engine.Register("/get/coins", api.GetCoins(ee))
	engine.Register("/get/orders", api.GetOrders(ee))

	// utxos handler
	engine.Register("/get/utxos", api.GetUtxos(ee))

	// output history handler
	engine.Register("/get/output", api.GetOutput(ee))

	// transaction handler
	engine.Register("/inject/tx", api.InjectTx(ee))
	engine.Register("/get/tx", api.GetTx(ee))
	engine.Register("/get/rawtx", api.GetRawTx(ee))

	engine.Register("/admin/update/credit", api.UpdateCredit(ee))

	return engine
}
