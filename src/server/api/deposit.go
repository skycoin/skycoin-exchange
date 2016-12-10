package api

import (
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// GetNewAddress account create new address for depositing.
func GetNewAddress(ee engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) error {
		rlt := &pp.EmptyRes{}
		for {
			req := pp.GetDepositAddrReq{}
			if err := c.BindJSON(&req); err != nil {
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

			ct := req.GetCoinType()
			// get the new address for depositing
			addr := ee.GetNewAddress(ct)

			// add the new address to engin for watching it's utxos.
			at.AddDepositAddress(ct, addr)
			ee.WatchAddress(ct, addr)

			ds := pp.GetDepositAddrRes{
				Result:   pp.MakeResultWithCode(pp.ErrCode_Success),
				CoinType: req.CoinType,
				Address:  &addr,
			}

			return c.SendJSON(&ds)
		}

		return c.Error(rlt)
	}
}
