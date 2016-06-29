package server

import "github.com/gin-gonic/gin"

func NewRouter(svr *Server) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		v1.POST("/account", CreateAccount(svr)) // create account

		v1.POST("/account/deposit", GetNewAddress(svr)) // get new address from account.
	}

	// r.Use(Encrypt())

	return r
}
