package wallet

import (
	"encoding/hex"

	"github.com/skycoin/skycoin-exchange/src/coin"
	skycoin "github.com/skycoin/skycoin-exchange/src/coin/skycoin"
)

// SkyWallet skycoin wallet.
type SkyWallet struct {
	walletBase
}

// NewSkyWltCreator wallet generator
func NewSkyWltCreator() Creator {
	return func() Walleter {
		return &SkyWallet{}
	}
}

// GetCoinType return the wallet coin type.
func (sk SkyWallet) GetCoinType() coin.Type {
	return coin.Skycoin
}

// Copy return the copy of self.
func (sk SkyWallet) Copy() Walleter {
	return &SkyWallet{
		sk.walletBase.Copy(),
	}
}

// NewAddresses generate skycoin addresses.
func (sk *SkyWallet) NewAddresses(num int) ([]coin.AddressEntry, error) {
	entries := []coin.AddressEntry{}
	defer func() {
		sk.AddressEntries = append(sk.AddressEntries, entries...)
	}()

	if sk.Seed == sk.InitSeed {
		sk.Seed, entries = skycoin.GenerateAddresses([]byte(sk.Seed), num)
		return entries, nil
	}

	s, err := hex.DecodeString(sk.Seed)
	if err != nil {
		return entries, err
	}
	sk.Seed, entries = skycoin.GenerateAddresses(s, num)
	return entries, nil
}
