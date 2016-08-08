package api

import (
	"errors"
	"fmt"
	"strings"

	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/net"
	"github.com/skycoin/skycoin-exchange/src/server/order"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

func CreateOrder(egn engine.Exchange) net.HandlerFunc {
	return func(c *net.Context) {
		rlt := &pp.EmptyRes{}
		req := &pp.OrderReq{}
		for {
			tp, err := order.TypeFromStr(req.GetType())
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

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
				logger.Info("account:%s decrease %s:%d", acnt.GetID(), ct, bal)
				if err := acnt.DecreaseBalance(ct, bal); err != nil {
					rlt = pp.MakeErrRes(err)
					break
				}
			}

			odr := order.New(aid, tp, req.GetPrice(), req.GetAmount())
			oid, err := egn.AddOrder(req.GetCoinPair(), *odr)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			success = true
			logger.Info(fmt.Sprintf("new %s order:%d", tp, oid))
			res := pp.OrderRes{
				Result:    pp.MakeResultWithCode(pp.ErrCode_Success),
				AccountId: req.AccountId,
				OrderId:   &oid,
			}
			reply(c, res)
			return
		}
		c.JSON(rlt)
	}
}

func GetOrders(egn engine.Exchange) net.HandlerFunc {
	return func(c *net.Context) {
		rlt := &pp.EmptyRes{}
		for {
			req := pp.GetOrderReq{}
			if err := c.BindJSON(&req); err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			tp, err := order.TypeFromStr(req.GetType())
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
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
			c.JSON(res)
			return
		}

		c.JSON(rlt)
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
