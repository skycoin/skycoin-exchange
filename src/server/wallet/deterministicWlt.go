package wallet

import (
	"sync"

	"github.com/skycoin/skycoin/src/cipher"
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
		seckeys := cipher.GenerateDeterministicKeyPairs([]byte(self.Seed), num)
		for i, sec := range seckeys {
			addr := AddressEntry{}
			pub := cipher.PubKeyFromSecKey(sec)
			addr.Address = cipher.BitcoinAddressFromPubkey(pub)
			addr.Pubkey = pub.Hex()
			addr.Seckey = cipher.BitcoinWalletImportFormatFromSeckey(sec)
			addrEntries[i] = addr
		}
		self.addAddresses(coinType, addrEntries)
		return addrEntries
	case Skycoin:
	}
	return []AddressEntry{}
}

func (self *DeterministicWallet) GetBalance(addr string) (string, error) {
	return "", nil
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
