// get unspent output from blockexplorer.com
package bitcoin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin/src/cipher"
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

	for _, a := range addrs {
		if !validateAddress(a) {
			return []Utxo{}, fmt.Errorf("invalid bitcoin address %v", a)
		}
	}

	url := fmt.Sprintf("https://blockexplorer.com/api/addrs/%s/utxo", strings.Join(addrs, ","))
	rsp, err := http.Get(url)
	if err != nil {
		return []Utxo{}, fmt.Errorf("get utxo from blockexplorer.com failed")
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

	if strings.ToLower(string(d)) == "not found" {
		return nil, fmt.Errorf("not found")
	}

	tx := pp.Tx{}
	if err := json.Unmarshal(d, &tx.Btc); err != nil {
		return nil, err
	}
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

type balanceResult struct {
	balance uint64
	err     error
}

func getBalanceExplr(addrs []string) (uint64, error) {
	var wg sync.WaitGroup

	valueChan := make(chan balanceResult, len(addrs))

	for _, addr := range addrs {
		// verify the address.
		_, err := cipher.BitcoinDecodeBase58Address(addr)
		if err != nil {
			return 0, err
		}

		wg.Add(1)
		go func(addr string, wg *sync.WaitGroup, vc chan balanceResult) {
			defer wg.Done()
			d, err := getDataOfUrl(fmt.Sprintf("https://blockexplorer.com/api/addr/%s/balance", addr))
			if err != nil {
				vc <- balanceResult{0, err}
				return
			}
			v, err := strconv.ParseUint(string(d), 10, 64)
			if err != nil {
				vc <- balanceResult{0, err}
			}
			vc <- balanceResult{v, nil}
		}(addr, &wg, valueChan)
	}
	wg.Wait()
	close(valueChan)
	var totalBal uint64
	for v := range valueChan {
		if v.err != nil {
			return 0, v.err
		}
		totalBal += v.balance
	}
	return totalBal, nil
}
