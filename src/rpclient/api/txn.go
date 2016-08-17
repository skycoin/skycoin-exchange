package api

import (
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// TxnHandler transaction ahndler.
func InjectTx(se Servicer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rlt *pp.EmptyRes
		for {
			if r.Method != "POST" {
				logger.Error("require POST method")
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			// get account key.
			_, key, err := getAccountAndKey(r)
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrRes(err)
				break
			}

			// get tx
			tx := r.URL.Query().Get("tx")
			if tx == "" {
				logger.Error("empty tx")
				rlt = pp.MakeErrRes(errors.New("empty tx"))
				break
			}

			// get coin type
			tp := r.URL.Query().Get("coin_type")
			if tp == "" {
				logger.Error("empty coin type")
				rlt = pp.MakeErrRes(errors.New("empty coin type"))
				break
			}

			txb, err := hex.DecodeString(tx)
			if err != nil {
				logger.Error("error tx")
				rlt = pp.MakeErrRes(errors.New("error tx"))
				break
			}

			req := pp.InjectTxnReq{
				CoinType: pp.PtrString(tp),
				Tx:       txb,
			}

			enc_req, err := makeEncryptReq(&req, se.GetServKey().Hex(), key)
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			resp, err := sknet.Get(se.GetServAddr(), "/auth/get/orders", enc_req)
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			v, err := decodeRsp(resp.Body, se.GetServKey().Hex(), key, pp.InjectTxnRes{})
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			sendJSON(w, v)
			return
		}
		sendJSON(w, rlt)
	}
}
