package skycoin_interface

import (
	"encoding/hex"

	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/wallet"
)

// Wallet skycoin wallet struct
type Wallet struct {
	wallet.Wallet
}

// newWalletCtor creates coin creator in skycoin ledger
func newWalletCtor() wallet.Creator {
	return func() wallet.Walleter {
		return &Wallet{
			Wallet: wallet.Wallet{
				Type: "skycoin",
			},
		}
	}
}

// NewAddresses generate skycoin addresses.
func (wlt *Wallet) NewAddresses(num int) ([]coin.AddressEntry, error) {
	entries := []coin.AddressEntry{}
	if wlt.Seed == wlt.InitSeed {
		wlt.Seed, entries = GenerateAddresses([]byte(wlt.Seed), num)
		wlt.AddressEntries = append(wlt.AddressEntries, entries...)
		return entries, nil
	}

	s, err := hex.DecodeString(wlt.Seed)
	if err != nil {
		return entries, err
	}
	wlt.Seed, entries = GenerateAddresses(s, num)
	wlt.AddressEntries = append(wlt.AddressEntries, entries...)
	return entries, nil
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
