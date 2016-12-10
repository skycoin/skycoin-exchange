package api

import (
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// IsAdmin middleware for checking if the account is admin.
func IsAdmin(ee engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) error {
		var rlt *pp.EmptyRes
		for {
			req := pp.UpdateCreditReq{}
			if err := c.BindJSON(&req); err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			if !ee.IsAdmin(req.GetPubkey()) {
				logger.Error("not admin")
				rlt = pp.MakeErrResWithCode(pp.ErrCode_UnAuthorized)
				break
			}
			return c.Next()
		}
		return c.Error(rlt)
	}
}

// UpdateCredit update credit.
func UpdateCredit(ee engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) error {
		var rlt *pp.EmptyRes
		for {
			req := pp.UpdateCreditReq{}
			if err := c.BindJSON(&req); err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			// validate the dst pubkey.
			dstPubkey := req.GetDst()
			if err := validatePubkey(dstPubkey); err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongPubkey)
				break
			}

			// get account.
			a, err := ee.GetAccount(dstPubkey)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongPubkey)
				break
			}

			// get coin type.
			if err := a.SetBalance(req.GetCoinType(), req.GetAmount()); err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}
			ee.SaveAccount()
			res := pp.UpdateCreditRes{
				Result: pp.MakeResultWithCode(pp.ErrCode_Success),
			}

			return c.SendJSON(&res)
		}
		return c.Error(rlt)
	}
}
