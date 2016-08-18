package bitcoin_interface

import (
	"io"

	"github.com/skycoin/skycoin-exchange/src/server/coin"
)

type Gateway struct {
}

func (gw *Gateway) GetTx(txid string) (coin.Transaction, error) {
	return getTxVerboseExplr(txid)
}

func (gw *Gateway) GetRawTx(txid string) (string, error) {
	return getRawtxExplr(txid)
}

func (gw *Gateway) DecodeRawTx(r io.Reader) (coin.Transaction, error) {
	return nil, nil
}

func (gw *Gateway) InjectTx(tx coin.Transaction) (string, error) {
	return "new bitcoin transaction", nil
}
