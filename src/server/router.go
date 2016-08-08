package server

// func NewRouter(ee engine.Exchange) *gin.Engine {
// 	r := gin.New()
// 	r.Use(Logger())
// 	r.Use(gin.Recovery())
//
// 	// the Authorize middle will decrypt the request, and encrypt the response.
// 	v1 := r.Group("/api/v1")
// 	authReq := v1.Group("/", api.Authorize(ee))
// 	{
// 		authReq.POST("/accounts", api.CreateAccount(ee)) // create account
// 		authReq.POST("/deposit_address", api.GetNewAddress(ee))
// 		authReq.POST("/account/withdrawal", api.Withdraw(ee))
// 		authReq.POST("/account/balance", api.GetBalance(ee))
//
// 		authReq.POST("/account/order/:type", api.CreateOrder(ee))
// 	}
//
// 	v1.POST("/orders/:type", api.GetOrders(ee))
// 	// v1.GET("/coins", api.GetCoins(ee))
// 	return r
// }
