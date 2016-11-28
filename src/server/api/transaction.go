package api

import (
	"errors"

	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// InjectTx inject transaction.
func InjectTx(egn engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) error {
		var rlt *pp.EmptyRes
		for {
			req := pp.InjectTxnReq{}
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

			// get coin gateway
			gateway, err := coin.GetGateway(cp)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			// inject tx.
			txid, err := gateway.InjectTx(req.GetTx())
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			res := pp.InjectTxnRes{
				Result: pp.MakeResultWithCode(pp.ErrCode_Success),
				Txid:   pp.PtrString(txid),
			}
			return c.SendJSON(&res)
		}
		return c.Error(rlt)
	}
}

// GetTx get transaction by id.
func GetTx(egn engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) error {
		var rlt *pp.EmptyRes
		for {
			req := pp.GetTxReq{}
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

			gateway, err := coin.GetGateway(cp)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			// brief validate transaction id
			if !gateway.ValidateTxid(req.GetTxid()) {
				rlt = pp.MakeErrRes(errors.New("invalid transaction id"))
				break
			}

			tx, err := gateway.GetTx(req.GetTxid())
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			res := pp.GetTxRes{
				Result:   pp.MakeResultWithCode(pp.ErrCode_Success),
				CoinType: req.CoinType,
				Tx:       tx,
			}
			return c.SendJSON(&res)
		}
		return c.Error(rlt)
	}
}

// GetRawTx return rawtx of specifc tx.
func GetRawTx(egn engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) error {
		var rlt *pp.EmptyRes
		for {
			req := pp.GetRawTxReq{}
			if err := c.BindJSON(&req); err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			cp, err := coin.TypeFromStr(req.GetCoinType())
			if err != nil {
				rlt = pp.MakeErrRes(err)
				break
			}

			gateway, err := coin.GetGateway(cp)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			rawtx, err := gateway.GetRawTx(req.GetTxid())
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			res := pp.GetRawTxRes{
				Result:   pp.MakeResultWithCode(pp.ErrCode_Success),
				CoinType: req.CoinType,
				Rawtx:    pp.PtrString(rawtx),
			}
			return c.SendJSON(&res)
		}
		return c.Error(rlt)
	}
}
