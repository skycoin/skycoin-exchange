package api

import (
	"time"

	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/sknet"
	"github.com/skycoin/skycoin/src/cipher"
)

// CreateAccount create account with specific pubkey,
func CreateAccount(ee engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) {
		errRlt := &pp.EmptyRes{}
		for {
			req := pp.CreateAccountReq{}
			if err := getRequest(c, &req); err != nil {
				logger.Error(err.Error())
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			if _, err := cipher.PubKeyFromHex(req.GetPubkey()); err != nil {
				logger.Error(err.Error())
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongAccountId)
				break
			}

			// create account with pubkey.
			if _, err := ee.CreateAccountWithPubkey(req.GetPubkey()); err != nil {
				logger.Error(err.Error())
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			res := pp.CreateAccountRes{
				Result:    pp.MakeResultWithCode(pp.ErrCode_Success),
				Pubkey:    req.Pubkey,
				CreatedAt: pp.PtrInt64(time.Now().Unix()),
			}

			reply(c, res)
			return
		}

		c.JSON(errRlt)
	}
}
