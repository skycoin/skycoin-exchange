package api

import (
	"errors"
	"time"

	"github.com/btcsuite/btcd/wire"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/account"
	bitcoin "github.com/skycoin/skycoin-exchange/src/server/coin_interface/bitcoin"
	skycoin "github.com/skycoin/skycoin-exchange/src/server/coin_interface/skycoin"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/server/net"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

var ChooseUtxoTmout = 5 * time.Second

type BtcTxResult struct {
	Tx         *wire.MsgTx
	UsingUtxos []bitcoin.Utxo
	ChangeAddr string
}

type SkyTxResult struct {
	Tx         *skycoin.Transaction
	UsingUtxos []skycoin.Utxo
	ChangeAddr string
}

// Withdraw api handler for generating withdraw transaction.
func Withdraw(ee engine.Exchange) net.HandlerFunc {
	return func(c *net.Context) {
		errRlt := &pp.EmptyRes{}
		for {
			rp := NewReqParams()

			wr := pp.WithdrawalReq{}
			if err := getRequest(c, &wr); err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			if _, err := cipher.PubKeyFromHex(wr.GetAccountId()); err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongAccountId)
				break
			}

			a, err := ee.GetAccount(wr.GetAccountId())
			if err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_NotExits)
				break
			}

			ct, err := wallet.CoinTypeFromStr(wr.GetCoinType())
			if err != nil {
				errRlt = pp.MakeErrRes(err)
				break
			}
			rp.Values["engine"] = ee
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
		c.JSON(errRlt)
	}
}

func withdrawlWork(c *net.Context, rp *ReqParams) (*pp.WithdrawalRes, *pp.EmptyRes) {
	ct := rp.Values["cointype"].(wallet.CoinType)
	switch ct {
	case wallet.Bitcoin:
		return btcWithdraw(rp)
	case wallet.Skycoin:
		return skyWithdrawl(rp)
	default:
		return nil, pp.MakeErrRes(errors.New("unknow coin type"))
	}
}

