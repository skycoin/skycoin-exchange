package aynrandcoin

import (
	skycoin "github.com/skycoin/skycoin-exchange/src/coin/skycoin"
	"github.com/skycoin/skycoin-exchange/src/wallet"
)

// Type represents aynrandcoin coin type
var Type = "aynrandcoin"

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

// Aynrandcoin will implement coin.Gateway interface
type Aynrandcoin struct {
	skycoin.Skycoin // embeded from skycoin , as all apis are the same as skycoin
}

// New creates a aynrand instance.
func New(nodeAddr string) *Aynrandcoin {
	return &Aynrandcoin{Skycoin: skycoin.Skycoin{NodeAddress: nodeAddr}}
}

// Symbol returns the aynrandcoin symbol
func (ayn Aynrandcoin) Symbol() string {
	return "ARC"
}

// Type returns aynrandcoin type
func (ayn Aynrandcoin) Type() string {
	return Type
}
