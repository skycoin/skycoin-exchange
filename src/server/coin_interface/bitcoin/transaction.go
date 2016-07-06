package bitcoin_interface

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

	//fmt.Printf("%s\n", b)
	txinfo := &blockChainInfoTx{}
	err = json.Unmarshal(b, txinfo)
	if err != nil {
		log.Fatal(err)
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

	if err != nil {
		return nil, nil, err
	}

	outpoint := wire.NewOutPoint(hash, vout)

	subscript, err := hex.DecodeString(blkChnTxOut.ScriptHex)
	if err != nil {
		return nil, nil, err
	}

	oldTxOut := wire.NewTxOut(int64(amnt), subscript)
	return oldTxOut, outpoint, nil
}

// NewTransaction create transaction
func NewRawTransaction(utxos []UnspentOutput, toAddr string, fee int64, defaultNet *chaincfg.Params) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx()
	var inCoins int64

	for _, utxo := range utxos {
		txid, err := wire.NewShaHashFromStr(utxo.GetTxid())
		if err != nil {
			return nil, err
		}
		rawFundingTx, err := lookupTxid(txid)
		if err != nil {
			return nil, err
		}
		oldTxOut, outpoint, err := getFundingParams(rawFundingTx, utxo.GetVout())
		inCoins += oldTxOut.Value

		txin := createTxIn(outpoint)
		tx.AddTxIn(txin)
	}

	// convert toAddr stringt o btcutil.Address
	addr, err := btcutil.DecodeAddress(toAddr, defaultNet)
	if err != nil {
		return nil, err
	}
	// create tx out.
	txout := createTxOut(inCoins, addr, fee)
	tx.AddTxOut(txout)

	return tx, nil
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
func createTxOut(inCoin int64, addr btcutil.Address, fee int64) *wire.TxOut {
	// Pay the minimum network fee so that nodes will broadcast the tx.
	outCoin := inCoin - fee
	// Take the address and generate a PubKeyScript out of it
	script, err := txscript.PayToAddrScript(addr)
	if err != nil {
		log.Fatal(err)
	}
	txout := wire.NewTxOut(outCoin, script)
	return txout
}
