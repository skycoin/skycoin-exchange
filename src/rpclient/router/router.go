package router

import (
	"net/http"

	"github.com/skycoin/skycoin-exchange/src/rpclient/api"
	"github.com/skycoin/skycoin-exchange/src/rpclient/model"
)

func New(cli *model.Client) *http.ServeMux {
	mux := http.NewServeMux()
	// base handlers.
	mux.Handle("/api/v1/coins", api.GetCoins(cli))
	mux.Handle("/api/v1/accounts", api.CreateAccount(cli))
	mux.Handle("/api/v1/account/deposit_address", api.GetNewAddress(cli))
	mux.Handle("/api/v1/account/balance", api.GetBalance(cli))
	mux.Handle("/api/v1/account/withdrawal", api.Withdraw(cli))

	// order handlers
	mux.Handle("/api/v1/account/order", api.CreateOrder(cli))
	mux.Handle("/api/v1/orders/bid", api.GetBidOrders(cli))
	mux.Handle("/api/v1/orders/ask", api.GetAskOrders(cli))

	return mux
}
