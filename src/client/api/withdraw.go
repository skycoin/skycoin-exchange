package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/skycoin/skycoin-exchange/src/client/account"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// Withdraw transaction.
func Withdraw(se Servicer) httprouter.Handle {
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
				err := errors.New("coin_type empty")
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			amount := r.FormValue("amount")
			if amount == "" {
				rlt = pp.MakeErrRes(errors.New("amount empty"))
				break
			}

			toAddr := r.FormValue("toaddr")
			if toAddr == "" {
				rlt = pp.MakeErrRes(errors.New("toaddr empty"))
				break
			}

			amt, err := strconv.ParseUint(amount, 10, 64)
			if err != nil {
				rlt = pp.MakeErrRes(err)
				break
			}
			req := pp.WithdrawalReq{
				Pubkey:        &a.Pubkey,
				CoinType:      &cp,
				Coins:         &amt,
				OutputAddress: &toAddr,
			}

			var res pp.WithdrawalRes
			if err := sknet.EncryGet(se.GetServAddr(), "/withdrawl", req, &res); err != nil {
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
