package bitcoin_interface

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type BlockChainInfoTxOut struct {
	Value     int    `json:"value"`
	ScriptHex string `json:"script"`
}

type blockChainInfoTx struct {
	Ver     int                   `json:"ver"`
	Hash    string                `json:"hash"`
	Outputs []BlockChainInfoTxOut `json:"out"`
}

type sendTxJson struct {
	RawTx string `json:"rawtx"`
}

// NewTransaction create transaction,
// utxos is an interface which need to be a slice type, and each item
// of the slice is an UtxoWithPrivkey interface.
// outAddrs is the output address array.
// using the api of blockchain.info to get the raw trasaction info of txid.
func NewTransaction(utxos interface{}, outAddrs []UtxoOut) (*wire.MsgTx, error) {
	s := reflect.ValueOf(utxos)
	if s.Kind() != reflect.Slice {
		return nil, errors.New("error utxo type")
	}

	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	tx := wire.NewMsgTx()
	oldTxOuts := make([]*wire.TxOut, len(ret))
	for i, r := range ret {
		utxo := r.(UtxoWithkey)
		txid, err := wire.NewShaHashFromStr(utxo.GetTxid())
		if err != nil {
			return nil, err
		}
		rawFundingTx, err := lookupTxid(txid)
		if err != nil {
			return nil, err
		}
		oldTxOut, outpoint, err := getFundingParams(rawFundingTx, utxo.GetVout())
		if err != nil {
			return nil, err
		}
		oldTxOuts[i] = oldTxOut

		txin := createTxIn(outpoint)
		tx.AddTxIn(txin)
	}

	if len(outAddrs) > 2 {
		return nil, errors.New("out address more than 2")
	}

	for _, out := range outAddrs {
		addr, err := btcutil.DecodeAddress(out.Addr, &chaincfg.MainNetParams)
		if err != nil {
			return nil, fmt.Errorf("decode address %s, faild, %s", out.Addr, err)
		}
		txout := createTxOut(out.Value, addr)
		tx.AddTxOut(txout)
	}

	// sign the transaction
	for i, r := range ret {
		utxo := r.(UtxoWithkey)
		sig, err := signRawTransaction(tx, i, utxo.GetPrivKey(), oldTxOuts[i].PkScript)
		if err != nil {
			return nil, err
		}
		tx.TxIn[i].SignatureScript = sig
	}
	return tx, nil
}

// func BroadcastTx(tx []byte) (string, error) {
// 	t := wire.MsgTx{}
// 	if err := t.Deserialize(bytes.NewBuffer(tx)); err != nil {
// 		return "", err
// 	}
//
// 	return broadcastTx(&t)
// }

// BroadcastTx tries to send the transaction using an api that will broadcast
// a submitted transaction on behalf of the user.
//
// The transaction is broadcast to the bitcoin network using this API:
//    https://github.com/bitpay/insight-api
//
func BroadcastTx(tx *wire.MsgTx) (string, error) {
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	tx.Serialize(buf)
	hexstr := hex.EncodeToString(buf.Bytes())

	url := "https://insight.bitpay.com/api/tx/send"
	contentType := "application/json"

	// fmt.Printf("Sending transaction to: %s\n", url)
	sendTxJson := &sendTxJson{RawTx: hexstr}
	j, err := json.Marshal(sendTxJson)
	if err != nil {
		return "", fmt.Errorf("Broadcasting the tx failed: %v", err)
	}
	buf = bytes.NewBuffer(j)
	resp, err := http.Post(url, contentType, buf)
	if err != nil {
		return "", fmt.Errorf("Broadcasting the tx failed: %v", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	v := struct {
		Txid string `json:"txid"`
	}{}
	json.Unmarshal(b, &v)
	// fmt.Printf("The sending api responded with:\n%s\n", b)
	return v.Txid, nil
}

func DumpTxBytes(tx *wire.MsgTx) []byte {
	b := bytes.Buffer{}
	tx.Serialize(&b)
	return b.Bytes()
}

// signRawTransaction requires a transaction, a private key, and the bytes of the raw
// scriptPubKey. It will then generate a signature over all of the outputs of
// the provided tx. This is the last step of creating a valid transaction.
func signRawTransaction(tx *wire.MsgTx, index int, wifPrivKey string, scriptPubKey []byte) ([]byte, error) {
	wif, err := btcutil.DecodeWIF(wifPrivKey)
	if err != nil {
		return []byte{}, err
	}

	// The all important signature. Each input is documented below.
	scriptSig, err := txscript.SignatureScript(
		tx,                  // The tx to be signed.
		index,               // The index of the txin the signature is for.
		scriptPubKey,        // The other half of the script from the PubKeyHash.
		txscript.SigHashAll, // The signature flags that indicate what the sig covers.
		wif.PrivKey,         // The key to generate the signature with.
		true,                // The compress sig flag. This saves space on the blockchain.
	)
	if err != nil {
		return []byte{}, err
	}
	return scriptSig, nil
}

// Uses the txid of the target funding transaction and asks blockchain.info's
// api for information (in json) related to that transaction.
func lookupTxid(hash *wire.ShaHash) (*blockChainInfoTx, error) {
	url := "https://blockchain.info/rawtx/" + hash.String()
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Tx Lookup failed: %v", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("TxInfo read failed: %s", err)
	}

	txinfo := &blockChainInfoTx{}
	err = json.Unmarshal(b, txinfo)
	if err != nil {
		return nil, err
	}

	if txinfo.Ver != 1 {
		return nil, fmt.Errorf("Blockchain.info's response seems bad: %v", txinfo)
	}

	return txinfo, nil
}

// getFundingParams pulls the relevant transaction information from the json returned by blockchain.info
// To generate a new valid transaction all of the parameters of the TxOut we are
// spending from must be used.
func getFundingParams(rawtx *blockChainInfoTx, vout uint32) (*wire.TxOut, *wire.OutPoint, error) {
	blkChnTxOut := rawtx.Outputs[vout]

	hash, err := wire.NewShaHashFromStr(rawtx.Hash)
	if err != nil {
		return nil, nil, err
	}

	// Then convert it to a btcutil amount
	amnt := btcutil.Amount(int64(blkChnTxOut.Value))

	outpoint := wire.NewOutPoint(hash, vout)

	subscript, err := hex.DecodeString(blkChnTxOut.ScriptHex)
	if err != nil {
		return nil, nil, err
	}

	oldTxOut := wire.NewTxOut(int64(amnt), subscript)
	return oldTxOut, outpoint, nil
}

// createTxIn pulls the outpoint out of the funding TxOut and uses it as a reference
// for the txin that will be placed in a new transaction.
func createTxIn(outpoint *wire.OutPoint) *wire.TxIn {
	// The second arg is the txin's signature script, which we are leaving empty
	// until the entire transaction is ready.
	txin := wire.NewTxIn(outpoint, []byte{})
	return txin
}

// createTxOut generates a TxOut that can be added to a transaction.
func createTxOut(outCoins uint64, addr btcutil.Address) *wire.TxOut {
	// Take the address and generate a PubKeyScript out of it
	script, err := txscript.PayToAddrScript(addr)
	if err != nil {
		log.Fatal(err)
	}
	txout := wire.NewTxOut(int64(outCoins), script)
	return txout
}
