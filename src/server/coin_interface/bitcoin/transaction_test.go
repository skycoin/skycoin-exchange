package bitcoin_interface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Vout struct {
	Addr  string `json:"addr"`
	Value uint64 `json:"value"`
}

type TData struct {
	Address string `json:"address"`
	Key     string `json:"key"`
	OutAddr []Vout `json:"vout"`
	Fee     uint64 `json:"fee"`
}

func TestNewRawTransaction(t *testing.T) {
	d, err := ioutil.ReadFile("test.json")
	assert.Nil(t, err)
	td := TData{}
	err = json.Unmarshal(d, &td)
	assert.Nil(t, err)

	utxos, err := GetUnspentOutputs([]string{td.Address})
	assert.Nil(t, err)
	outAddr := make([]UtxoOut, len(td.OutAddr))
	for i, o := range td.OutAddr {
		outAddr[i].Addr = o.Addr
		outAddr[i].Value = o.Value
	}

	lastIdx := len(outAddr) - 1
	outAddr[lastIdx].Value = outAddr[lastIdx].Value - td.Fee

	utks := make([]UtxoWithkey, len(utxos))
	for i, utxo := range utxos {
		bk := NewUtxoWithKey(utxo, td.Key)
		utks[i] = bk
	}

	tx, err := NewTransaction(utks, outAddr)
	assert.Nil(t, err)
	b := bytes.Buffer{}
	tx.Serialize(&b)
	fmt.Printf("%x\n", b.Bytes())
}
