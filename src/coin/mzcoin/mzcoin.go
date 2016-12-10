package mzcoin

import (
	skycoin "github.com/skycoin/skycoin-exchange/src/coin/skycoin"
)

// Type represents mzcoin coin type
var Type = "mzcoin"

// Mzcoin will implement coin.Gateway interface
type Mzcoin struct {
	skycoin.Skycoin        // embeded from skycoin , as all apis are the same as skycoin
	NodeAddress     string // mzcoin's node address
}

// Symbol returns the mzcoin symbol
func (mz Mzcoin) Symbol() string {
	return "MZC"
}

// Type returns mzcoin type
func (mz Mzcoin) Type() string {
	return Type
}
