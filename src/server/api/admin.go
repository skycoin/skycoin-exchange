package api

import (
	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// IsAdmin middleware for checking if the account is admin.
func IsAdmin(ee engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) {
		var rlt *pp.EmptyRes
		for {
			req := pp.UpdateCreditReq{}
			if err := getRequest(c, &req); err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			if !ee.IsAdmin(req.GetPubkey()) {
				logger.Error("not admin")
				rlt = pp.MakeErrResWithCode(pp.ErrCode_UnAuthorized)
				break
			}

			c.Next()
			return
		}
		c.JSON(rlt)
	}
}

// UpdateCredit update credit.
func UpdateCredit(ee engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) {
		var rlt *pp.EmptyRes
		for {
			req := pp.UpdateCreditReq{}
			if err := getRequest(c, &req); err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			// validate the dst pubkey.
			dstPubkey := req.GetDstPubkey()
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
			cp, err := coin.TypeFromStr(req.GetCoinType())
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			if err := a.SetBalance(cp, req.GetCredit()); err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}
			res := pp.UpdateCreditRes{
				Result: pp.MakeResultWithCode(pp.ErrCode_Success),
			}

			reply(c, res)
			return
		}
		c.JSON(rlt)
	}
}
