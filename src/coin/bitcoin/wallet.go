package bitcoin_interface

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
func (bt *Wallet) NewAddresses(num int) ([]coin.AddressEntry, error) {
	entries := []coin.AddressEntry{}
	defer func() {
		bt.AddressEntries = append(bt.AddressEntries, entries...)
	}()

	if bt.Seed == bt.InitSeed {
		bt.Seed, entries = GenerateAddresses([]byte(bt.Seed), num)
		return entries, nil
	}

	s, err := hex.DecodeString(bt.Seed)
	if err != nil {
		return entries, err
	}
	bt.Seed, entries = GenerateAddresses(s, num)
	return entries, nil
}