func btcWithdraw(rp *ReqParams) (*pp.WithdrawalRes, *pp.EmptyRes) {
	ee := rp.Values["engine"].(engine.Exchange)
	acnt := rp.Values["account"].(account.Accounter)
	amt := rp.Values["amt"].(uint64)
	ct := rp.Values["cointype"].(wallet.CoinType)
	toAddr := rp.Values["toAddr"].(string)
	// verify the toAddr
	if _, err := cipher.BitcoinDecodeBase58Address(toAddr); err != nil {
		return nil, pp.MakeErrRes(errors.New("invalid bitcoin address"))
	}
	var success bool
	var btcTxRlt *BtcTxResult
	var err error
	if err := acnt.DecreaseBalance(ct, amt+ee.GetBtcFee()); err != nil {
		return nil, pp.MakeErrRes(err)
	}
	defer func() {
		if !success {
			go func() {
				if btcTxRlt != nil {
					ee.PutUtxos(wallet.Bitcoin, btcTxRlt.UsingUtxos)
				}
				acnt.IncreaseBalance(ct, amt+ee.GetBtcFee())
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
	if err != nil {
		logger.Error(err.Error())
		return nil, pp.MakeErrResWithCode(pp.ErrCode_BroadcastTxFail)
	}

	success = true
	if btcTxRlt.ChangeAddr != "" {
		logger.Debug("change address:%s", btcTxRlt.ChangeAddr)
		ee.WatchAddress(ct, btcTxRlt.ChangeAddr)
	}

	id := acnt.GetID()
	resp := pp.WithdrawalRes{
		Result:    pp.MakeResultWithCode(pp.ErrCode_Success),
		AccountId: &id,
		NewTxid:   &newTxid,
	}
	return &resp, nil
}

func skyWithdrawl(rp *ReqParams) (*pp.WithdrawalRes, *pp.EmptyRes) {
	ee := rp.Values["engine"].(engine.Exchange)
	acnt := rp.Values["account"].(account.Accounter)
	amt := rp.Values["amt"].(uint64)
	ct := rp.Values["cointype"].(wallet.CoinType)
	toAddr := rp.Values["toAddr"].(string)

	if err := skycoin.VerifyAmount(amt); err != nil {
		return nil, pp.MakeErrRes(err)
	}

	// verify the toAddr
	if _, err := cipher.DecodeBase58Address(toAddr); err != nil {
		return nil, pp.MakeErrRes(errors.New("invalid skycoin address"))
	}

	var success bool
	var skyTxRlt *SkyTxResult
	var err error
	if err := acnt.DecreaseBalance(ct, amt); err != nil {
		return nil, pp.MakeErrRes(err)
	}
	defer func() {
		if !success {
			go func() {
				if skyTxRlt != nil {
					ee.PutUtxos(wallet.Skycoin, skyTxRlt.UsingUtxos)
				}
				acnt.IncreaseBalance(ct, amt)
			}()
		} else {
			//TODO: handle the saving failure.
			ee.SaveAccount()
		}
	}()

	skyTxRlt, err = createSkyWithdrawTx(ee, amt, toAddr)
	if err != nil {
		return nil, pp.MakeErrRes(errors.New("failed to create withdrawal tx"))
	}

	newTxid, err := skycoin.BroadcastTx(*skyTxRlt.Tx)
	if err != nil {
		logger.Error(err.Error())
		return nil, pp.MakeErrResWithCode(pp.ErrCode_BroadcastTxFail)
	}

	success = true
	if skyTxRlt.ChangeAddr != "" {
		logger.Debug("change address:%s", skyTxRlt.ChangeAddr)
		ee.WatchAddress(ct, skyTxRlt.ChangeAddr)
	}

	id := acnt.GetID()
	resp := pp.WithdrawalRes{
		Result:    pp.MakeResultWithCode(pp.ErrCode_Success),
		AccountId: &id,
		NewTxid:   &newTxid,
	}
	return &resp, nil
}

// createBtcWithdrawTx create withdraw transaction.
// amount is the number of coins that want to withdraw.
// toAddr is the address that the coins will be sent to.
func createBtcWithdrawTx(egn engine.Exchange, amount uint64, toAddr string) (*BtcTxResult, error) {
	uxs, err := egn.ChooseUtxos(wallet.Bitcoin, amount+egn.GetBtcFee(), 5*time.Second)
	if err != nil {
		return nil, err
	}
	utxos := uxs.([]bitcoin.Utxo)

	for _, u := range utxos {
		logger.Debug("using utxos: txid:%s vout:%d", u.GetTxid(), u.GetVout())
	}

	var success bool
	defer func() {
		if !success {
			go func() { egn.PutUtxos(wallet.Bitcoin, utxos) }()
		}
	}()

	var totalAmounts uint64
	for _, u := range utxos {
		totalAmounts += u.GetAmount()
	}
	fee := egn.GetBtcFee()
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

	logger.Debug("creating transaction...")
	tx, err := bitcoin.NewTransaction(utxoKeys, outAddrs)
	if err != nil {
		logger.Error(err.Error())
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

func createSkyWithdrawTx(egn engine.Exchange, amount uint64, toAddr string) (*SkyTxResult, error) {
	uxs, err := egn.ChooseUtxos(wallet.Skycoin, amount, 5*time.Second)
	if err != nil {
		return nil, err
	}
	utxos := uxs.([]skycoin.Utxo)

	for _, u := range utxos {
		logger.Debug("using skycoin utxos:%s", u.GetHash())
	}

	var success bool
	defer func() {
		if !success {
			go func() { egn.PutUtxos(wallet.Skycoin, utxos) }()
		}
	}()

	var totalAmounts uint64
	var totalHours uint64
	for _, u := range utxos {
		totalAmounts += u.GetCoins()
		totalHours += u.GetHours()
	}

	outAddrs := []skycoin.UtxoOut{}
	chgAmt := totalAmounts - amount
	chgHours := totalHours / 4
	chgAddr := ""
	if chgAmt > 0 {
		// generate a change address
		chgAddr = egn.GetNewAddress(wallet.Skycoin)
		outAddrs = append(outAddrs,
			skycoin.MakeUtxoOutput(toAddr, amount, chgHours/2),
			skycoin.MakeUtxoOutput(chgAddr, chgAmt, chgHours/2))
	} else {
		outAddrs = append(outAddrs, skycoin.MakeUtxoOutput(toAddr, amount, chgHours/2))
	}

	keys := make([]cipher.SecKey, len(utxos))
	for i, u := range utxos {
		k, err := egn.GetAddrPrivKey(wallet.Skycoin, u.GetAddress())
		if err != nil {
			panic(err)
		}
		keys[i] = cipher.MustSecKeyFromHex(k)
	}

	logger.Debug("creating skycoin transaction...")
	tx := skycoin.NewTransaction(utxos, keys, outAddrs)
	if err := tx.Verify(); err != nil {
		return nil, err
	}

	success = true
	rlt := SkyTxResult{
		Tx:         tx,
		UsingUtxos: utxos[:],
		ChangeAddr: chgAddr,
	}
	return &rlt, nil
}

func makeBtcUtxoWithkeys(utxos []bitcoin.Utxo, egn engine.Exchange) ([]bitcoin.UtxoWithkey, error) {
	utxoks := make([]bitcoin.UtxoWithkey, len(utxos))
	for i, u := range utxos {
		key, err := egn.GetAddrPrivKey(wallet.Bitcoin, u.GetAddress())
		if err != nil {
			return []bitcoin.UtxoWithkey{}, err
		}
		utxoks[i] = bitcoin.NewUtxoWithKey(u, key)
	}
	return utxoks, nil
}
