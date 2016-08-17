package coin_interface

import (
	"fmt"
	"io"

	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
)

type AddressEntry struct {
	Address string
	Public  string
	Secret  string
}

// CoinHandlers records the handlers for different coin.
var gateways = map[wallet.CoinType]Gateway{}

// Transaction tx interface
type Transaction interface {
	Serialize() (string, error)
	Deserialize(r io.Reader) error
	ToPPTx() *pp.Tx // translate to *pp.Tx
}

// CoinHandler interface for handlering all coin relevance things.
type Gateway interface {
	TxHandler
}

type TxHandler interface {
	GetTx(txid string) (Transaction, error)
	GetRawTx(txid string) (string, error)
	DecodeRawTx(rawtx string) (Transaction, error)
}

func RegisterGateway(tp wallet.CoinType, gw Gateway) {
	if _, ok := gateways[tp]; ok {
		panic(fmt.Errorf("%s gateway already registered"))
	}
	coinHandlers[tp] = gw
}

func GetGateway(tp wallet.CoinType) (Gateway, error) {
	if c, ok := gateways[tp]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("%s handler not registerd")
}
