package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	skycoin "github.com/skycoin/skycoin-exchange/src/coin/skycoin"
)

var rawtx string

// type Transaction struct {
// 	Length    uint32        //length prefix
// 	Type      uint8         //transaction type
// 	InnerHash cipher.SHA256 //inner hash SHA256 of In[],Out[]

// 	Sigs []cipher.Sig        //list of signatures, 64+1 bytes each
// 	In   []cipher.SHA256     //ouputs being spent
// 	Out  []TransactionOutput //ouputs being created
// }

type TxReadable struct {
	Length    uint32 //length prefix
	Type      uint8  //transaction type
	InnerHash string //inner hash SHA256 of In[],Out[]

	Sigs []string      //list of signatures, 64+1 bytes each
	In   []string      //ouputs being spent
	Out  []OutReadable //ouputs being created
}

type OutReadable struct {
	Address string //address to send to
	Coins   uint64 //amount to be sent in coins
	Hours   uint64 //amount to be sent in coin hours
}

func ToReadable(tx *skycoin.Transaction) *TxReadable {
	rdTx := &TxReadable{
		Length:    tx.Length,
		Type:      tx.Type,
		InnerHash: tx.InnerHash.Hex(),
		Sigs:      make([]string, len(tx.Sigs)),
		In:        make([]string, len(tx.In)),
		Out:       make([]OutReadable, len(tx.Out)),
	}

	for i, sig := range tx.Sigs {
		rdTx.Sigs[i] = sig.Hex()
	}

	for i, in := range tx.In {
		rdTx.In[i] = in.Hex()
	}

	for i, o := range tx.Out {
		rdTx.Out[i] = OutReadable{
			Address: o.Address.String(),
			Coins:   o.Coins,
			Hours:   o.Hours,
		}
	}
	return rdTx
}

func main() {
	flag.StringVar(&rawtx, "rawtx", "", "rawtx")
	flag.Parse()

	tx := skycoin.Transaction{}
	b, err := hex.DecodeString(rawtx)
	if err != nil {
		log.Fatal(err)
	}

	if err := tx.Deserialize(bytes.NewBuffer(b)); err != nil {
		log.Fatal(err)
	}

	readable := ToReadable(&tx)

	v, err := json.MarshalIndent(readable, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(v))
}
