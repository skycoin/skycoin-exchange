package rpclient

import "github.com/gin-gonic/gin"

func NewRouter(cli Client) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// the Authorize middle will decrypt the request, and encrypt the response.
	v1 := r.Group("/api/v1")
	{
		v1.POST("/accounts", CreateAccount(cli))
		v1.GET("/account/deposit_address", GetNewAddress(cli))
		v1.GET("/account/balance", GetBalance(cli))
		v1.GET("/account/withdrawal", Withdraw(cli))

		v1.POST("/account/bid", BidOrder(cli))
		v1.POST("/account/ask", AskOrder(cli))
		v1.GET("orders/:type", GetOrderBook(cli))
	}
	return r
}
