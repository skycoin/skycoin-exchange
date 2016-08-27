package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/skycoin/skycoin-exchange/src/client/account"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// GetCoins get supported coins from exchange server.
func GetCoins(se Servicer) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		rlt := &pp.EmptyRes{}
		for {
			pubkey, err := getPubkey(r)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			a, err := account.Get(pubkey)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			rq := pp.GetCoinsReq{
				AccountId: pp.PtrString(pubkey),
			}

			req, err := makeEncryptReq(&rq, se.GetServKey().Hex(), a.Seckey)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			rsp, err := sknet.Get(se.GetServAddr(), "/auth/get/coins", req)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			res, err := decodeRsp(rsp.Body, se.GetServKey().Hex(), a.Seckey, &pp.CoinsRes{})
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
