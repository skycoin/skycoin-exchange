package api

import (
	"errors"
	"time"

	"github.com/btcsuite/btcd/wire"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/account"
	bitcoin "github.com/skycoin/skycoin-exchange/src/server/coin_interface/bitcoin"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

var ChooseUtxoTmout = 5 * time.Second

type BtcTxResult struct {
	Tx         *wire.MsgTx
	UsingUtxos []bitcoin.Utxo
	ChangeAddr string
}

// Withdraw api handler for generating withdraw transaction.
func Withdraw(ee engine.Exchange) gin.HandlerFunc {
	return func(c *gin.Context) {
		errRlt := &pp.EmptyRes{}
		for {
			rp := NewReqParams()

			wr := pp.WithdrawalReq{}
			if err := getRequest(c, &wr); err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			// convert to cipher.PubKey
			pubkey := pp.BytesToPubKey(wr.GetAccountId())
			if err := pubkey.Verify(); err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongAccountId)
				break
			}

			a, err := ee.GetAccount(account.AccountID(pubkey))
			if err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_NotExits)
				break
			}

			ct, err := wallet.ConvertCoinType(wr.GetCoinType())
			if err != nil {
				errRlt = pp.MakeErrRes(err)
				break
			}
			rp.Values["engine"] = ee
			rp.Values["pubkey"] = pubkey
			rp.Values["account"] = a
			rp.Values["cointype"] = ct
			rp.Values["amt"] = wr.GetCoins()
			rp.Values["toAddr"] = wr.GetOutputAddress()

			resp, rlt := withdrawlWork(c, rp)
			if rlt != nil {
				errRlt = rlt
				break
			}

			reply(c, *resp)
			return
		}
		c.JSON(200, *errRlt)
	}
}

func withdrawlWork(c *gin.Context, rp *ReqParams) (*pp.WithdrawalRes, *pp.EmptyRes) {
	ee := rp.Values["engine"].(engine.Exchange)
	acnt := rp.Values["account"].(account.Accounter)
	amt := rp.Values["amt"].(uint64)
	ct := rp.Values["cointype"].(wallet.CoinType)
	toAddr := rp.Values["toAddr"].(string)

	switch ct {
	case wallet.Bitcoin:
		var success bool
		var btcTxRlt *BtcTxResult
		var err error
		if err := acnt.DecreaseBalance(ct, amt+ee.GetFee()); err != nil {
			return nil, pp.MakeErrRes(err)
		}
		defer func() {
			if !success {
				go func() {
					if btcTxRlt != nil {
						ee.BtcPutUtxos(btcTxRlt.UsingUtxos)
					}
					acnt.IncreaseBalance(ct, amt+ee.GetFee())
				}()
			} else {
				//TODO: handle the saving failure.
				ee.SaveAccount()
			}
		}()

		btcTxRlt, err = createBtcWithdrawTx(ee, amt, toAddr)
		if err != nil {
			return nil, pp.MakeErrRes(errors.New("failed to create withdrawal tx"))
		}

		newTxid, err := bitcoin.BroadcastTx(btcTxRlt.Tx)
		// newTxid, err := "123", errors.New("broadcast tx not support yet")
		// newTxid, err := "123", nil
		if err != nil {
			glog.Error(err)
			// errRlt = pp.MakeErrResWithCode(pp.ErrCode_BroadcastTxFail)
			return nil, pp.MakeErrResWithCode(pp.ErrCode_BroadcastTxFail)
		}

		success = true
		if btcTxRlt.ChangeAddr != "" {
			glog.Info("change address:", btcTxRlt.ChangeAddr)
			ee.AddWatchAddress(ct, btcTxRlt.ChangeAddr)
		}

		pk := cipher.PubKey(acnt.GetID())
		resp := pp.WithdrawalRes{
			AccountId: pk[:],
			NewTxid:   &newTxid,
		}
		return &resp, nil
	case wallet.Skycoin:
		return nil, pp.MakeErrRes(errors.New("skycoin withdrawal not support yet"))
	default:
		return nil, pp.MakeErrRes(errors.New("unknow coin type"))
	}
}

// createBtcWithdrawTx create withdraw transaction.
// amount is the number of coins that want to withdraw.
// toAddr is the address that the coins will be sent to.
func createBtcWithdrawTx(egn engine.Exchange, amount uint64, toAddr string) (*BtcTxResult, error) {
	utxos, err := egn.BtcChooseUtxos(amount+egn.GetFee(), 5*time.Second)
	if err != nil {
		return nil, err
	}

	for _, u := range utxos {
		glog.Info("using utxos:", u.GetTxid(), " ", u.GetVout())
	}

	var success bool
	defer func() {
		if !success {
			go func() { egn.BtcPutUtxos(utxos) }()
		}
	}()

	var totalAmounts uint64
	for _, u := range utxos {
		totalAmounts += u.GetAmount()
	}
	fee := egn.GetFee()
	outAddrs := []bitcoin.UtxoOut{}
	chgAmt := totalAmounts - fee - amount
	chgAddr := ""
	if chgAmt > 0 {
		// generate a change address
		chgAddr = egn.GetNewAddress(wallet.Bitcoin)
		// egn.AddWatchAddress(coinType, chgAddr)
		outAddrs = append(outAddrs,
			bitcoin.UtxoOut{Addr: toAddr, Value: amount},
			bitcoin.UtxoOut{Addr: chgAddr, Value: chgAmt})
	} else {
		outAddrs = append(outAddrs, bitcoin.UtxoOut{Addr: toAddr, Value: amount})
	}
	// change utxo to UtxoWithkey
	utxoKeys, err := makeBtcUtxoWithkeys(utxos, egn)
	if err != nil {
		return nil, err
	}

	glog.Info("creating transaction...")
	tx, err := bitcoin.NewTransaction(utxoKeys, outAddrs)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	success = true
	rlt := BtcTxResult{
		Tx:         tx,
		UsingUtxos: utxos[:],
		ChangeAddr: chgAddr,
	}
	return &rlt, nil
}

func makeBtcUtxoWithkeys(utxos []bitcoin.Utxo, egn engine.Exchange) ([]bitcoin.UtxoWithkey, error) {
	utxoks := make([]bitcoin.UtxoWithkey, len(utxos))
	for i, u := range utxos {
		key, err := egn.GetPrivKey(wallet.Bitcoin, u.GetAddress())
		if err != nil {
			return []bitcoin.UtxoWithkey{}, err
		}
		utxoks[i] = bitcoin.NewUtxoWithKey(u, key)
	}
	return utxoks, nil
}
