package server

import "github.com/gin-gonic/gin"

func NewRouter(svr Server) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// the Authorize middle will decrypt the request, and encrypt the response.
	v1 := r.Group("/api/v1", Authorize(svr))
	{
		v1.POST("/accounts", CreateAccount(svr)) // create account
		v1.POST("/deposit_address", GetNewAddress(svr))
		v1.POST("/account/withdrawal", Withdraw(svr))
		v1.GET("/account/:id/balance", GetBalance(svr))
	}
	return r
}
