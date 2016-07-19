package api

import (
	"bytes"
	"errors"
	"time"

	"github.com/btcsuite/btcd/wire"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/account"
	"github.com/skycoin/skycoin-exchange/src/server/coin_interface/bitcoin"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
)

var ChooseUtxoTmout = 1 * time.Second

// Withdraw api handler for generating withdraw transaction.
func Withdraw(ee engine.Exchange) gin.HandlerFunc {
	return func(c *gin.Context) {
		errRlt := &pp.EmptyRes{}
		for {
			wr := pp.WithdrawalReq{}
			err := getRequest(c, &wr)
			if err != nil {
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

			var success bool
			tx, usingUtxos, err := generateWithdrawlTx(ee, a, ct, wr.GetCoins(), wr.GetOutputAddress())
			if err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			defer func() {
				if !success {
					ee.PutUtxos(ct, usingUtxos)
				}
			}()

			switch ct {
			case wallet.Bitcoin:
				t := wire.MsgTx{}
				if err := t.Deserialize(bytes.NewBuffer(tx)); err != nil {
					glog.Error(err)
					errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
					break
				}

				newTxid, err := bitcoin_interface.BroadcastTx(&t)
				if err != nil {
					errRlt = pp.MakeErrResWithCode(pp.ErrCode_BroadcastTxFail)
					break
				}

				success = true
				a.DecreaseBalance(ct, wr.GetCoins())
				resp := pp.WithdrawalRes{
					AccountId: wr.AccountId,
					NewTxid:   &newTxid,
				}

				reply(c, resp)
				return
				// TODO:
				// case wallet.Skycoin:
			}
		}
		c.JSON(200, *errRlt)
	}
}

// generateWithdrawlTx create withdraw transaction.
// act is the user that want to withdraw coins, it's balance need to be checked.
// coinType specific which kind of coin the user want to withdraw.
// amount is the number of coins that want to withdraw.
// toAddr is the address that the coins will be sent to.
func generateWithdrawlTx(egn engine.Exchange, act account.Accounter, coinType wallet.CoinType, amount uint64, toAddr string) ([]byte, []bitcoin_interface.UtxoWithkey, error) {
	bal := act.GetBalance(coinType)
	fee := egn.GetFee()
	if bal < amount+fee {
		return []byte{}, []bitcoin_interface.UtxoWithkey{}, errors.New("balance is not sufficient")
	}

	utxos, err := egn.ChooseUtxos(coinType, amount, ChooseUtxoTmout)
	if err != nil {
		return []byte{}, []bitcoin_interface.UtxoWithkey{}, err
	}

	var totalAmounts uint64
	for _, u := range utxos {
		totalAmounts += u.GetAmount()
	}

	outAddrs := []bitcoin_interface.UtxoOut{}
	chgAmt := totalAmounts - fee - amount
	if chgAmt > 0 {
		// generate a change address
		chgAddr := egn.GetNewAddress(coinType)
		egn.AddWatchAddress(coinType, chgAddr)
		outAddrs = append(outAddrs, bitcoin_interface.UtxoOut{Addr: toAddr, Value: amount}, bitcoin_interface.UtxoOut{Addr: chgAddr, Value: chgAmt})
	} else {
		outAddrs = append(outAddrs, bitcoin_interface.UtxoOut{Addr: toAddr, Value: amount})
	}

	tx, err := bitcoin_interface.NewTransaction(utxos, outAddrs)
	if err != nil {
		return []byte{}, []bitcoin_interface.UtxoWithkey{}, err
	}

	return bitcoin_interface.DumpTxBytes(tx), utxos, nil
}
