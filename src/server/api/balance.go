package api

import (
	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// GetBalance return balance of specific account.
func GetBalance(ee engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) {
		rlt := &pp.EmptyRes{}
		for {
			req := pp.GetBalanceReq{}
			if err := getRequest(c, &req); err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			// validate pubkey
			pubkey := req.GetPubkey()
			if err := validatePubkey(pubkey); err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongPubkey)
				break
			}

			a, err := ee.GetAccount(pubkey)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_NotExits)
				break
			}

			ct, err := coin.TypeFromStr(req.GetCoinType())
			if err != nil {
				rlt = pp.MakeErrRes(err)
				break
			}

			bal := a.GetBalance(ct)
			bres := pp.GetBalanceRes{
				Result:   pp.MakeResultWithCode(pp.ErrCode_Success),
				CoinType: req.CoinType,
				Balance:  &bal,
			}
			reply(c, bres)
			return
		}

		c.JSON(rlt)
	}
}
