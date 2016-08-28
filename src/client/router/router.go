package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/skycoin/skycoin-exchange/src/client/api"
)

// New return router
func New(se api.Servicer) *httprouter.Router {
	rt := httprouter.New()
	registerBaseHandlers(rt, se)
	registerOrderHandlers(rt, se)
	registerUtxoHandlers(rt, se)
	registerTxnHandlers(rt, se)
	registerWalletHandlers(rt, se)
	registerAdminHandlers(rt, se)
	return rt
}

// base handlers.
func registerBaseHandlers(rt *httprouter.Router, se api.Servicer) {
	rt.GET("/api/v1/coins", api.GetCoins(se))
	rt.POST("/api/v1/accounts", api.CreateAccount(se))
	rt.PUT("/api/v1/account/session", api.ActiveAccount(se))
	rt.POST("/api/v1/account/deposit_address", api.GetDepositAddress(se))
	rt.GET("/api/v1/account/balance", api.GetBalance(se))
	rt.POST("/api/v1/account/withdrawal", api.Withdraw(se))
}

// order handlers
func registerOrderHandlers(rt *httprouter.Router, se api.Servicer) {
	rt.POST("/api/v1/account/order", api.CreateOrder(se))
	rt.GET("/api/v1/orders/bid", api.GetBidOrders(se))
	rt.GET("/api/v1/orders/ask", api.GetAskOrders(se))
}

// utxos handlers
func registerUtxoHandlers(rt *httprouter.Router, se api.Servicer) {
	rt.GET("/api/v1/utxos", api.GetUtxos(se))
}

// transaction handlers.
func registerTxnHandlers(rt *httprouter.Router, se api.Servicer) {
	rt.POST("/api/v1/inject_tx", api.InjectTx(se))
	rt.GET("/api/v1/tx", api.GetTx(se))
	rt.GET("/api/v1/rawtx", api.GetRawTx(se))
}

// wallet handlers.
func registerWalletHandlers(rt *httprouter.Router, se api.Servicer) {
	rt.POST("/api/v1/wallet", api.CreateWallet(se))
	rt.POST("/api/v1/wallet/addresses", api.NewAddress(se))
	rt.GET("/api/v1/wallet/addresses", api.GetAddresses(se))
	rt.GET("/api/v1/wallet/address/key", api.GetKeys(se))
}

// admin handlers.
func registerAdminHandlers(rt *httprouter.Router, se api.Servicer) {
	rt.PUT("/api/v1/admin/credit", api.UpdateCredit(se))
}
