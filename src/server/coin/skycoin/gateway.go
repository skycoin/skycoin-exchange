package skycoin_interface

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/skycoin/skycoin-exchange/src/pp"
)

type Gateway struct{}

func (gw *Gateway) GetTx(txid string) (*pp.Tx, error) {
	url := fmt.Sprintf("%s/transaction?txid=%s", ServeAddr, txid)
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	txRlt := pp.Tx{}
	if err := json.NewDecoder(rsp.Body).Decode(&txRlt.Sky); err != nil {
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

// func (gw *Gateway) DecodeRawTx(r io.Reader) (coin.Transaction, error) {
// 	return nil, nil
// }

func (gw *Gateway) InjectTx(rawtx string) (string, error) {
	return "new skycoin transaction", nil
}
