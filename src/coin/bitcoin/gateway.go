package bitcoin_interface

import "github.com/skycoin/skycoin-exchange/src/pp"

// Gateway bitcoin gateway which implements the interface of coin.Gateway.
type Gateway struct{}

// GetTx get bitcoin transaction of specific txid.
func (gw *Gateway) GetTx(txid string) (*pp.Tx, error) {
	return getTxVerboseExplr(txid)
}

// GetRawTx get bitcoin raw transaction of specific txid.
func (gw *Gateway) GetRawTx(txid string) (string, error) {
	return getRawtxExplr(txid)
}

// InjectTx inject bitcoin raw transaction.
func (gw *Gateway) InjectTx(rawtx string) (string, error) {
	return BroadcastTx(rawtx)
}

// GetBalance get balance of specific addresses.
func (gw *Gateway) GetBalance(addrs []string) (pp.Balance, error) {
	v, err := getBalanceExplr(addrs)
	if err != nil {
		return pp.Balance{}, err
	}
	return pp.Balance{Amount: pp.PtrUint64(v)}, nil
}
