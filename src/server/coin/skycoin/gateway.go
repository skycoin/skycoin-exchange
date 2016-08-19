package skycoin_interface

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/skycoin/skycoin-exchange/src/server/coin"
)

type Gateway struct{}

func (gw *Gateway) GetTx(txid string) (coin.Transaction, error) {
	url := fmt.Sprintf("%s/transaction?txid=%s", ServeAddr, txid)
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	txRlt := TxRawResult{}
	if err := json.NewDecoder(rsp.Body).Decode(&txRlt); err != nil {
		return nil, err
	}
	return &txRlt, nil
}

// GetRawTx get raw tx by txid.
func (gw *Gateway) GetRawTx(txid string) (string, error) {
	url := fmt.Sprintf("%s/rawtx?txid=%s", ServeAddr, txid)
	rsp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	res := struct {
		Rawtx string `json:"rawtx"`
	}{}
	if err := json.NewDecoder(rsp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res.Rawtx, nil
}

func (gw *Gateway) DecodeRawTx(r io.Reader) (coin.Transaction, error) {
	return nil, nil
}

func (gw *Gateway) InjectTx(tx coin.Transaction) (string, error) {
	return "new skycoin transaction", nil
}
