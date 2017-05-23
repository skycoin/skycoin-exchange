package shellcoin

import (
	skycoin "github.com/skycoin/skycoin-exchange/src/coin/skycoin"
	"github.com/skycoin/skycoin-exchange/src/wallet"
)

// Type represents shellcoin coin type
var Type = "shellcoin"

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

// Shellcoin will implement coin.Gateway interface
type Shellcoin struct {
	skycoin.Skycoin // embeded from skycoin , as all apis are the same as skycoin
}

// New creates a shellcoin instance.
func New(nodeAddr string) *Shellcoin {
	return &Shellcoin{Skycoin: skycoin.Skycoin{NodeAddress: nodeAddr}}
}

// Symbol returns the shellcoin symbol
func (sh Shellcoin) Symbol() string {
	return "SC2"
}

// Type returns shellcoin type
func (sh Shellcoin) Type() string {
	return Type
}
