package api

import (
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// GetCoins get supported coins.
func GetCoins(egn engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) error {
		coins := pp.CoinsRes{
			Result: pp.MakeResultWithCode(pp.ErrCode_Success),
			Coins:  egn.GetSupportCoins(),
		}
		return c.SendJSON(&coins)
	}
}
