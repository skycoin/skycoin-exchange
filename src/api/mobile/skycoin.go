package mobile

import (
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"

	"github.com/skycoin/skycoin-exchange/src/coin"
	skycoin "github.com/skycoin/skycoin-exchange/src/coin/skycoin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/sknet"
	"github.com/skycoin/skycoin/src/cipher"
)

type skyNode struct {
	NodeAddr string
}

func (sn skyNode) getOutputs(addrs []string) ([]*pp.SkyUtxo, error) {
	// get uxout of the address
	_, s := cipher.GenerateKeyPair()
	sknet.SetKey(s.Hex())

	req := pp.GetUtxoReq{
		CoinType:  pp.PtrString("skycoin"),
		Addresses: addrs,
	}
	res := pp.GetUtxoRes{}
	if err := sknet.EncryGet(sn.NodeAddr, "/auth/get/utxos", req, &res); err != nil {
		return nil, err
	}

	if !res.Result.GetSuccess() {
		return nil, fmt.Errorf("get utxos failed: %v", res.Result.GetReason())
	}

	return res.SkyUtxos, nil
}

func (sn skyNode) GetBalance(addrs []string) (uint64, error) {
	utxos, err := sn.getOutputs(addrs)
	if err != nil {
		return 0, err
	}

	var bal uint64
	for _, u := range utxos {
		bal += u.GetCoins()
	}

	return bal, nil
}

func (sn skyNode) ValidateAddr(address string) error {
	_, err := cipher.DecodeBase58Address(address)
	return err
}

func (sn skyNode) CreateRawTx(txIns []coin.TxIn, keys []cipher.SecKey, txOuts interface{}) (string, error) {
	tx := skycoin.Transaction{}
	for _, in := range txIns {
		tx.PushInput(cipher.MustSHA256FromHex(in.Txid))
	}

	s := reflect.ValueOf(txOuts)
	if s.Kind() != reflect.Slice {
		return "", errors.New("error tx out type")
	}
	outs := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		outs[i] = s.Index(i).Interface()
	}

	if len(outs) > 2 {
		return "", errors.New("out address more than 2")
	}

	for _, o := range outs {
		out := o.(skycoin.TxOut)
		if (out.Coins % 1e6) != 0 {
			return "", errors.New("skycoin coins must be multiple of 1e6")
		}
		tx.PushOutput(out.Address, out.Coins, out.Hours)
	}
	tx.SignInputs(keys)
	tx.UpdateHeader()

	d, err := tx.Serialize()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(d), nil
}

func (sn skyNode) BroadcastTx(rawtx string) (string, error) {
	gw := skycoin.Gateway{}
	return gw.InjectTx(rawtx)
}

func (sn skyNode) PrepareTx(addrs []string, toAddr string, amt uint64) ([]coin.TxIn, []string, interface{}, error) {
	outMap := make(map[string][]*pp.SkyUtxo)
	utxos, err := sn.getOutputs(addrs)
	if err != nil {
		return nil, nil, nil, err
	}

	for _, u := range utxos {
		outMap[u.GetAddress()] = append(outMap[u.GetAddress()], u)
	}

	allUtxos := []*pp.SkyUtxo{}
	var allBal uint64
	var allHours uint64
	for _, utxos := range outMap {
		bal, hour := func(utxos []*pp.SkyUtxo) (uint64, uint64) {
			var bal uint64
			var hour uint64
			for _, u := range utxos {
				if u.GetCoins() == 0 {
					continue
				}
				bal += u.GetCoins()
				hour += u.GetHours()
			}
			return bal, hour
		}(utxos)

		allUtxos = append(allUtxos, utxos...)

		allHours += hour
		allBal += bal
		if allBal >= amt {
			break
		}
	}

	if allBal >= amt {
		txIns := make([]coin.TxIn, len(allUtxos))
		inAddrs := make([]string, len(allUtxos))
		for i, u := range allUtxos {
			txIns[i] = coin.TxIn{
				Txid: u.GetHash(),
			}
			inAddrs[i] = u.GetAddress()
		}

		var txOut []skycoin.TxOut
		chgAmt := allBal - amt
		chgHours := allHours / 4
		chgAddr := addrs[0]
		if chgAmt > 0 {
			txOut = append(txOut,
				makeSkyTxOut(toAddr, amt, chgHours/2),
				makeSkyTxOut(chgAddr, chgAmt, chgHours/2))
		} else {
			txOut = append(txOut, makeSkyTxOut(toAddr, amt, chgHours/2))
		}
		return txIns, inAddrs, txOut, nil
	}

	return nil, nil, nil, errors.New("balance is not sufficient")
}

func makeSkyTxOut(addr string, coins uint64, hours uint64) skycoin.TxOut {
	out := skycoin.TxOut{}
	out.Address = cipher.MustDecodeBase58Address(addr)
	out.Coins = coins
	out.Hours = hours
	return out
}
