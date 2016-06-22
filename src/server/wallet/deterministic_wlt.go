package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"

	bitcoin "github.com/skycoin/skycoin-exchange/src/server/coin_interface/bitcoin"
	"github.com/skycoin/skycoin/src/util"
)

const (
	DETER_META_ID = iota
	DETER_META_SEED
	DETER_META_INIT_SEED
	DETER_META_WALLET_TYPE
)

var DeterMetaStr = []string{
	DETER_META_ID:          "wallet_id",
	DETER_META_SEED:        "seed",
	DETER_META_INIT_SEED:   "init_seed",
	DETER_META_WALLET_TYPE: "wallet_type",
}

// DeterministicWallet generate and store addresses for various coin types.
// This wallet is not thread safe.
type DeterministicWallet struct {
	ID             string // wallet id
	InitSeed       string // Init seed, used to recover the wallet.
	Seed           string // seed used to create new address.
	AddressEntries map[string][]AddressEntry
	// addrLock       sync.Mutex // a lock, for protecting the writing, reading of the Addresses in wallet.
	// fileLock       sync.Mutex // lock for protecting wallet file.
}

// NewAddresses generate new addresses base on the coin type, and then store the address.
// NewAddress must be Sequential，cause the seed.
// this function is not thread safe, should not be used concrrently.
func (self *DeterministicWallet) NewAddresses(coinType CoinType, num int) ([]AddressEntry, error) {
	switch coinType {
	case Bitcoin:
		addrEntries := make([]AddressEntry, num)
		if self.Seed == self.InitSeed {
			sd, entries := bitcoin.GenerateAddresses([]byte(self.Seed), num)
			self.Seed = sd
			addressEntryCopy(&addrEntries, entries)
		} else {
			s, err := hex.DecodeString(self.Seed)
			if err != nil {
				return []AddressEntry{}, err
			}
			sd, entries := bitcoin.GenerateAddresses(s, num)
			self.Seed = sd
			addressEntryCopy(&addrEntries, entries)
		}
		self.addAddresses(coinType, addrEntries)
		// save automaticaly after new addressess are added.
		self.save(dataDir)
		return addrEntries, nil
	case Skycoin:
	default:
	}
	return []AddressEntry{}, nil
}

// Save the wallet
func (self *DeterministicWallet) save(dir string) error {
	// self.fileLock.Lock()
	// defer self.fileLock.Unlock()
	w := self.toWalletBase()
	return util.SaveJSON(filepath.Join(dir, self.GetID()), w, 0600)
}

func (self *DeterministicWallet) SetID(id string) {
	self.ID = id
}

func (self *DeterministicWallet) GetID() string {
	return self.ID
}

func (self *DeterministicWallet) toWalletBase() WalletBase {
	w := WalletBase{
		Meta: map[string]string{
			DeterMetaStr[DETER_META_ID]:          self.ID,
			DeterMetaStr[DETER_META_SEED]:        self.Seed,
			DeterMetaStr[DETER_META_INIT_SEED]:   self.InitSeed,
			DeterMetaStr[DETER_META_WALLET_TYPE]: Deterministic.String()},
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
		id        string
		seed      string
		init_seed string
		ok        bool
	)

	if id, ok = w.Meta[DeterMetaStr[DETER_META_ID]]; !ok {
		return nil, errors.New("invalid wallet meta info, empty wallet_id")
	}

	if seed, ok = w.Meta[DeterMetaStr[DETER_META_SEED]]; !ok {
		return nil, errors.New("invalid wallet meta info, empty seed")
	}

	if init_seed, ok = w.Meta[DeterMetaStr[DETER_META_INIT_SEED]]; !ok {
		return nil, errors.New("invalid wallet meta info, empty init seed")
	}

	wlt := &DeterministicWallet{
		ID:             id,
		Seed:           seed,
		InitSeed:       init_seed,
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
	// self.addrLock.Lock()
	self.AddressEntries[coinType.String()] = append(self.AddressEntries[coinType.String()], addrs...)
	// self.addrLock.Unlock()
}

func addressEntryCopy(dst *[]AddressEntry, src []bitcoin.AddressEntry) {
	for i, e := range src {
		(*dst)[i] = AddressEntry{
			Address: e.Address,
			Public:  e.Public,
			Secret:  e.Secret}
	}
}
