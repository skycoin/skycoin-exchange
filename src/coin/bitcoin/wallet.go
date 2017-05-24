package bitcoin

import (
	"encoding/hex"

	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/wallet"
)

// Wallet represents the bitcoin wallet
type Wallet struct {
	wallet.Wallet
}

func init() {
	// Register wallet creator
	wallet.RegisterCreator(Type, func() wallet.Walleter {
		return &Wallet{
			Wallet: wallet.Wallet{
				Type: Type,
			},
		}
	})
}

// NewAddresses generate bitcoin addresses.
func (wlt *Wallet) NewAddresses(num int) ([]coin.AddressEntry, error) {
	entries := []coin.AddressEntry{}
	defer func() {
		wlt.AddressEntries = append(wlt.AddressEntries, entries...)
	}()

	if wlt.Seed == wlt.InitSeed {
		wlt.Seed, entries = GenerateAddresses([]byte(wlt.Seed), num)
		return entries, nil
	}

	s, err := hex.DecodeString(wlt.Seed)
	if err != nil {
		return entries, err
	}
	wlt.Seed, entries = GenerateAddresses(s, num)
	return entries, nil
}

// Copy returns copy of self
func (wlt *Wallet) Copy() wallet.Walleter {
	return &Wallet{
		Wallet: wlt.Wallet.Copy(),
	}
}
