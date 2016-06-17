package wallet

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	bitcoin "github.com/skycoin/skycoin-exchange/src/server/coin_interface/bitcoin"
	"github.com/skycoin/skycoin/src/util"
)

const (
	DETER_META_ID = iota
	DETER_META_SEED
	DETER_META_WALLET_TYPE
)

var DeterMetaStr = []string{
	DETER_META_ID:          "wallet_id",
	DETER_META_SEED:        "seed",
	DETER_META_WALLET_TYPE: "wallet_type",
}

// DeterministicWallet, generate and store addresses for various coin types.
type DeterministicWallet struct {
	// WalletBase
	ID             string                    `json:"wallet_id"` // wallet id
	Seed           string                    `json:"seed"`      // used to generate address
	WalletType     string                    `json:"wallet_type"`
	AddressEntries map[string][]AddressEntry `json:"addresses"`
	addrLock       sync.Mutex                // a lock, for protecting the writing, reading of the Addresses in wallet.
	fileLock       sync.Mutex                // lock for protecting wallet file.
}

// GenerateAddress, generate new addresses base on the and coin type, and then store the address.
func (self *DeterministicWallet) NewAddresses(coinType CoinType, num int) []AddressEntry {
	switch coinType {
	case Bitcoin:
		addrEntries := make([]AddressEntry, num)
		entries := bitcoin.GenerateAddresses(self.Seed, num)
		for i, entry := range entries {
			addrEntries[i] = AddressEntry{
				Address: entry.Address,
				Public:  entry.Public,
				Secret:  entry.Secret}
		}
		self.addAddresses(coinType, addrEntries)
		// save automaticaly after new addressess are added.
		self.Save(dataDir)
		return addrEntries
	case Skycoin:
	default:
	}
	return []AddressEntry{}
}

// GetBalance of specific address.
func (self *DeterministicWallet) GetBalance(addr string, coinType CoinType) (string, error) {
	switch coinType {
	case Bitcoin:
		return bitcoin.GetBalance(addr)
	default:
		return "", fmt.Errorf("unknow coin type:%d", coinType)
	}
}

// save the wallet
func (self *DeterministicWallet) Save(dir string) error {
	w := self.ToWalletBase()
	self.fileLock.Lock()
	defer self.fileLock.Unlock()
	return util.SaveJSON(filepath.Join(dir, self.GetID()), w, 0600)
}

func (self *DeterministicWallet) SetID(id string) {
	self.ID = id
}

func (self *DeterministicWallet) GetID() string {
	return self.ID
}

func (self *DeterministicWallet) ToWalletBase() WalletBase {
	w := WalletBase{
		Meta: map[string]string{
			DeterMetaStr[DETER_META_ID]:          self.ID,
			DeterMetaStr[DETER_META_SEED]:        self.Seed,
			DeterMetaStr[DETER_META_WALLET_TYPE]: self.WalletType},
		AddressEntries: make(map[string][]AddressEntry),
	}

	for k, entries := range self.AddressEntries {
		newEntries := make([]AddressEntry, len(entries))
		for i, e := range entries {
			newEntries[i] = e
		}
		w.AddressEntries[k] = newEntries
	}
	return w
}

func newDeterministicWalletFromBase(w *WalletBase) (*DeterministicWallet, error) {
	var (
		id         string
		seed       string
		wallettype string
		ok         bool
	)

	if id, ok = w.Meta[DeterMetaStr[DETER_META_ID]]; !ok {
		return nil, errors.New("invalid wallet meta info, empty wallet_id")
	}

	if seed, ok = w.Meta[DeterMetaStr[DETER_META_SEED]]; !ok {
		return nil, errors.New("invalid wallet meta info, empty seed")
	}

	if wallettype, ok = w.Meta[DeterMetaStr[DETER_META_WALLET_TYPE]]; !ok {
		return nil, errors.New("invalid wallet meta info, empty wallettype")
	}

	wlt := &DeterministicWallet{
		ID:             id,
		Seed:           seed,
		WalletType:     wallettype,
		AddressEntries: make(map[string][]AddressEntry),
	}

	for k, entries := range w.AddressEntries {
		newEntries := make([]AddressEntry, len(entries))
		for i, e := range entries {
			newEntries[i] = e
		}
		wlt.AddressEntries[k] = newEntries
	}

	if err := validateWallet(wlt); err != nil {
		return nil, fmt.Errorf("invalide wallet, error:%s", err)
	}
	return wlt, nil
}

// TODO: validate the id, seed, wallet_type, and addressEntries.
func validateWallet(wlt *DeterministicWallet) error {
	return nil
}

func (self *DeterministicWallet) addAddresses(coinType CoinType, addrs []AddressEntry) {
	self.addrLock.Lock()
	self.AddressEntries[CoinStr[coinType]] = append(self.AddressEntries[CoinStr[coinType]], addrs...)
	self.addrLock.Unlock()
}
