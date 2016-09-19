package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/skycoin/skycoin-exchange/src/client/account"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// CreateOrder create order through exchange server.
// mode: POST
// url: /api/v1/account/order?coin_pair=[:coin_pair]&type=[:type]&price=[:price]&amt=[:amt]
// params:
// 		coin_pair: order coin pair.
// 		type: order type, can be bid or ask.
// 		price: price.
// 		amt: amount.
func CreateOrder(se Servicer) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		rlt := &pp.EmptyRes{}
		for {
			rawReq, err := makeOrderReq(r)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			a, err := account.GetActive()
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			rawReq.Pubkey = pp.PtrString(a.Pubkey)
			req, err := makeEncryptReq(rawReq, se.GetServKey().Hex(), a.Seckey)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			resp, err := sknet.Get(se.GetServAddr(), fmt.Sprintf("/auth/create/order"), req)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			v, err := decodeRsp(resp.Body, se.GetServKey().Hex(), a.Seckey, &pp.OrderRes{})
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			sendJSON(w, v)
			return
		}
		sendJSON(w, rlt)
	}
}

func makeOrderReq(r *http.Request) (*pp.OrderReq, error) {
	// get coin_pair
	cp := r.FormValue("coin_pair")
	if cp == "" {
		return nil, errors.New("coin_pair is empty")
	}

	// get order type
	tp := r.FormValue("type")
	if tp == "" {
		return nil, errors.New("type is empty")
	}

	// get price
	pc := r.FormValue("price")
	if pc == "" {
		return nil, errors.New("price is empty")
	}
	price, err := strconv.ParseUint(pc, 10, 64)
	if err != nil {
		return nil, err
	}

	// get amount
	amt := r.FormValue("amt")
	if amt == "" {
		return nil, errors.New("amt is empty")
	}
	amount, err := strconv.ParseUint(amt, 10, 64)
	if err != nil {
		return nil, err
	}

	return &pp.OrderReq{
		CoinPair: pp.PtrString(cp),
		Type:     pp.PtrString(tp),
		Price:    pp.PtrUint64(price),
		Amount:   pp.PtrUint64(amount),
	}, nil
}

// GetBidOrders get bid orders through exchange server.
func GetBidOrders(se Servicer) httprouter.Handle {
	return getOrders(se, "bid")
}

// GetAskOrders get ask orders through exchange server.
func GetAskOrders(se Servicer) httprouter.Handle {
	return getOrders(se, "ask")
}

func getOrders(se Servicer, tp string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		rlt := &pp.EmptyRes{}
		for {
			a, err := account.GetActive()
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}
			cp := r.FormValue("coin_pair")
			st := r.FormValue("start")
			ed := r.FormValue("end")
			if cp == "" || st == "" || ed == "" || tp == "" {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			start, err := strconv.ParseInt(st, 10, 64)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			end, err := strconv.ParseInt(ed, 10, 64)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			getOrderReq := pp.GetOrderReq{
				CoinPair: &cp,
				Type:     pp.PtrString(tp),
				Start:    &start,
				End:      &end,
			}

			req, err := makeEncryptReq(&getOrderReq, se.GetServKey().Hex(), a.Seckey)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			resp, err := sknet.Get(se.GetServAddr(), "/auth/get/orders", req)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			res, err := decodeRsp(resp.Body, se.GetServKey().Hex(), a.Seckey, &pp.GetOrderRes{})
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			sendJSON(w, res)
			return
		}
		sendJSON(w, rlt)
	}
}
