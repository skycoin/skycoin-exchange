package bitcoin

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TVout struct {
	Addr  string `json:"addr"`
	Value uint64 `json:"value"`
}

type TData struct {
	Address string  `json:"address"`
	Key     string  `json:"key"`
	OutAddr []TVout `json:"vout"`
	Fee     uint64  `json:"fee"`
}

// func TestNewRawTransaction(t *testing.T) {
// 	d, err := ioutil.ReadFile("test.json")
// 	assert.Nil(t, err)
// 	td := TData{}
// 	err = json.Unmarshal(d, &td)
// 	assert.Nil(t, err)
//
// 	utxos, err := GetUnspentOutputs([]string{td.Address})
// 	assert.Nil(t, err)
// 	outAddr := make([]UtxoOut, len(td.OutAddr))
// 	for i, o := range td.OutAddr {
// 		outAddr[i].Addr = o.Addr
// 		outAddr[i].Value = o.Value
// 	}
//
// 	lastIdx := len(outAddr) - 1
// 	outAddr[lastIdx].Value = outAddr[lastIdx].Value - td.Fee
//
// 	utks := make([]UtxoWithkey, len(utxos))
// 	for i, utxo := range utxos {
// 		bk := NewUtxoWithKey(utxo, td.Key)
// 		utks[i] = bk
// 	}
//
// 	tx, err := NewTransaction(utks, outAddr)
// 	assert.Nil(t, err)
// 	// b := bytes.Buffer{}
// 	rawTx, err := tx.Serialize()
// 	assert.Nil(t, err)
// 	fmt.Printf("%x\n", rawTx)
// 	txMsg, err := json.MarshalIndent(tx, "", " ")
// 	assert.Nil(t, err)
// 	fmt.Printf("%s\n", string(txMsg))
// }

// func TestDecodeRawTx(t *testing.T) {
// 	tx := Transaction{}
// 	rawtx := "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff6403a6ab05e4b883e5bda9e7a59ee4bb99e9b1bc76a3a2bb0e9c92f06e4a6349de9ccc8fbe0fad11133ed73c78ee12876334c13c02000000f09f909f2f4249503130302f4d696e65642062792073647a6861626364000000000000000000000000000000005f77dba4015ca34297000000001976a914c825a1ecf2a6830c4401620c3a16f1995057c2ab88acfe75853a"
// 	d, err := hex.DecodeString(rawtx)
// 	assert.Nil(t, err)
//
// 	err = tx.Deserialize(bytes.NewBuffer(d))
// 	assert.Nil(t, err)
// 	v, err := json.MarshalIndent(tx, "", " ")
// 	assert.Nil(t, err)
// 	fmt.Println(string(v))
// }

func TestGetTxVerbose(t *testing.T) {
	tx, err := getTxVerboseExplr("69be3a3b98541e609f5a4935f94c92012d2b3e3437e9508770ba2257f532142f")
	assert.Nil(t, err)
	v, err := json.MarshalIndent(tx, "", " ")
	assert.Nil(t, err)
	fmt.Println(string(v))
}
