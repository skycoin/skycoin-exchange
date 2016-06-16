package wallet

import (
	"fmt"
	"sync"

	bitcoin "github.com/skycoin/skycoin-exchange/src/server/coin_interface/bitcoin"
)

// DeterministicWallet, generate and store addresses for various coin types.
type DeterministicWallet struct {
	ID        string                      // wallet id
	Seed      string                      // used to generate address
	Addresses map[CoinType][]AddressEntry // key is coin type, value is address entry list.
	addrLock  sync.Mutex                  // a lock, for protecting the writing, reading of the Addresses in wallet.
}

// GenerateAddress, generate new addresses base on the and coin type, and then store the address.
func (self *DeterministicWallet) NewAddresses(coinType CoinType, num int) []AddressEntry {
	switch coinType {
	case Bitcoin:
		addrEntries := make([]AddressEntry, num)
		entries := bitcoin.GenerateAddresses(self.Seed, num)
		for i, entry := range entries {
			addrEntries[i] = AddressEntry{
				CoinType: Bitcoin,
				Address:  entry.Address,
				Pubkey:   entry.Public,
				Seckey:   entry.Secret}
		}
		self.addAddresses(coinType, addrEntries)
		return addrEntries
	case Skycoin:
	}
	return []AddressEntry{}
}

// GetBalance of specific address.
func (self *DeterministicWallet) GetBalance(addr AddressEntry) (string, error) {
	switch addr.CoinType {
	case Bitcoin:
		return bitcoin.GetBalance(addr.Address)
	default:
		return "", fmt.Errorf("unknow coin type:%d", addr.CoinType)
	}
}

func (self *DeterministicWallet) addAddresses(coinType CoinType, addrs []AddressEntry) {
	self.addrLock.Lock()
	self.Addresses[coinType] = append(self.Addresses[coinType], addrs...)
	self.addrLock.Unlock()
}

func (self *DeterministicWallet) SetID(id string) {
	self.ID = id
}

func (self *DeterministicWallet) GetID() string {
	return self.ID
}
