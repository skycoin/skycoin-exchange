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
			a, err := account.GetActive()
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

			req := pp.GetAccountBalanceReq{
				Pubkey:   pp.PtrString(a.Pubkey),
				CoinType: pp.PtrString(cp),
			}

			var res pp.GetAccountBalanceRes

			if err := sknet.EncryGet(se.GetServAddr(), "/get/account/balance", req, &res); err != nil {
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
