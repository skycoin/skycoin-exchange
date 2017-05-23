package suncoin

import (
	skycoin "github.com/skycoin/skycoin-exchange/src/coin/skycoin"
)

// Type represents mzcoin coin type
var Type = "suncoin"

// Suncoin will implement coin.Gateway interface
type Suncoin struct {
	skycoin.Skycoin // embeded from skycoin , as all apis are the same as skycoin
}

// New creates a mzcoin instance.
func New(nodeAddr string) *Suncoin {
	return &Suncoin{Skycoin: skycoin.Skycoin{NodeAddress: nodeAddr}}
}

// Symbol returns the suncoin symbol
func (sun Suncoin) Symbol() string {
	return "SUN"
}

// Type returns mzcoin type
func (sun Suncoin) Type() string {
	return Type
}
