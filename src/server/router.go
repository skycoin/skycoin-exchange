package server

import (
	"github.com/gin-gonic/gin"
	"github.com/skycoin/skycoin-exchange/src/server/api"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
)

func NewRouter(ee engine.Exchange) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// the Authorize middle will decrypt the request, and encrypt the response.
	v1 := r.Group("/api/v1", api.Authorize(ee))
	{
		v1.POST("/accounts", api.CreateAccount(ee)) // create account
		v1.POST("/deposit_address", api.GetNewAddress(ee))
		v1.POST("/account/withdrawal", api.Withdraw(ee))
		v1.POST("/account/balance", api.GetBalance(ee))
	}
	return r
}
