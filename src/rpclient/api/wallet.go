package api

import (
	"net/http"

	"github.com/skycoin/skycoin-exchange/src/pp"
)

// CreateWallet api for creating local wallet.
// mode: POST
// url: /api/v1/wallet
// url params:
//    coin_type: bitcoin or skycoin
//    seed: wallet seed.
func CreateWallet(se Servicer) http.HandlerFunc {
	return func(c http.ResponseWriter, r *http.Request) {
		rlt = pp.EmptyRes{}
		for {
			// check method
			if r.Method != "POST" {
				rlt = pp.MakeErrRes("require POST method")
				break
			}

			// get coin type
			cp := r.FormValue("coin_type")
			if cp == "" {
				rlt = pp.MakeErrRes("no coin type")
				break
			}

			// get seed
			sd := r.FormValue("seed")
			if sd == "" {
				rlt = pp.MakeErrRes("no seed")
				break
			}

		}
	}
}
