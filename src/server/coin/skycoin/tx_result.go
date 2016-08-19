package skycoin_interface

import (
	"encoding/json"

	"github.com/skycoin/skycoin/src/visor"
)

type TxRawResult struct {
	visor.ReadableTransaction
}

func (tx *TxRawResult) Bytes() ([]byte, error) {
	return json.Marshal(tx)
}
