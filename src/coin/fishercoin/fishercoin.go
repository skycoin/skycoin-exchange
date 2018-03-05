package fishercoin

import (
	skycoin "github.com/skycoin/skycoin-exchange/src/coin/skycoin"
	"github.com/skycoin/skycoin-exchange/src/wallet"
)

// Type represents fishercoin coin type
var Type = "fishercoin"

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

// Fishercoin will implement coin.Gateway interface
type Fishercoin struct {
	skycoin.Skycoin // embeded from skycoin , as all apis are the same as skycoin
}

// New creates a fishercoin instance.
func New(nodeAddr string) *Fishercoin {
	return &Fishercoin{Skycoin: skycoin.Skycoin{NodeAddress: nodeAddr}}
}

// Symbol returns the fishercoin symbol
func (fsc Fishercoin) Symbol() string {
	return "FSC"
}

// Type returns fishercoin type
func (fsc Fishercoin) Type() string {
	return Type
}
