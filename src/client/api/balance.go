package api

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/skycoin/skycoin-exchange/src/client/account"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// GetBalance get balance of specific account through exchange server.
func GetBalance(se Servicer) httprouter.Handle {
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

			cp := r.FormValue("coin_type")
			if cp == "" {
				err := errors.New("coin type empty")
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			gbr := pp.GetBalanceReq{
				AccountId: pp.PtrString(pubkey),
				CoinType:  pp.PtrString(cp),
			}

			req, err := makeEncryptReq(&gbr, se.GetServKey().Hex(), a.Seckey)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			resp, err := sknet.Get(se.GetServAddr(), "/auth/get/balance", req)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			res, err := decodeRsp(resp.Body, se.GetServKey().Hex(), a.Seckey, &pp.GetBalanceRes{})
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
