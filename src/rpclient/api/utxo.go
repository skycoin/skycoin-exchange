package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// GetUtxos get utxos through exchange server.
func GetUtxos(se Servicer) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var rlt *pp.EmptyRes
		for {
			if r.Method != "GET" {
				logger.Error("require GET method")
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			_, key, err := getAccountAndKey(r)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			tp := r.FormValue("coin_type")
			if tp == "" {
				logger.Error("coin type empty")
				rlt = pp.MakeErrRes(errors.New("coin type empty"))
				break
			}

			addrs := r.FormValue("addrs")
			if addrs == "" {
				logger.Error("addrs empty")
				rlt = pp.MakeErrRes(errors.New("addrs empty"))
				break
			}
			addrArray := strings.Split(addrs, ",")
			for i, addr := range addrArray {
				addrArray[i] = strings.Trim(addr, " ")
			}

			req := pp.GetUtxoReq{
				CoinType:  pp.PtrString(tp),
				Addresses: addrArray,
			}
			encReq, err := makeEncryptReq(&req, se.GetServKey().Hex(), key)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			resp, err := sknet.Get(se.GetServAddr(), "/auth/get/utxos", encReq)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			res, err := decodeRsp(resp.Body, se.GetServKey().Hex(), key, &pp.GetUtxoRes{})
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
