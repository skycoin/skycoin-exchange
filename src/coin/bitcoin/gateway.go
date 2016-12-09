package bitcoin_interface

import (
	"bytes"
	"encoding/hex"
	"errors"
	"reflect"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/pp"
)

// Bitcoin implements the interface of coin.Gateway.
type Bitcoin struct{}

// GetTx get bitcoin transaction of specific txid.
func (btc Bitcoin) GetTx(txid string) (*pp.Tx, error) {
	return getTxVerboseExplr(txid)
}

// GetRawTx get bitcoin raw transaction of specific txid.
func (btc Bitcoin) GetRawTx(txid string) (string, error) {
	return getRawtxExplr(txid)
}

// InjectTx inject bitcoin raw transaction.
func (btc Bitcoin) InjectTx(rawtx string) (string, error) {
	return BroadcastTx(rawtx)
}

// GetBalance get balance of specific addresses.
func (btc Bitcoin) GetBalance(addrs []string) (pp.Balance, error) {
	v, err := getBalanceExplr(addrs)
	if err != nil {
		return pp.Balance{}, err
	}
	return pp.Balance{Amount: pp.PtrUint64(v)}, nil
}

// CreateRawTx create bitcoin raw transaction.
func (btc Bitcoin) CreateRawTx(txIns []coin.TxIn, txOuts interface{}) (string, error) {
	tx := wire.NewMsgTx()
	oldTxOuts := make([]*wire.TxOut, len(txIns))
	for i, in := range txIns {
		txid, err := chainhash.NewHashFromStr(in.Txid)
		// txid, err := chainhash.NewShaHashFromStr(in.Txid)
		if err != nil {
			return "", err
		}
		rawFundingTx, err := lookupTxid(txid)
		if err != nil {
			return "", err
		}
		oldTxOut, outpoint, err := getFundingParams(rawFundingTx, in.Vout)
		if err != nil {
			return "", err
		}
		oldTxOuts[i] = oldTxOut

		txin := createTxIn(outpoint)
		tx.AddTxIn(txin)
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
		out := o.(TxOut)
		addr, err := btcutil.DecodeAddress(out.Addr, &chaincfg.MainNetParams)
		if err != nil {
			return "", err
		}
		txout := createTxOut(out.Value, addr)
		tx.AddTxOut(txout)
	}

	t := Transaction{*tx}
	d, err := t.Serialize()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(d), nil
}

// SignRawTx sign bitcoin transaction.
func (btc Bitcoin) SignRawTx(rawtx string, getKey coin.GetPrivKey) (string, error) {
	// decode the rawtx
	tx := Transaction{}
	d, err := hex.DecodeString(rawtx)
	if err != nil {
		return "", err
	}

	if err := tx.Deserialize(bytes.NewBuffer(d)); err != nil {
		return "", err
	}

	// get scriptPubkey and addr of the inputs.
	for i, t := range tx.TxIn {
		txid := t.PreviousOutPoint.Hash.String()
		index := t.PreviousOutPoint.Index
		// get the scriptPubkey and addr.
		vt, err := getTxVerboseExplr(txid)
		if err != nil {
			return "", err
		}
		outs := vt.GetBtc().GetVout()
		if int(index) > len(outs) {
			return "", errors.New("error rawtx")
		}
		addr := outs[index].GetScriptPubkey().GetAddresses()
		scriptPubkey := outs[index].GetScriptPubkey().GetHex()
		sp, err := hex.DecodeString(scriptPubkey)
		if err != nil {
			return "", err
		}

		// get private key of specific address in wallet.
		wltPrivKey, err := getKey(addr[0])
		if err != nil {
			return "", err
		}

		sig, err := signRawTx(&tx, i, wltPrivKey, sp)
		if err != nil {
			return "", err
		}

		tx.TxIn[i].SignatureScript = sig
	}
	txb, err := tx.Serialize()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txb), nil
}

// ValidateTxid check if the bitcoin transaction id is validated.
func (btc *Bitcoin) ValidateTxid(txid string) bool {
	if len(txid) != 64 {
		return false
	}

	_, err := hex.DecodeString(txid)
	return err == nil
}
