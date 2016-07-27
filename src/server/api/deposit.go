package api

import (
	"github.com/gin-gonic/gin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/account"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
)

// GetNewAddress account create new address for depositing.
func GetNewAddress(ee engine.Exchange) gin.HandlerFunc {
	return func(c *gin.Context) {
		errRlt := &pp.EmptyRes{}
		for {
			dar := pp.GetDepositAddrReq{}
			err := getRequest(c, &dar)
			if err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			// convert to cipher.PubKey
			pubkey := pp.BytesToPubKey(dar.GetAccountId())
			if err := pubkey.Verify(); err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongAccountId)
				break
			}

			at, err := ee.GetAccount(account.AccountID(pubkey))
			if err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_NotExits)
				break
			}

			ct, err := wallet.ConvertCoinType(dar.GetCoinType())
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
				AccountId: dar.AccountId,
				CoinType:  dar.CoinType,
				Address:   &addr,
			}

			reply(c, ds)
			return
		}

		c.JSON(200, *errRlt)
	}
}
