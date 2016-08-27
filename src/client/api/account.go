package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/skycoin/skycoin-exchange/src/client/account"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// CreateAccount handle the request of creating account.
func CreateAccount(se Servicer) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// generate account pubkey/privkey pair, pubkey is the account id.
		errRlt := &pp.EmptyRes{}
		for {
			a := account.New()
			r := pp.CreateAccountReq{
				Pubkey: pp.PtrString(a.Pubkey),
			}

			req, err := makeEncryptReq(&r, se.GetServKey().Hex(), a.Seckey)
			if err != nil {
				logger.Error(err.Error())
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			rsp, err := sknet.Get(se.GetServAddr(), "/auth/create/account", req)
			if err != nil {
				logger.Error(err.Error())
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			res, err := decodeRsp(rsp.Body, se.GetServKey().Hex(), a.Seckey, &pp.CreateAccountRes{})
			if err != nil {
				logger.Error(err.Error())
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			acntRes := res.(*pp.CreateAccountRes)
			if !acntRes.GetResult().GetSuccess() {
				sendJSON(w, res)
			} else {
				ret := struct {
					Result    pp.Result `json:"result"`
					Pubkey    string    `json:"pubkey"`
					CreatedAt int64     `json:"created_at"`
				}{
					Result:    *acntRes.Result,
					Pubkey:    a.Pubkey,
					CreatedAt: acntRes.GetCreatedAt(),
				}
				account.Set(a)
				sendJSON(w, &ret)
			}
			return
		}
		sendJSON(w, errRlt)
	}
}
