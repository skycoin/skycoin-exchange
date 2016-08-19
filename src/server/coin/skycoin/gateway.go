package skycoin_interface

import (
	"io"

	"github.com/skycoin/skycoin-exchange/src/server/coin"
)

type Gateway struct{}

func (gw *Gateway) GetTx(txid string) (coin.Transaction, error) {
	return GetTxByID(txid)
}

func (gw *Gateway) GetRawTx(txid string) (string, error) {
	return "skycoin hello world", nil
}

func (gw *Gateway) DecodeRawTx(r io.Reader) (coin.Transaction, error) {
	return nil, nil
}

func (gw *Gateway) InjectTx(tx coin.Transaction) (string, error) {
	return "new skycoin transaction", nil
}
