package rpclient

import "net/http"

func NewRouter(cli Client) *http.ServeMux {
	mux := http.NewServeMux()
	// base handlers.
	mux.Handle("/api/v1/coins", GetCoinsTcp(cli))
	mux.Handle("/api/v1/accounts", CreateAccount(cli))
	mux.Handle("/api/v1/account/deposit_address", GetNewAddress(cli))
	mux.Handle("/api/v1/account/balance", GetBalance(cli))
	mux.Handle("/api/v1/account/withdrawal", Withdraw(cli))

	// order handlers
	mux.Handle("/api/v1/account/order/bid", CreateBidOrder(cli))
	mux.Handle("/api/v1/account/order/ask", CreateAskOrder(cli))
	mux.Handle("/api/v1/orders/bid", GetBidOrders(cli))
	mux.Handle("/api/v1/orders/ask", GetAskOrders(cli))

	return mux
}
