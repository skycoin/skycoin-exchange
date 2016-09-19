package api

import (
	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// GetNewAddress account create new address for depositing.
func GetNewAddress(ee engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) {
		rlt := &pp.EmptyRes{}
		for {
			req := pp.GetDepositAddrReq{}
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

			at, err := ee.GetAccount(pubkey)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_NotExits)
				break
			}

			cp, err := coin.TypeFromStr(req.GetCoinType())
			if err != nil {
				rlt = pp.MakeErrRes(err)
				break
			}

			// get the new address for depositing
			addr := ee.GetNewAddress(cp)

			// add the new address to engin for watching it's utxos.
			at.AddDepositAddress(cp, addr)
			ee.WatchAddress(cp, addr)

			ds := pp.GetDepositAddrRes{
				Result:   pp.MakeResultWithCode(pp.ErrCode_Success),
				CoinType: req.CoinType,
				Address:  &addr,
			}

			reply(c, ds)
			return
		}

		c.JSON(rlt)
	}
}
