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
			cp := r.FormValue("coin_type")
			if cp == "" {
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
				CoinType:  pp.PtrString(cp),
				Addresses: addrArray,
			}

			var res pp.GetUtxoRes
			if err := sknet.EncryGet(se.GetServAddr(), "/get/utxos", req, &res); err != nil {
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
