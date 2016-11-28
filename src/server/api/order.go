package api

import (
	"errors"
	"fmt"
	"strings"

	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/order"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// CreateOrder create specifc order.
func CreateOrder(egn engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) error {
		rlt := &pp.EmptyRes{}
		req := &pp.OrderReq{}
		for {
			if err := c.BindJSON(req); err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				logger.Error(err.Error())
				break
			}

			// validate pubkey
			pubkey := req.GetPubkey()
			if err := validatePubkey(pubkey); err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongPubkey)
				break
			}

			// get order type
			op, err := order.TypeFromStr(req.GetType())
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				logger.Error(err.Error())
				break
			}

			// find the account
			acnt, err := egn.GetAccount(pubkey)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongPubkey)
				logger.Error(err.Error())
				break
			}

			cp, bal, err := needBalance(op, req)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				logger.Error(err.Error())
				break
			}

			if acnt.GetBalance(cp) < bal {
				err := fmt.Errorf("%s balance is not sufficient", cp)
				rlt = pp.MakeErrRes(err)
				logger.Debug(err.Error())
				break
			}

			var success bool
			if op == order.Bid {
				defer func() {
					if success {
						egn.SaveAccount()
					} else {
						acnt.IncreaseBalance(cp, bal)
					}
				}()
				// decrease the balance, in case of double use the coins.
				logger.Info("account:%s decrease %s:%d", acnt.GetID(), cp, bal)
				if err := acnt.DecreaseBalance(cp, bal); err != nil {
					rlt = pp.MakeErrRes(err)
					logger.Error(err.Error())
					break
				}
			}

			odr := order.New(pubkey, op, req.GetPrice(), req.GetAmount())
			oid, err := egn.AddOrder(req.GetCoinPair(), *odr)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			success = true
			logger.Info(fmt.Sprintf("new %s order:%d", op, oid))
			res := pp.OrderRes{
				Result:  pp.MakeResultWithCode(pp.ErrCode_Success),
				OrderId: &oid,
			}
			return c.SendJSON(&res)
		}
		return c.Error(rlt)
	}
}

// GetOrders get order list.
func GetOrders(egn engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) error {
		rlt := &pp.EmptyRes{}
		for {
			req := pp.GetOrderReq{}
			if err := c.BindJSON(&req); err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			op, err := order.TypeFromStr(req.GetType())
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				logger.Error(err.Error())
				break
			}
			ords, err := egn.GetOrders(req.GetCoinPair(), op, req.GetStart(), req.GetEnd())
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				logger.Error(err.Error())
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
			return c.SendJSON(&res)
		}
		return c.Error(rlt)
	}
}

func needBalance(tp order.Type, req *pp.OrderReq) (coin.Type, uint64, error) {
	pair := strings.Split(req.GetCoinPair(), "/")
	if len(pair) != 2 {
		return -1, 0, errors.New("error coin pair")
	}

	mainCt, err := coin.TypeFromStr(pair[0])
	if err != nil {
		return -1, 0, err
	}
	subCt, err := coin.TypeFromStr(pair[1])
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
