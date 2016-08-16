package api

import (
	"encoding/json"

	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/coin_interface/bitcoin"
	"github.com/skycoin/skycoin-exchange/src/server/coin_interface/skycoin"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

func GetUtxos(egn engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) {
		var req pp.GetUtxoReq
		var rlt *pp.EmptyRes
		for {
			if err := getRequest(c, &req); err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				logger.Error("%s", err.Error())
				break
			}

			tp, err := wallet.CoinTypeFromStr(req.GetCoinType())
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				logger.Error("%s", err.Error())
				break
			}

			var utxos interface{}
			switch tp {
			case wallet.Bitcoin:
				utxos, err = bitcoin_interface.GetUnspentOutputs(req.GetAddresses())
			case wallet.Skycoin:
				utxos, err = skycoin_interface.GetUnspentOutputs(req.GetAddresses())
			}

			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			v, err := utxos2Str(utxos)
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			res := pp.GetUtxoRes{
				CoinType: req.CoinType,
				Utxos:    pp.PtrString(v),
			}
			reply(c, &res)
		}
		c.JSON(rlt)
	}
}

func utxos2Str(utxos interface{}) (string, error) {
	d, err := json.Marshal(utxos)
	if err != nil {
		return "", err
	}
	return string(d), nil
}
