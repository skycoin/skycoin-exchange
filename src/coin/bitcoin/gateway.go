package bitcoin_interface

import "github.com/skycoin/skycoin-exchange/src/pp"

type Gateway struct {
}

func (gw *Gateway) GetTx(txid string) (*pp.Tx, error) {
	return getTxVerboseExplr(txid)
}

func (gw *Gateway) GetRawTx(txid string) (string, error) {
	return getRawtxExplr(txid)
}

func (gw *Gateway) InjectTx(rawtx string) (string, error) {
	return BroadcastTx(rawtx)
}
