package mobile

import "github.com/skycoin/skycoin-exchange/src/coin"
import "github.com/skycoin/skycoin/src/cipher"

type noder interface {
	GetBalance(addrs []string) (uint64, error)
	ValidateAddr(addr string) error
	PrepareTx(addrs []string, toAddr string, amt uint64) ([]coin.TxIn, []string, interface{}, error)
	CreateRawTx(txIns []coin.TxIn, keys []cipher.SecKey, txOuts interface{}) (string, error)
	BroadcastTx(rawtx string) (string, error)
}
