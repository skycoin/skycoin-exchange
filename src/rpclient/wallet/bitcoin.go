package wallet

import (
	"os"
	"path/filepath"

	"github.com/skycoin/skycoin-exchange/src/server/coin"
	"github.com/skycoin/skycoin/src/util"
)

// BtcWallet bitcoin wallet.
type BtcWallet struct {
	walletBase
}

// GetID return wallet id.
func (bt BtcWallet) GetID() string {
	return bt.ID
}

// GetCoinType return the wallet coin type.
func (bt BtcWallet) GetCoinType() coin.Type {
	return coin.Bitcoin
}

// NewAddresses generate bitcoin addresses.
func (bt *BtcWallet) NewAddresses(num int) ([]coin.AddressEntry, error) {
	return []coin.AddressEntry{}, nil
}

// GetAddresses return all wallet addresses.
func (bt *BtcWallet) GetAddresses() []string {
	return []string{}
}

// GetAddressEntries get all address enties.
func (bt *BtcWallet) GetAddressEntries() []coin.AddressEntry {
	return []coin.AddressEntry{}
}

// GetAddrEntyByAddr get address entry of specific address.
func (bt *BtcWallet) GetAddrEntyByAddr(addr string) (coin.AddressEntry, error) {
	return coin.AddressEntry{}, nil
}

// Save save the wallet
func (bt *BtcWallet) Save() error {
	return util.SaveJSON(filepath.Join(wltDir, bt.ID), bt, 0600)
}

// Clear remove wallet file from local disk.
func (bt *BtcWallet) Clear() error {
	path := filepath.Join(wltDir, bt.ID)
	return os.RemoveAll(path)
}
