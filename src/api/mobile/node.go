package mobile

import "github.com/skycoin/skycoin-exchange/src/coin"

type noder interface {
	GetBalance(addrs []string) (uint64, error)
	ValidateAddr(addr string) error
	PrepareTx(params interface{}) ([]coin.TxIn, interface{}, error)
	CreateRawTx(txIns []coin.TxIn, getKey coin.GetPrivKey, txOuts interface{}) (string, error)
	BroadcastTx(rawtx string) (string, error)
	GetTransactionByID(txid string) (string, error)
	GetNodeAddr() string
}
