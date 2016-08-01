package api

import (
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/order"
)

func BidOrder(egn engine.Exchange) gin.HandlerFunc {
	return func(c *gin.Context) {
		rlt := &pp.EmptyRes{}
		req := &pp.OrderReq{}
		for {
			if err := getRequest(c, req); err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			aid := hex.EncodeToString(req.GetAccountId())
			odr := order.New(aid, order.Bid, req.GetPrice(), req.GetAmount())
			oid, err := egn.AddOrder(req.GetCoinPair(), *odr)
			if err != nil {
				glog.Error(err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			glog.Info("new order:%d", oid)
			res := pp.OrderRes{
				AccountId: req.AccountId,
				OrderId:   &oid,
			}
			reply(c, res)
			return
		}

		c.JSON(200, *rlt)
	}
}

func AskOrder(egn engine.Exchange) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetOrderbook(egn engine.Exchange) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
