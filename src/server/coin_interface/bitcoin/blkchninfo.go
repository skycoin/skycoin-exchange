// using the api and data struct of blockchain.info to handle the unspent output.
package bitcoin_interface

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// unspent output response struct of blockchain.info.
type BlkChnUtxoRsp struct {
	Utxos []BlkChnUtxo `json:"unspent_outputs"`
}

type BlkChnUtxo struct {
	Tx_hash            string `json:"tx_hash"` // the previous transaction id
	Tx_hash_big_endian string `json:"tx_hash_big_endian"`
	Tx_index           uint64 `json:"tx_index"`
	Tx_output_n        uint64 `json:"tx_output_n"` // the output index of previous transaction
	Script             string `json:"script"`      // pubkey script
	Value              uint64 `json:"value"`       // the bitcoin amount in satoshis
	Value_hex          string `json:"value_hex"`   // alisa the Value, in hex format.
	Confirmations      uint64 `json:"confirmations"`
}

// BlkChnUtxo with private key
type BlkChnUtxoWithkey struct {
	BlkChnUtxo
	Privkey string
}

func (bo BlkChnUtxo) GetTxid() string {
	return bo.Tx_hash_big_endian
}

func (bo BlkChnUtxo) GetVout() uint32 {
	return uint32(bo.Tx_output_n)
}

func (bo BlkChnUtxo) GetAmount() uint64 {
	return bo.Value
}

func (bk BlkChnUtxoWithkey) GetPrivKey() string {
	return bk.Privkey
}

// GetUtxosBlkChnInfo get unspent outputs from blockchain.info
// https://blockchain.info/unspent?active=1SakrZuzQmGwn7MSiJj5awqJZjSYeBWC3
func getUtxosBlkChnInfo(addr string) []Utxo {
	if AddressValid(addr) != nil {
		log.Fatal("Address is invalid")
	}
	url := fmt.Sprintf("https://blockchain.info/unspent?active=%s", addr)
	// fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Get url:%s fail, error:%s", addr, err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Read data from resp body fail, error:%s", err)
	}
	resp.Body.Close()
	// fmt.Println("status:", resp.StatusCode)

	// parse the JSON.
	utxoResp := BlkChnUtxoRsp{}
	err = json.Unmarshal(data, &utxoResp)
	if err != nil {
		log.Fatalf("unmasharl fail, error:%s", err)
	}

	utxos := make([]Utxo, len(utxoResp.Utxos))
	for i, u := range utxoResp.Utxos {
		utxos[i] = u
	}
	return utxos
}
