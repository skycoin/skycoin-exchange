package api

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/order"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

func BidOrder(egn engine.Exchange) gin.HandlerFunc {
	return addOrder(order.Bid, egn)
}

func AskOrder(egn engine.Exchange) gin.HandlerFunc {
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
					Id:        &ords[i].ID,
					Type:      req.Type,
					Price:     &ords[i].Price,
					Amount:    &ords[i].Amount,
					RestAmt:   &ords[i].RestAmt,
					CreatedAt: &ords[i].CreatedAt,
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
			aid := req.GetAccountId()
			// find the account
			if _, err := cipher.PubKeyFromHex(aid); err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongAccountId)
				break
			}

			acnt, err := egn.GetAccount(aid)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongAccountId)
				break
			}

			ct, bal, err := needBalance(tp, req)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			if acnt.GetBalance(ct) < bal {
				rlt = pp.MakeErrRes(fmt.Errorf("%s balance is not sufficient", ct))
				break
			}

			var success bool
			if tp == order.Bid {
				defer func() {
					if success {
						egn.SaveAccount()
					} else {
						acnt.IncreaseBalance(ct, bal)
					}
				}()
				// decrease the balance, in case of double use the coins.
				glog.Info(fmt.Sprintf("account:%s decrease %s:%d", acnt.GetID(), ct, bal))
				if err := acnt.DecreaseBalance(ct, bal); err != nil {
					rlt = pp.MakeErrRes(err)
					break
				}
			}

			odr := order.New(aid, tp, req.GetPrice(), req.GetAmount())
			oid, err := egn.AddOrder(req.GetCoinPair(), *odr)
			if err != nil {
				glog.Error(err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			success = true
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

func needBalance(tp order.Type, req *pp.OrderReq) (wallet.CoinType, uint64, error) {
	pair := strings.Split(req.GetCoinPair(), "/")
	if len(pair) != 2 {
		return -1, 0, errors.New("error coin pair")
	}

	mainCt, err := wallet.ConvertCoinType(pair[0])
	if err != nil {
		return -1, 0, err
	}
	subCt, err := wallet.ConvertCoinType(pair[1])
	if err != nil {
		return -1, 0, err
	}

	switch tp {
	case order.Bid:
		return subCt, req.GetPrice() * req.GetAmount(), nil
	case order.Ask:
		return mainCt, req.GetAmount(), nil
	default:
		return -1, 0, errors.New("unknow order type")
	}
}
