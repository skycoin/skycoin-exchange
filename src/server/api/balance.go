package api

import (
	"strings"

	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// GetBalance return balance of specific account.
func GetAccountBalance(ee engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) {
		rlt := &pp.EmptyRes{}
		for {
			req := pp.GetAccountBalanceReq{}
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
			bres := pp.GetAccountBalanceRes{
				Result:  pp.MakeResultWithCode(pp.ErrCode_Success),
				Balance: &bal,
			}
			reply(c, bres)
			return
		}

		c.JSON(rlt)
	}
}

// GetAddrBalance get balance of specific address.
func GetAddrBalance(ee engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) {
		var rlt *pp.EmptyRes
		for {
			req := pp.GetAddrBalanceReq{}
			if err := c.BindJSON(&req); err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			cp, err := coin.TypeFromStr(req.GetCoinType())
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			gw, err := coin.GetGateway(cp)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			addrs := strings.Split(req.GetAddrs(), ",")
			b, err := gw.GetBalance(addrs)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			res := pp.GetAddrBalanceRes{
				Result:  pp.MakeResultWithCode(pp.ErrCode_Success),
				Balance: pp.PtrUint64(b),
			}

			c.JSON(&res)
			return
		}
		c.JSON(rlt)
	}
}
