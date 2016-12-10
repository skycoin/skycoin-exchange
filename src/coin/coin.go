package coin

import "github.com/skycoin/skycoin-exchange/src/pp"

// Gateway coin gateway, once a coin implemented this interface,
// then this coin can be registered in this exchange system.
type Gateway interface {
	TxHandler
	Symbol() string // return the coin symbol, SKY, BTC, MZC, etc.
	Type() string   // return the coin type, skycoin, bitcoin, etc.
	// GetBalance interface for getting balance, the return value is an interface{}, cause
	// the balance struct of skycoin and bitcoin are not the same.
	GetBalance(addrs []string) (pp.Balance, error)
	GetOutput(hash string) (interface{}, error)
	GetUtxos(addrs []string) (interface{}, error)
}

// TxHandler transaction handler interface for gateway.
type TxHandler interface {
	GetTx(txid string) (*pp.Tx, error)
	GetRawTx(txid string) (string, error)
	InjectTx(rawtx string) (string, error)
	CreateRawTx(txIns []TxIn, txOuts interface{}) (string, error)
	SignRawTx(rawtx string, getKey GetPrivKey) (string, error)
	ValidateTxid(txid string) bool
}

// TxIn records the tx vin info, txid is the prevous txid, Index is the out index in previous tx.
type TxIn struct {
	Txid    string
	Address string
	Vout    uint32
}

// GetPrivKey is a callback func used for SignTx func to get relevant private key of specific address.
type GetPrivKey func(addr string) (string, error)

// AddressEntry represents the wallet address
type AddressEntry struct {
	Address string `json:"address"`
	Public  string `json:"pubkey"`
	Secret  string `json:"seckey"`
}
