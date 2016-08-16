package api

import (
	"bytes"
	"errors"

	"github.com/skycoin/skycoin-exchange/src/pp"
	bitcoin "github.com/skycoin/skycoin-exchange/src/server/coin_interface/bitcoin"
	skycoin "github.com/skycoin/skycoin-exchange/src/server/coin_interface/skycoin"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
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
			tp, err := wallet.CoinTypeFromStr(req.GetCoinType())
			if err != nil {
				rlt = pp.MakeErrRes(err)
				break
			}

			txid, err := injectTx(tp, req.GetTx())
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

// GetRawTxn get transaction by id or hex.
func GetRawTx(egn engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) {

	}
}

func injectTx(tp wallet.CoinType, tx []byte) (string, error) {
	switch tp {
	case wallet.Bitcoin:
		btctx := bitcoin.Transaction{}
		if err := btctx.Deserialize(bytes.NewBuffer(tx)); err != nil {
			return "", err
		}
		return bitcoin.BroadcastTx(&btctx)
	case wallet.Skycoin:
		sktx := skycoin.Transaction{}
		if err := sktx.Deserialize(tx); err != nil {
			return "", err
		}
		return skycoin.BroadcastTx(sktx)
	default:
		return "", errors.New("inject Txn failed, unknow coin type")
	}
}
