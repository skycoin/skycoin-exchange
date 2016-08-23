package api

import (
	"errors"
	"net/http"

	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/wallet"
)

// CreateWallet api for creating local wallet.
// mode: POST
// url: /api/v1/wallet
// url params:
//    coin_type: bitcoin or skycoin
//    seed: wallet seed.
func CreateWallet(se Servicer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rlt := &pp.EmptyRes{}
		for {
			// check method
			if r.Method != "POST" {
				rlt = pp.MakeErrRes(errors.New("require POST method"))
				break
			}

			// get coin type
			cp, err := coin.TypeFromStr(r.FormValue("coin_type"))
			if err != nil {
				rlt = pp.MakeErrRes(err)
				break
			}

			// get seed
			sd := r.FormValue("seed")
			if sd == "" {
				rlt = pp.MakeErrRes(errors.New("no seed"))
				break
			}

			wlt, err := wallet.New(cp, sd)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			res := struct {
				Result *pp.Result `json:"result"`
				ID     string     `json:"id"`
			}{
				Result: pp.MakeResultWithCode(pp.ErrCode_Success),
				ID:     wlt.GetID(),
			}
			sendJSON(w, &res)
			return
		}
		sendJSON(w, rlt)
	}
}
