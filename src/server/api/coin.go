package api

import (
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/net"
)

func GetCoins(egn engine.Exchange) net.HandlerFunc {
	return func(c *net.Context) {
		coins := pp.CoinsRes{
			Result: pp.MakeResultWithCode(pp.ErrCode_Success),
			Coins:  egn.GetSupportCoins(),
		}
		reply(c, &coins)
	}
}
