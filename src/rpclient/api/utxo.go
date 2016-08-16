package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

func GetUtxos(se Servicer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rlt *pp.EmptyRes
		for {
			_, key, err := getAccountAndKey(r)
			if err != nil {
				rlt = pp.MakeErrRes(err)
				break
			}

			tp := r.URL.Query().Get("coin_type")
			if tp == "" {
				rlt = pp.MakeErrRes(errors.New("coin type empty"))
				break
			}

			addrs := r.URL.Query().Get("addrs")
			if addrs == "" {
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
			enc_req, _ := pp.MakeEncryptReq(&req, se.GetServKey().Hex(), key)
			resp, err := sknet.Get(se.GetServAddr(), "/auth/get/utxos", enc_req)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			res := pp.EncryptRes{}
			json.NewDecoder(resp.Body).Decode(&res)

			// handle the response
			if res.Result.GetSuccess() {
				v := pp.GetUtxoRes{}
				pp.DecryptRes(res, se.GetServKey().Hex(), key, &v)
				sendJSON(w, &v)
				return
			} else {
				sendJSON(w, &res)
				return
			}
		}
		sendJSON(w, rlt)
	}
}
