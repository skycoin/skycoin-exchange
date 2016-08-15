package router

import (
	"net/http"

	"github.com/skycoin/skycoin-exchange/src/rpclient/api"
)

func New(se api.Servicer) *http.ServeMux {
	mux := http.NewServeMux()
	// base handlers.
	mux.Handle("/api/v1/coins", api.GetCoins(se))
	mux.Handle("/api/v1/accounts", api.CreateAccount(se))
	mux.Handle("/api/v1/account/deposit_address", api.GetNewAddress(se))
	mux.Handle("/api/v1/account/balance", api.GetBalance(se))
	mux.Handle("/api/v1/account/withdrawal", api.Withdraw(se))

	// order handlers
	mux.Handle("/api/v1/account/order", api.CreateOrder(se))
	mux.Handle("/api/v1/orders/bid", api.GetBidOrders(se))
	mux.Handle("/api/v1/orders/ask", api.GetAskOrders(se))

	return mux
}
