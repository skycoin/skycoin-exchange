package api

import (
	"github.com/skycoin/skycoin-exchange/src/coin"
	bitcoin "github.com/skycoin/skycoin-exchange/src/coin/bitcoin"
	skycoin "github.com/skycoin/skycoin-exchange/src/coin/skycoin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
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

			tp, err := coin.TypeFromStr(req.GetCoinType())
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				logger.Error("%s", err.Error())
				break
			}
			res, err := getUtxos(tp, req.GetAddresses())
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				logger.Error("%s", err.Error())
				break
			}
			res.Result = pp.MakeResultWithCode(pp.ErrCode_Success)
			reply(c, res)
			return
		}
		c.JSON(rlt)
	}
}

func getUtxos(tp coin.Type, addrs []string) (*pp.GetUtxoRes, error) {
	var res pp.GetUtxoRes
	switch tp {
	case coin.Bitcoin:
		utxos, err := bitcoin.GetUnspentOutputs(addrs)
		if err != nil {
			return nil, err
		}
		btcUxs := make([]*pp.BtcUtxo, len(utxos))
		for i, u := range utxos {
			btcUxs[i] = &pp.BtcUtxo{
				Address: pp.PtrString(u.GetAddress()),
				Txid:    pp.PtrString(u.GetTxid()),
				Vout:    pp.PtrUint32(u.GetVout()),
				Amount:  pp.PtrUint64(u.GetAmount()),
			}
		}
		res.BtcUtxos = btcUxs
	case coin.Skycoin:
		utxos, err := skycoin.GetUnspentOutputs(addrs)
		if err != nil {
			return nil, err
		}
		skyUxs := make([]*pp.SkyUtxo, len(utxos))
		for i, u := range utxos {
			skyUxs[i] = &pp.SkyUtxo{
				Hash:    pp.PtrString(u.GetHash()),
				SrcTx:   pp.PtrString(u.GetSrcTx()),
				Address: pp.PtrString(u.GetAddress()),
				Coins:   pp.PtrUint64(u.GetCoins()),
				Hours:   pp.PtrUint64(u.GetHours()),
			}
		}
		res.SkyUtxos = skyUxs
	}
	return &res, nil
}
