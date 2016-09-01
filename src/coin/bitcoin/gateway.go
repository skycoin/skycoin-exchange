package bitcoin_interface

import (
	"bytes"
	"encoding/hex"
	"errors"

	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/pp"
)

// Gateway bitcoin gateway which implements the interface of coin.Gateway.
type Gateway struct{}

// GetTx get bitcoin transaction of specific txid.
func (gw Gateway) GetTx(txid string) (*pp.Tx, error) {
	return getTxVerboseExplr(txid)
}

// GetRawTx get bitcoin raw transaction of specific txid.
func (gw Gateway) GetRawTx(txid string) (string, error) {
	return getRawtxExplr(txid)
}

// InjectTx inject bitcoin raw transaction.
func (gw Gateway) InjectTx(rawtx string) (string, error) {
	return BroadcastTx(rawtx)
}

// GetBalance get balance of specific addresses.
func (gw Gateway) GetBalance(addrs []string) (pp.Balance, error) {
	v, err := getBalanceExplr(addrs)
	if err != nil {
		return pp.Balance{}, err
	}
	return pp.Balance{Amount: pp.PtrUint64(v)}, nil
}

// SignRawTx sign bitcoin transaction.
func (gw Gateway) SignRawTx(rawtx string, getKey coin.GetPrivKey) (string, error) {
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
	return "", err
}
