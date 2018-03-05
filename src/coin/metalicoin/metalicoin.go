package metalicoin

import (
	skycoin "github.com/skycoin/skycoin-exchange/src/coin/skycoin"
	"github.com/skycoin/skycoin-exchange/src/wallet"
)

// Type represents metalicoin coin type
var Type = "metalicoin"

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

// Metalicoin will implement coin.Gateway interface
type Metalicoin struct {
	skycoin.Skycoin // embeded from skycoin , as all apis are the same as skycoin
}

// New creates a metalicoin instance.
func New(nodeAddr string) *Metalicoin {
	return &Metalicoin{Skycoin: skycoin.Skycoin{NodeAddress: nodeAddr}}
}

// Symbol returns the metalicoin symbol
func (mt Metalicoin) Symbol() string {
	return "MTC"
}

// Type returns metalicoin type
func (mt Metalicoin) Type() string {
	return Type
}
