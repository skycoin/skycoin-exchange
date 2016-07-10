// get unspent output from blockexplorer.com
package bitcoin_interface

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type BlkExplrUtxo struct {
	Address      string  `json:"address"`
	Txid         string  `json:"txid"`
	Vout         uint32  `json:"vout"`
	ScriptPubkey string  `json:"criptPubKey"`
	Amount       float64 `json:"amount"`
	Confirms     uint64  `json:"confirmations"`
}

func (be BlkExplrUtxo) GetTxid() string {
	return be.Txid
}

func (be BlkExplrUtxo) GetVout() uint32 {
	return be.Vout
}

func (be BlkExplrUtxo) GetAmount() uint64 {
	return uint64(be.Amount * 100000000)
}

func (be BlkExplrUtxo) GetAddress() string {
	return be.Address
}

// BlkChnUtxo with private key
type BlkExplrUtxoWithkey struct {
	BlkExplrUtxo
	Privkey string
}

func (beu BlkExplrUtxoWithkey) GetPrivKey() string {
	return beu.Privkey
}

func getUtxosBlkExplr(addrs []string) ([]Utxo, error) {
	url := fmt.Sprintf("https://blockexplorer.com/api/addrs/%s/utxo", strings.Join(addrs, ","))
	rsp, err := http.Get(url)
	if err != nil {
		return []Utxo{}, err
	}

	if rsp.StatusCode != 200 {
		return []Utxo{}, errors.New("get unspent output from blockexplorer.com failed")
	}
	defer rsp.Body.Close()
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return []Utxo{}, err
	}
	us := []BlkExplrUtxo{}
	err = json.Unmarshal(data, &us)
	if err != nil {
		return []Utxo{}, err
	}

	utxos := make([]Utxo, len(us))
	for i, u := range us {
		utxos[i] = u
	}
	return utxos, nil
}
