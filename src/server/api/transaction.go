package api

import (
	"bytes"

	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/coin"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// InjectTx inject transaction.
func InjectTx(egn engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) {
		var rlt *pp.EmptyRes
		for {
			req := pp.InjectTxnReq{}
			if err := getRequest(c, &req); err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			tp, err := coin.TypeFromStr(req.GetCoinType())
			if err != nil {
				rlt = pp.MakeErrRes(err)
				break
			}

			// get coin gateway
			gateway, err := coin.GetGateway(tp)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			// decode tx string into structed tx.
			tx, err := gateway.DecodeRawTx(bytes.NewBuffer(req.GetTx()))
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			// inject tx.
			txid, err := gateway.InjectTx(tx)
			// txid, err := injectTx(tp, req.GetTx())
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			res := pp.InjectTxnRes{
				Result: pp.MakeResultWithCode(pp.ErrCode_Success),
				Txid:   pp.PtrString(txid),
			}
			reply(c, &res)
			return
		}
		c.JSON(rlt)
	}
}

// GetTx get transaction by id.
func GetTx(egn engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) {
		var rlt *pp.EmptyRes
		for {
			req := pp.GetTxReq{}
			if err := getRequest(c, &req); err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			tp, err := coin.TypeFromStr(req.GetCoinType())
			if err != nil {
				rlt = pp.MakeErrRes(err)
				break
			}

			gateway, err := coin.GetGateway(tp)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			tx, err := gateway.GetTx(req.GetTxid())
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			res := pp.GetTxRes{
				Result:   pp.MakeResultWithCode(pp.ErrCode_Success),
				CoinType: req.CoinType,
				Tx:       tx.ToPPTx(),
			}
			reply(c, &res)
			return
		}
		c.JSON(rlt)
	}
}

// func injectTx(tp coin.Type, tx []byte) (string, error) {
// 	switch tp {
// 	case coin.Bitcoin:
// 		btctx := bitcoin.Transaction{}
// 		if err := btctx.Deserialize(bytes.NewBuffer(tx)); err != nil {
// 			return "", err
// 		}
// 		return bitcoin.BroadcastTx(&btctx)
// 	case coin.Skycoin:
// 		sktx := skycoin.Transaction{}
// 		if err := sktx.Deserialize(tx); err != nil {
// 			return "", err
// 		}
// 		return skycoin.BroadcastTx(sktx)
// 	default:
// 		return "", errors.New("inject Txn failed, unknow coin type")
// 	}
// }
