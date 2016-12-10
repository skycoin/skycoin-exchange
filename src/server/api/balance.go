package api

import (
	"strings"

	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// GetBalance return balance of specific account.
func GetAccountBalance(ee engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) error {
		rlt := &pp.EmptyRes{}
		for {
			req := pp.GetAccountBalanceReq{}
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

			a, err := ee.GetAccount(pubkey)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_NotExits)
				break
			}

			bal := a.GetBalance(req.GetCoinType())
			bres := pp.GetAccountBalanceRes{
				Result:  pp.MakeResultWithCode(pp.ErrCode_Success),
				Balance: &pp.Balance{Amount: pp.PtrUint64(bal)},
			}
			return c.SendJSON(&bres)
		}
		return c.Error(rlt)
	}
}

// GetAddrBalance get balance of specific address.
func GetAddrBalance(ee engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) error {
		var rlt *pp.EmptyRes
		for {
			req := pp.GetAddrBalanceReq{}
			if err := c.BindJSON(&req); err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			coin, err := ee.GetCoin(req.GetCoinType())
			if err != nil {
				rlt = pp.MakeErrRes(err)
				logger.Error(err.Error())
				break
			}

			addrs := strings.Split(req.GetAddrs(), ",")
			b, err := coin.GetBalance(addrs)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}
			res := pp.GetAddrBalanceRes{
				Result:  pp.MakeResultWithCode(pp.ErrCode_Success),
				Balance: &b,
			}

			return c.SendJSON(&res)
		}
		return c.Error(rlt)
	}
}
