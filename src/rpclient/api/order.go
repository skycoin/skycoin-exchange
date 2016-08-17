package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

func CreateOrder(se Servicer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rlt := &pp.EmptyRes{}
		for {
			if r.Method != "POST" {
				logger.Error("require POST method")
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			rawReq := pp.OrderReq{}
			if err := bindJSON(r, &rawReq); err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			id, key, err := getAccountAndKey(r)
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrRes(err)
				break
			}

			rawReq.AccountId = &id
			req, err := makeEncryptReq(&rawReq, se.GetServKey().Hex(), key)
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			resp, err := sknet.Get(se.GetServAddr(), fmt.Sprintf("/auth/create/order"), req)
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			v, err := decodeRsp(resp.Body, se.GetServKey().Hex(), key, &pp.OrderRes{})
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			sendJSON(w, v)
			return
		}
		sendJSON(w, rlt)
	}
}

func GetBidOrders(se Servicer) http.HandlerFunc {
	return getOrders(se, "bid")
}

func GetAskOrders(se Servicer) http.HandlerFunc {
	return getOrders(se, "ask")
}

func getOrders(se Servicer, tp string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rlt := &pp.EmptyRes{}
		for {
			if r.Method != "GET" {
				logger.Error("require GET method")
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			_, key, err := getAccountAndKey(r)
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrRes(err)
				break
			}
			cp := r.URL.Query().Get("coin_pair")
			st := r.URL.Query().Get("start")
			ed := r.URL.Query().Get("end")
			if cp == "" || st == "" || ed == "" || tp == "" {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			start, err := strconv.ParseInt(st, 10, 64)
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			end, err := strconv.ParseInt(ed, 10, 64)
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			getOrderReq := pp.GetOrderReq{
				CoinPair: &cp,
				Type:     pp.PtrString(tp),
				Start:    &start,
				End:      &end,
			}

			req, err := makeEncryptReq(&getOrderReq, se.GetServKey().Hex(), key)
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			resp, err := sknet.Get(se.GetServAddr(), "/auth/get/orders", req)
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			res, err := decodeRsp(resp.Body, se.GetServKey().Hex(), key, &pp.GetOrderRes{})
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			sendJSON(w, res)
			return
		}
		sendJSON(w, rlt)
	}
}
