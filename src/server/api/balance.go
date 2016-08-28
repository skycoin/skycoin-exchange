package api

import (
	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/sknet"
	"github.com/skycoin/skycoin/src/cipher"
)

// GetBalance return balance of specific account.
func GetBalance(ee engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) {
		errRlt := &pp.EmptyRes{}
		for {
			breq := pp.GetBalanceReq{}
			if err := getRequest(c, &breq); err != nil {
				logger.Error(err.Error())
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			// convert to cipher.PubKey
			if _, err := cipher.PubKeyFromHex(breq.GetPubkey()); err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongAccountId)
				break
			}

			a, err := ee.GetAccount(breq.GetPubkey())
			if err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_NotExits)
				break
			}

			ct, err := coin.TypeFromStr(breq.GetCoinType())
			if err != nil {
				errRlt = pp.MakeErrRes(err)
				break
			}

			bal := a.GetBalance(ct)
			bres := pp.GetBalanceRes{
				Result:   pp.MakeResultWithCode(pp.ErrCode_Success),
				CoinType: breq.CoinType,
				Balance:  &bal,
			}
			reply(c, bres)
			return
		}

		c.JSON(errRlt)
	}
}
