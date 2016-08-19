// get unspent output from blockexplorer.com
package bitcoin_interface

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/skycoin/skycoin-exchange/src/pp"
)

type BlkExplrUtxo struct {
	Address      string `json:"address"`
	Txid         string `json:"txid"`
	Vout         uint32 `json:"vout"`
	ScriptPubkey string `json:"criptPubKey"`
	Amount       uint64 `json:"satoshis"`
	Confirms     uint64 `json:"confirmations"`
}

func (be BlkExplrUtxo) GetTxid() string {
	return be.Txid
}

func (be BlkExplrUtxo) GetVout() uint32 {
	return be.Vout
}

func (be BlkExplrUtxo) GetAmount() uint64 {
	return be.Amount
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
	if len(addrs) == 0 {
		return []Utxo{}, nil
	}
	url := fmt.Sprintf("https://blockexplorer.com/api/addrs/%s/utxo", strings.Join(addrs, ","))
	rsp, err := http.Get(url)
	if err != nil {
		return []Utxo{}, errors.New("get utxo from blockexplorer.com failed")
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

// get tx verbose from blockexplorer.com
func getTxVerboseExplr(txid string) (*pp.Tx, error) {
	d, err := getDataOfUrl(fmt.Sprintf("https://blockexplorer.com/api/tx/%s", txid))
	if err != nil {
		return nil, err
	}
	tx := pp.Tx{}
	if err := json.Unmarshal(d, &tx.Btc); err != nil {
		return nil, err
	}
	logger.Debug("%v", tx)
	return &tx, nil
}

func getRawtxExplr(txid string) (string, error) {
	d, err := getDataOfUrl(fmt.Sprintf("https://blockexplorer.com/api/rawtx/%s", txid))
	if err != nil {
		return "", err
	}
	v := struct {
		Rawtx string `json:"rawtx"`
	}{}
	if err := json.Unmarshal(d, &v); err != nil {
		return "", err
	}
	return v.Rawtx, nil
}
