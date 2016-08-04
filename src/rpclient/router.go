package rpclient

import "github.com/gin-gonic/gin"

func NewRouter(cli Client) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		v1.POST("/accounts", CreateAccount(cli))
		v1.GET("/account/deposit_address", GetNewAddress(cli))
		v1.GET("/account/balance", GetBalance(cli))
		v1.GET("/account/withdrawal", Withdraw(cli))

		v1.POST("/account/order/:type", CreateOrder(cli))
		v1.GET("orders/:type", GetOrders(cli))
		v1.GET("/coins", GetCoins(cli))
	}
	return r
}
