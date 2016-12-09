package api

import (
	"github.com/skycoin/skycoin-exchange/src/coin"
	bitcoin "github.com/skycoin/skycoin-exchange/src/coin/bitcoin"
	skycoin "github.com/skycoin/skycoin-exchange/src/coin/skycoin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/sknet"
)

// GetUtxos get unspent output of specific address.
func GetUtxos(egn engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) error {
		var req pp.GetUtxoReq
		var rlt *pp.EmptyRes
		for {
			if err := c.BindJSON(&req); err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				logger.Error(err.Error())
				break
			}

			cp, err := coin.TypeFromStr(req.GetCoinType())
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				logger.Error(err.Error())
				break
			}
			res, err := getUtxos(cp, req.GetAddresses())
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				logger.Error(err.Error())
				break
			}
			res.Result = pp.MakeResultWithCode(pp.ErrCode_Success)
			return c.SendJSON(&res)
		}
		return c.Error(rlt)
	}
}

// GetOutput  gets skycoin output by output hash id.
func GetOutput(egn engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) error {
		var rlt *pp.EmptyRes
		for {
			req := pp.GetOutputReq{}
			if err := c.BindJSON(&req); err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			output, err := skycoin.GetOutput(req.GetHash())
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			res := pp.GetOutputRes{
				Result: pp.MakeResultWithCode(pp.ErrCode_Success),
				Output: output,
			}
			return c.SendJSON(&res)
		}
		return c.Error(rlt)
	}
}

func getUtxos(nodeAddr string, cp coin.Type, addrs []string) (*pp.GetUtxoRes, error) {
	var res pp.GetUtxoRes
	switch cp {
	case coin.Bitcoin:
		utxos, err := bitcoin.GetUnspentOutputs(nodeAddr, addrs)
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
