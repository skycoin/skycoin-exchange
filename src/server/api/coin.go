package api

import (
	"github.com/gin-gonic/gin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
)

func GetCoins(egn engine.Exchange) gin.HandlerFunc {
	return func(c *gin.Context) {
		coins := pp.CoinsRes{
			Result: pp.MakeResultWithCode(pp.ErrCode_Success),
			Coins:  egn.GetSupportCoins(),
		}
		c.JSON(200, coins)
	}
}
