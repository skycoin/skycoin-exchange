package api

import (
	"encoding/hex"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/order"
)

func BidOrder(egn engine.Exchange) gin.HandlerFunc {
	// TODO: check the balance.
	return addOrder(order.Bid, egn)
}

func AskOrder(egn engine.Exchange) gin.HandlerFunc {
	// TODO: check the balance.
	return addOrder(order.Ask, egn)
}

func GetOrders(egn engine.Exchange) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get start and end.
		rlt := &pp.EmptyRes{}
		for {
			req := pp.GetOrderReq{}
			if err := c.BindJSON(&req); err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			tp := order.MustTypeFromStr(req.GetType())
			ords, err := egn.GetOrders(req.GetCoinPair(), tp, req.GetStart(), req.GetEnd())
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			res := pp.GetOrderRes{
				CoinPair: req.CoinPair,
				Type:     req.Type,
				Orders:   make([]*pp.Order, len(ords)),
			}

			for i := range ords {
				res.Orders[i] = &pp.Order{
					AccountId:   &ords[i].AccountID,
					Id:          &ords[i].ID,
					Type:        req.Type,
					Price:       &ords[i].Price,
					Amount:      &ords[i].Amount,
					RestAmt:     &ords[i].RestAmt,
					CreatedTime: &ords[i].CreatedTime,
				}
			}

			res.Result = pp.MakeResultWithCode(pp.ErrCode_Success)
			c.JSON(200, res)
			return
		}

		c.JSON(200, rlt)
	}
}

func addOrder(tp order.Type, egn engine.Exchange) gin.HandlerFunc {
	return func(c *gin.Context) {
		rlt := &pp.EmptyRes{}
		req := &pp.OrderReq{}
		for {
			if err := getRequest(c, req); err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			aid := hex.EncodeToString(req.GetAccountId())
			odr := order.New(aid, tp, req.GetPrice(), req.GetAmount())
			oid, err := egn.AddOrder(req.GetCoinPair(), *odr)
			if err != nil {
				glog.Error(err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			glog.Info(fmt.Sprintf("new %s order:%d", tp, oid))
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
