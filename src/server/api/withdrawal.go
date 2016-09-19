package api

import (
	"encoding/hex"
	"errors"
	"time"

	"github.com/skycoin/skycoin-exchange/src/coin"
	bitcoin "github.com/skycoin/skycoin-exchange/src/coin/bitcoin"
	skycoin "github.com/skycoin/skycoin-exchange/src/coin/skycoin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/account"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/sknet"
	"github.com/skycoin/skycoin/src/cipher"
)

// ChooseUtxoTmout max time that will be allowed in choosing sufficient utxos.
var ChooseUtxoTmout = 5 * time.Second

// BtcTxResult be used in creating bitcoin withdraw transaction.
type BtcTxResult struct {
	Tx         *bitcoin.Transaction
	UsingUtxos []bitcoin.Utxo
	ChangeAddr string
}

// SkyTxResult be used in creating skycoin withdraw transaction.
type SkyTxResult struct {
	Tx         *skycoin.Transaction
	UsingUtxos []skycoin.Utxo
	ChangeAddr string
}

// Withdraw api handler for generating withdraw transaction.
func Withdraw(ee engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) {
		errRlt := &pp.EmptyRes{}
		for {
			rp := NewReqParams()

			req := pp.WithdrawalReq{}
			if err := getRequest(c, &req); err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			// validate pubkey
			pubkey := req.GetPubkey()
			if err := validatePubkey(pubkey); err != nil {
				logger.Error(err.Error())
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongPubkey)
				break
			}

			if _, err := cipher.PubKeyFromHex(req.GetPubkey()); err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongPubkey)
				break
			}

			a, err := ee.GetAccount(req.GetPubkey())
			if err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_NotExits)
				break
			}

			ct, err := coin.TypeFromStr(req.GetCoinType())
			if err != nil {
				errRlt = pp.MakeErrRes(err)
				break
			}
			rp.Values["engine"] = ee
			rp.Values["account"] = a
			rp.Values["cointype"] = ct
			rp.Values["amt"] = req.GetCoins()
			rp.Values["toAddr"] = req.GetOutputAddress()

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

func withdrawlWork(c *sknet.Context, rp *ReqParams) (*pp.WithdrawalRes, *pp.EmptyRes) {
	ct := rp.Values["cointype"].(coin.Type)
	switch ct {
	case coin.Bitcoin:
		return btcWithdraw(rp)
	case coin.Skycoin:
		return skyWithdrawl(rp)
	default:
		return nil, pp.MakeErrRes(errors.New("unknow coin type"))
	}
}

func btcWithdraw(rp *ReqParams) (*pp.WithdrawalRes, *pp.EmptyRes) {
	ee := rp.Values["engine"].(engine.Exchange)
	acnt := rp.Values["account"].(account.Accounter)
	amt := rp.Values["amt"].(uint64)
	ct := rp.Values["cointype"].(coin.Type)
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
					ee.PutUtxos(coin.Bitcoin, btcTxRlt.UsingUtxos)
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

	rawtx, err := btcTxRlt.Tx.Serialize()
	if err != nil {
		return nil, pp.MakeErrRes(errors.New("tx serialize failed"))
	}

	newTxid, err := bitcoin.BroadcastTx(hex.EncodeToString(rawtx))
	if err != nil {
		logger.Error(err.Error())
		return nil, pp.MakeErrResWithCode(pp.ErrCode_BroadcastTxFail)
	}

	success = true
	if btcTxRlt.ChangeAddr != "" {
		logger.Debug("change address:%s", btcTxRlt.ChangeAddr)
		ee.WatchAddress(ct, btcTxRlt.ChangeAddr)
	}

	resp := pp.WithdrawalRes{
		Result:  pp.MakeResultWithCode(pp.ErrCode_Success),
		NewTxid: &newTxid,
	}
	return &resp, nil
}

