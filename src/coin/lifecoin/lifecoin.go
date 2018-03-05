package lifecoin

import (
	skycoin "github.com/skycoin/skycoin-exchange/src/coin/skycoin"
	"github.com/skycoin/skycoin-exchange/src/wallet"
)

// Type represents lifecoin coin type
var Type = "lifecoin"

func init() {
	// Register wallet creator
	wallet.RegisterCreator(Type, func() wallet.Walleter {
		return &skycoin.Wallet{
			Wallet: wallet.Wallet{
				Type: Type,
			},
		}
	})
}

// Lifecoin will implement coin.Gateway interface
type Lifecoin struct {
	skycoin.Skycoin // embeded from skycoin , as all apis are the same as skycoin
}

// New creates a lifecoin instance.
func New(nodeAddr string) *Lifecoin {
	return &Lifecoin{Skycoin: skycoin.Skycoin{NodeAddress: nodeAddr}}
}

// Symbol returns the lifecoin symbol
func (lfc Lifecoin) Symbol() string {
	return "LFC"
}

// Type returns lifecoin type
func (lfc Lifecoin) Type() string {
	return Type
}
