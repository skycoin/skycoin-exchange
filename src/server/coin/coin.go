package coin

import (
	"fmt"

	"github.com/skycoin/skycoin-exchange/src/pp"
)

// gateways records the handlers for different coin.
var gateways = map[Type]Gateway{}

// Gateway coin gateway, once a coin implemented this interface,
// then this coin can be registered in this exchange system.
type Gateway interface {
	TxHandler
}

// TxHandler transaction handlers interface for gateway.
type TxHandler interface {
	GetTx(txid string) (*pp.Tx, error)
	GetRawTx(txid string) (string, error)
	InjectTx(rawtx string) (string, error)
}

// RegisterGateway register gateway for specific coin.
func RegisterGateway(tp Type, gw Gateway) {
	if _, ok := gateways[tp]; ok {
		panic(fmt.Errorf("%s gateway already registered"))
	}
	gateways[tp] = gw
}

// GetGateway get gateway of specific coin.
func GetGateway(tp Type) (Gateway, error) {
	if c, ok := gateways[tp]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("%s handler not registerd")
}

type AddressEntry struct {
	Address string `json:"address"`
	Public  string `json:"pubkey"`
	Secret  string `json:"seckey"`
}

// Type represents the coin type.
type Type int8

const (
	Bitcoin Type = iota
	Skycoin
	// Shellcoin
	// Ethereum
	// other coins...
)

var coinStr = []string{
	Bitcoin: "bitcoin",
	Skycoin: "skycoin",
}

func (c Type) String() string {
	switch c {
	case Bitcoin:
		return coinStr[c]
	case Skycoin:
		return coinStr[c]
	default:
		// return fmt.Sprintf("unknow coin type:%d", c)
		panic(fmt.Sprintf("unknow coin type:%d", c))
	}
}

//TypeFromStr convert string to Type
func TypeFromStr(ct string) (Type, error) {
	switch ct {
	case "bitcoin":
		return Bitcoin, nil
	case "skycoin":
		return Skycoin, nil
	default:
		return -1, fmt.Errorf("unknow coin type:%s", ct)
	}
}
