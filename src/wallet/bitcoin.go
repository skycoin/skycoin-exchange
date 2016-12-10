package wallet

import (
	"encoding/hex"

	"github.com/skycoin/skycoin-exchange/src/coin"
	bitcoin "github.com/skycoin/skycoin-exchange/src/coin/bitcoin"
)

// BtcWallet bitcoin wallet.
type BtcWallet struct {
	walletBase
}

// NewBtcWltCreator wallet generator
func NewBtcWltCreator() Creator {
	return func() Walleter {
		return &BtcWallet{}
	}
}

// GetType return the wallet coin type.
func (bt BtcWallet) GetType() string {
	return "bitcoin"
}

// Copy return the copy of self.
func (bt BtcWallet) Copy() Walleter {
	return &BtcWallet{
		bt.walletBase.Copy(),
	}
}

// NewAddresses generate bitcoin addresses.
func (bt *BtcWallet) NewAddresses(num int) ([]coin.AddressEntry, error) {
	entries := []coin.AddressEntry{}
	defer func() {
		bt.AddressEntries = append(bt.AddressEntries, entries...)
	}()

	if bt.Seed == bt.InitSeed {
		bt.Seed, entries = bitcoin.GenerateAddresses([]byte(bt.Seed), num)
		return entries, nil
	}

	s, err := hex.DecodeString(bt.Seed)
	if err != nil {
		return entries, err
	}
	bt.Seed, entries = bitcoin.GenerateAddresses(s, num)
	return entries, nil
}
