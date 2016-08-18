package api

import (
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/coin"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/sknet"
	"github.com/skycoin/skycoin/src/cipher"
)

// GetNewAddress account create new address for depositing.
func GetNewAddress(ee engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) {
		errRlt := &pp.EmptyRes{}
		for {
			dar := pp.GetDepositAddrReq{}
			err := getRequest(c, &dar)
			if err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			// convert to cipher.PubKey
			if _, err := cipher.PubKeyFromHex(dar.GetAccountId()); err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongAccountId)
				break
			}

			at, err := ee.GetAccount(dar.GetAccountId())
			if err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_NotExits)
				break
			}

			ct, err := coin.TypeFromStr(dar.GetCoinType())
			if err != nil {
				errRlt = pp.MakeErrRes(err)
				break
			}

			// get the new address for depositing
			addr := ee.GetNewAddress(ct)

			// add the new address to engin for watching it's utxos.
			at.AddDepositAddress(ct, addr)
			ee.WatchAddress(ct, addr)

			ds := pp.GetDepositAddrRes{
				Result:    pp.MakeResultWithCode(pp.ErrCode_Success),
				AccountId: dar.AccountId,
				CoinType:  dar.CoinType,
				Address:   &addr,
			}

			reply(c, ds)
			return
		}

		c.JSON(errRlt)
	}
}
