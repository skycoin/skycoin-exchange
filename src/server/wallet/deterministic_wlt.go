package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/skycoin/skycoin-exchange/src/server/coin_interface"
	bitcoin "github.com/skycoin/skycoin-exchange/src/server/coin_interface/bitcoin"
	skycoin "github.com/skycoin/skycoin-exchange/src/server/coin_interface/skycoin"
	"github.com/skycoin/skycoin/src/util"
)

const (
	DETER_META_ID = iota
	DETER_META_INIT_SEED
	DETER_META_BTC_SEED
	DETER_META_SKY_SEED
	DETER_META_WALLET_TYPE
)

var DeterMetaStr = []string{
	DETER_META_ID:          "wallet_id",
	DETER_META_INIT_SEED:   "init_seed",
	DETER_META_BTC_SEED:    "btc_seed",
	DETER_META_SKY_SEED:    "sky_seed",
	DETER_META_WALLET_TYPE: "wallet_type",
}

// DeterministicWallet generate and store addresses for various coin types.
// This wallet is not thread safe.
type DeterministicWallet struct {
	ID             string                    // wallet id
	InitSeed       string                    // Init seed, used to recover the wallet.
	Seed           map[CoinType]string       // key: coin type, value: used to track the latset seed
	AddressEntries map[string][]AddressEntry // key: coin type, value: address entries.
	// addrLock       sync.Mutex // a lock, for protecting the writing, reading of the Addresses in wallet.
	// fileLock       sync.Mutex // lock for protecting wallet file.
}

// NewAddresses generate new addresses base on the coin type, and then store the address.
// NewAddress must be Sequential，cause the seed.
// this function is not thread safe, should not be used concrrently.
func (self *DeterministicWallet) NewAddresses(coinType CoinType, num int) ([]AddressEntry, error) {
	addrEntries := make([]AddressEntry, num)
	seed := self.Seed[coinType]
	switch coinType {
	case Bitcoin:
		if seed == self.InitSeed {
			sd, entries := bitcoin.GenerateAddresses([]byte(seed), num)
			self.Seed[Bitcoin] = sd
			addressEntryCopy(&addrEntries, entries)
		} else {
			s, err := hex.DecodeString(seed)
			if err != nil {
				return []AddressEntry{}, err
			}
			sd, entries := bitcoin.GenerateAddresses(s, num)
			self.Seed[Bitcoin] = sd
			addressEntryCopy(&addrEntries, entries)
		}
	case Skycoin:
		if seed == self.InitSeed {
			sd, entries := skycoin.GenerateAddresses([]byte(seed), num)
			self.Seed[Skycoin] = sd
			addressEntryCopy(&addrEntries, entries)
		} else {
			s, err := hex.DecodeString(seed)
			if err != nil {
				return []AddressEntry{}, err
			}
			sd, entries := skycoin.GenerateAddresses(s, num)
			self.Seed[Skycoin] = sd
			addressEntryCopy(&addrEntries, entries)
		}
	default:
		return addrEntries, fmt.Errorf("NewAddresses fail, unknow coin type:%d", coinType)
	}
	self.addAddresses(coinType, addrEntries)
	// save automaticaly after new addressess are added.
	self.save()
	return addrEntries, nil
}

func (self DeterministicWallet) GetCoinTypes() []CoinType {
	cts := []CoinType{}
	for ct, _ := range self.Seed {
		cts = append(cts, ct)
	}
	return cts
}

// Save the wallet
func (self *DeterministicWallet) save() error {
	// self.fileLock.Lock()
	// defer self.fileLock.Unlock()
	w := self.toWalletBase()
	return util.SaveJSON(filepath.Join(wltDir, self.ID), w, 0600)
}

func (self *DeterministicWallet) SetID(id string) {
	self.ID = id
}

func (self *DeterministicWallet) GetID() string {
	return self.ID
}

func (self DeterministicWallet) GetAddressEntries(bt CoinType) []AddressEntry {
	return self.AddressEntries[bt.String()]
}

func (self DeterministicWallet) GetAddresses(ct CoinType) []string {
	entries := self.AddressEntries[ct.String()]
	addrs := make([]string, len(entries))
	for i, entry := range entries {
		addrs[i] = entry.Address
	}
	return addrs
}

func (self *DeterministicWallet) toWalletBase() WalletBase {
	w := WalletBase{
		Meta: map[string]string{
			DeterMetaStr[DETER_META_ID]:          self.ID,
			DeterMetaStr[DETER_META_INIT_SEED]:   self.InitSeed,
			DeterMetaStr[DETER_META_BTC_SEED]:    self.Seed[Bitcoin],
			DeterMetaStr[DETER_META_SKY_SEED]:    self.Seed[Skycoin],
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

// GetAddressEntry get address entry by coin type and address.
func (self DeterministicWallet) GetAddressEntry(ct CoinType, addr string) (AddressEntry, error) {
	if aes, ok := self.AddressEntries[ct.String()]; ok {
		for _, a := range aes {
			if a.Address == addr {
				return a, nil
			}
		}
		return AddressEntry{}, errors.New("address not found")
	}
	return AddressEntry{}, errors.New("unknow coin type")
}

func newDeterministicWalletFromBase(w *WalletBase) (*DeterministicWallet, error) {
	var (
		id        string
		btc_seed  string
		sky_seed  string
		init_seed string
		ok        bool
	)

	if id, ok = w.Meta[DeterMetaStr[DETER_META_ID]]; !ok {
		return nil, errors.New("invalid wallet meta info, empty wallet_id")
	}

	if init_seed, ok = w.Meta[DeterMetaStr[DETER_META_INIT_SEED]]; !ok {
		return nil, errors.New("invalid wallet meta info, empty init seed")
	}

	if btc_seed, ok = w.Meta[DeterMetaStr[DETER_META_BTC_SEED]]; !ok {
		return nil, errors.New("invalid wallet meta info, empty btc seed")
	}

	if sky_seed, ok = w.Meta[DeterMetaStr[DETER_META_SKY_SEED]]; !ok {
		return nil, errors.New("invalid wallet meta info, empty sky seed")
	}

	wlt := &DeterministicWallet{
		ID: id,
		Seed: map[CoinType]string{
			Bitcoin: btc_seed, Skycoin: sky_seed},
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

func addressEntryCopy(dst *[]AddressEntry, src []coin_interface.AddressEntry) {
	for i, e := range src {
		(*dst)[i] = AddressEntry{
			Address: e.Address,
			Public:  e.Public,
			Secret:  e.Secret}
	}
}
