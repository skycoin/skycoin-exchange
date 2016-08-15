package api

import (
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin-exchange/src/sknet"
	"github.com/skycoin/skycoin/src/cipher"
)

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
			if _, err := cipher.PubKeyFromHex(breq.GetAccountId()); err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongAccountId)
				break
			}

			a, err := ee.GetAccount(breq.GetAccountId())
			if err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_NotExits)
				break
			}

			ct, err := wallet.CoinTypeFromStr(breq.GetCoinType())
			if err != nil {
				errRlt = pp.MakeErrRes(err)
				break
			}

			bal := a.GetBalance(ct)
			bres := pp.GetBalanceRes{
				Result:    pp.MakeResultWithCode(pp.ErrCode_Success),
				AccountId: breq.AccountId,
				CoinType:  breq.CoinType,
				Balance:   &bal,
			}
			reply(c, bres)
			return
		}

		c.JSON(errRlt)
	}
}
