package server

import "github.com/gin-gonic/gin"

func NewRouter(svr Server) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		v1.POST("/accounts", CreateAccount(svr)) // create account

		v1.POST("/authorization", Authorize(svr)) // authorize account

		authorized := v1.Group("/account", AuthRequired(svr), Security(svr))
		{
			authorized.POST("/deposit_address", GetNewAddress(svr)) // get new address from account.
			authorized.POST("/withdraw", Withdraw(svr))
		}

	}
	return r
}
