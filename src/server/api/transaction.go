package api

import (
	"errors"

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

			// get coin gateway
			coin, err := egn.GetCoin(req.GetCoinType())
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			// inject tx.
			txid, err := coin.InjectTx(req.GetTx())
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

			coin, err := egn.GetCoin(req.GetCoinType())
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			// brief validate transaction id
			if !coin.ValidateTxid(req.GetTxid()) {
				rlt = pp.MakeErrRes(errors.New("invalid transaction id"))
				break
			}

			tx, err := coin.GetTx(req.GetTxid())
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

			coin, err := egn.GetCoin(req.GetCoinType())
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}
			rawtx, err := coin.GetRawTx(req.GetTxid())
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