func skyWithdrawl(rp *ReqParams) (*pp.WithdrawalRes, *pp.EmptyRes) {
	ee := rp.Values["engine"].(engine.Exchange)
	acnt := rp.Values["account"].(account.Accounter)
	amt := rp.Values["amt"].(uint64)
	ct := rp.Values["cointype"].(coin.Type)
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
					ee.PutUtxos(coin.Skycoin, skyTxRlt.UsingUtxos)
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
	rawtx, err := skyTxRlt.Tx.Serialize()
	if err != nil {
		return nil, pp.MakeErrRes(errors.New("skycoin tx serialize failed"))
	}

	newTxid, err := skycoin.BroadcastTx(hex.EncodeToString(rawtx))
	if err != nil {
		logger.Error(err.Error())
		return nil, pp.MakeErrResWithCode(pp.ErrCode_BroadcastTxFail)
	}

	success = true
	if skyTxRlt.ChangeAddr != "" {
		logger.Debug("change address:%s", skyTxRlt.ChangeAddr)
		ee.WatchAddress(ct, skyTxRlt.ChangeAddr)
	}

	resp := pp.WithdrawalRes{
		Result:  pp.MakeResultWithCode(pp.ErrCode_Success),
		NewTxid: &newTxid,
	}
	return &resp, nil
}

// createBtcWithdrawTx create withdraw transaction.
// amount is the number of coins that want to withdraw.
// toAddr is the address that the coins will be sent to.
func createBtcWithdrawTx(egn engine.Exchange, amount uint64, toAddr string) (*BtcTxResult, error) {
	uxs, err := egn.ChooseUtxos(coin.Bitcoin, amount+egn.GetBtcFee(), ChooseUtxoTmout)
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
			go func() { egn.PutUtxos(coin.Bitcoin, utxos) }()
		}
	}()

	var totalAmounts uint64
	for _, u := range utxos {
		totalAmounts += u.GetAmount()
	}
	fee := egn.GetBtcFee()
	outAddrs := []bitcoin.TxOut{}
	chgAmt := totalAmounts - fee - amount
	chgAddr := ""
	if chgAmt > 0 {
		// generate a change address
		chgAddr = egn.GetNewAddress(coin.Bitcoin)
		// egn.AddWatchAddress(coinType, chgAddr)
		outAddrs = append(outAddrs,
			bitcoin.TxOut{Addr: toAddr, Value: amount},
			bitcoin.TxOut{Addr: chgAddr, Value: chgAmt})
	} else {
		outAddrs = append(outAddrs, bitcoin.TxOut{Addr: toAddr, Value: amount})
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
	uxs, err := egn.ChooseUtxos(coin.Skycoin, amount, ChooseUtxoTmout)
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
			go func() { egn.PutUtxos(coin.Skycoin, utxos) }()
		}
	}()

	var totalAmounts uint64
	var totalHours uint64
	for _, u := range utxos {
		totalAmounts += u.GetCoins()
		totalHours += u.GetHours()
	}

	outAddrs := []skycoin.TxOut{}
	chgAmt := totalAmounts - amount
	chgHours := totalHours / 4
	chgAddr := ""
	if chgAmt > 0 {
		// generate a change address
		chgAddr = egn.GetNewAddress(coin.Skycoin)
		outAddrs = append(outAddrs,
			skycoin.MakeUtxoOutput(toAddr, amount, chgHours/2),
			skycoin.MakeUtxoOutput(chgAddr, chgAmt, chgHours/2))
	} else {
		outAddrs = append(outAddrs, skycoin.MakeUtxoOutput(toAddr, amount, chgHours/2))
	}

	keys := make([]cipher.SecKey, len(utxos))
	for i, u := range utxos {
		k, err := egn.GetAddrPrivKey(coin.Skycoin, u.GetAddress())
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
		key, err := egn.GetAddrPrivKey(coin.Bitcoin, u.GetAddress())
		if err != nil {
			return []bitcoin.UtxoWithkey{}, err
		}
		utxoks[i] = bitcoin.NewUtxoWithKey(u, key)
	}
	return utxoks, nil
}
