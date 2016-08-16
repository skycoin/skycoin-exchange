package api

import (
	"encoding/json"
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
			rawReq := pp.OrderReq{}
			if err := BindJSON(r, &rawReq); err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			id, key, err := getAccountAndKey(r)
			if err != nil {
				rlt = pp.MakeErrRes(err)
				break
			}

			rawReq.AccountId = &id
			req, _ := pp.MakeEncryptReq(&rawReq, se.GetServKey().Hex(), key)
			resp, err := sknet.Get(se.GetServAddr(), fmt.Sprintf("/auth/create/order"), req)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			res := pp.EncryptRes{}
			json.NewDecoder(resp.Body).Decode(&res)

			// handle the response
			if res.Result.GetSuccess() {
				v := pp.OrderRes{}
				pp.DecryptRes(res, se.GetServKey().Hex(), key, &v)
				SendJSON(w, &v)
				return
			} else {
				SendJSON(w, &res)
				return
			}
		}
		SendJSON(w, rlt)
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
			_, key, err := getAccountAndKey(r)
			if err != nil {
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
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			end, err := strconv.ParseInt(ed, 10, 64)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			getOrderReq := pp.GetOrderReq{
				CoinPair: &cp,
				Type:     pp.PtrString(tp),
				Start:    &start,
				End:      &end,
			}

			req, err := pp.MakeEncryptReq(&getOrderReq, se.GetServKey().Hex(), key)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			resp, err := sknet.Get(se.GetServAddr(), "/auth/get/orders", req)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			res := pp.EncryptRes{}
			json.NewDecoder(resp.Body).Decode(&res)

			// handle the response
			if res.Result.GetSuccess() {
				v := pp.CoinsRes{}
				pp.DecryptRes(res, se.GetServKey().Hex(), key, &v)
				SendJSON(w, &v)
				return
			} else {
				SendJSON(w, &res)
				return
			}
		}
		SendJSON(w, rlt)
	}
}
