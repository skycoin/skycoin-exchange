package api

import (
	"encoding/hex"
	"errors"
	"fmt"
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

// ChooseUtxoTm max time that will be allowed in choosing sufficient utxos.
var ChooseUtxoTm = 5 * time.Second

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

func getWithdrawReqParams(c *sknet.Context, ee engine.Exchange) (*ReqParams, error) {
	rp := NewReqParams()
	req := pp.WithdrawalReq{}
	if err := c.BindJSON(&req); err != nil {
		return nil, err
	}

	// validate pubkey
	pubkey := req.GetPubkey()
	if err := validatePubkey(pubkey); err != nil {
		return nil, err
	}

	a, err := ee.GetAccount(pubkey)
	if err != nil {
		return nil, err
	}

	ct, err := coin.TypeFromStr(req.GetCoinType())
	if err != nil {
		return nil, err
	}

	rp.Values["account"] = a
	rp.Values["cointype"] = ct
	rp.Values["amt"] = req.GetCoins()
	rp.Values["outAddr"] = req.GetOutputAddress()
	return rp, nil
}

// Withdraw api for handlering withdraw process.
func Withdraw(ee engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) error {
		rlt := &pp.EmptyRes{}
		for {
			reqParam, err := getWithdrawReqParams(c, ee)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			cp := reqParam.Values["cointype"].(coin.Type)
			a := reqParam.Values["account"].(account.Accounter)
			amt := reqParam.Values["amt"].(uint64)
			outAddr := reqParam.Values["outAddr"].(string)

			// get handler for creating txIns and txOuts base on the coin type.
			createTxInOut, err := getTxInOutHandler(cp)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			// create txIns and txOuts.
			inOutSet, err := createTxInOut(ee, a, amt, outAddr)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			var success bool
			defer func() {
				if !success {
					// if not success, invoke the teardown, for putting back utxos, and reset balance.
					inOutSet.Teardown()
				}
			}()

			// get coin gateway.
			gw, err := coin.GetGateway(cp)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			// create raw tx
			rawtx, err := gw.CreateRawTx(inOutSet.TxIns, inOutSet.TxOuts)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			// sign the tx
			rawtx, err = gw.SignRawTx(rawtx, getAddrPrivKey(ee, cp))
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			// inject the transaction.
			txid, err := gw.InjectTx(rawtx)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			success = true
			resp := pp.WithdrawalRes{
				Result:  pp.MakeResultWithCode(pp.ErrCode_Success),
				NewTxid: &txid,
			}
			return c.SendJSON(&resp)
		}
		return c.Error(rlt)
	}
}

func getAddrPrivKey(ee engine.Exchange, cp coin.Type) coin.GetPrivKey {
	return func(addr string) (string, error) {
		return ee.GetAddrPrivKey(cp, addr)
	}
}

// txInOutHandler used to generate TxIns and txOuts.
type txInOutHandler func(ee engine.Exchange, a account.Accounter, amount uint64, outAddr string) (*txInOutResult, error)

// global txInOut handlers, if new coin type need to be supported, register here.
var txInOutHandlers = map[coin.Type]txInOutHandler{
	coin.Bitcoin: createBtcTxInOut,
	coin.Skycoin: createSkyTxInOut,
}

func getTxInOutHandler(cp coin.Type) (txInOutHandler, error) {
	if hd, ok := txInOutHandlers[cp]; ok {
		return hd, nil
	}
	return nil, fmt.Errorf("%s tx in handler not found", cp)
}

type txInOutResult struct {
	TxIns    []coin.TxIn // transaction in values.
	TxOuts   interface{} // transaction out values, must be a slice.
	Teardown func()      // function for put back the choosen utxos, and reset balance,etc.
}

func createBtcTxInOut(ee engine.Exchange, a account.Accounter, amount uint64, outAddr string) (*txInOutResult, error) {
	var rlt txInOutResult
	// verify the outAddr
	if _, err := cipher.BitcoinDecodeBase58Address(outAddr); err != nil {
		return nil, errors.New("invalid bitcoin address")
	}

	var err error
	// decrease balance and check if the balance is sufficient.
	if err := a.DecreaseBalance(coin.Bitcoin, amount+ee.GetBtcFee()); err != nil {
		return nil, err
	}

	var utxos []bitcoin.Utxo

	// choose sufficient utxos.
	uxs, err := ee.ChooseUtxos(coin.Bitcoin, amount+ee.GetBtcFee(), ChooseUtxoTm)
	if err != nil {
		return nil, err
	}
	utxos = uxs.([]bitcoin.Utxo)

	for _, u := range utxos {
		logger.Debug("using utxos: txid:%s vout:%d addr:%s", u.GetTxid(), u.GetVout(), u.GetAddress())
		rlt.TxIns = append(rlt.TxIns, coin.TxIn{
			Txid: u.GetTxid(),
			Vout: u.GetVout(),
		})
	}

	var totalAmounts uint64
	for _, u := range utxos {
		totalAmounts += u.GetAmount()
	}
	fee := ee.GetBtcFee()
	txOuts := []bitcoin.TxOut{}
	chgAmt := totalAmounts - fee - amount
	chgAddr := ""
	if chgAmt > 0 {
		// generate a change address
		chgAddr = ee.GetNewAddress(coin.Bitcoin)
		txOuts = append(txOuts,
			bitcoin.TxOut{Addr: outAddr, Value: amount},
			bitcoin.TxOut{Addr: chgAddr, Value: chgAmt})
	} else {
		txOuts = append(txOuts, bitcoin.TxOut{Addr: outAddr, Value: amount})
	}

	rlt.TxOuts = txOuts
	rlt.Teardown = func() {
		a.IncreaseBalance(coin.Bitcoin, amount+ee.GetBtcFee())
		ee.PutUtxos(coin.Bitcoin, utxos)
	}

	return &rlt, nil
}

func createSkyTxInOut(ee engine.Exchange, a account.Accounter, amount uint64, outAddr string) (*txInOutResult, error) {
	return nil, nil
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
	// decrease balance and check if the balance is sufficient.
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
	uxs, err := egn.ChooseUtxos(coin.Bitcoin, amount+egn.GetBtcFee(), ChooseUtxoTm)
	if err != nil {
		return nil, err
	}
	utxos := uxs.([]bitcoin.Utxo)

	for _, u := range utxos {
		logger.Debug("using utxos: txid:%s vout:%d addr:%s", u.GetTxid(), u.GetVout(), u.GetAddress())
	}

	var success bool
	defer func() {
		if !success {
			// put utxos back to pool if withdraw failed.
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
	uxs, err := egn.ChooseUtxos(coin.Skycoin, amount, ChooseUtxoTm)
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
