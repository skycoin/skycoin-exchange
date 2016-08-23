package wallet

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/skycoin/skycoin-exchange/src/coin"
)

type walletBase struct {
	ID             string              `json:"id"`                // wallet id
	InitSeed       string              `json:"init_seed"`         // Init seed, used to recover the wallet.
	Seed           string              `json:"seed"`              // used to track the latset seed
	AddressEntries []coin.AddressEntry `json:"entries,omitempty"` // address entries.
}

// GetID return wallet id.
func (wlt walletBase) GetID() string {
	return wlt.ID
}

// SetID set wallet id
func (wlt *walletBase) SetID(id string) {
	wlt.ID = id
}

// SetSeed initialize the wallet seed.
func (wlt *walletBase) SetSeed(seed string) {
	wlt.InitSeed = seed
	wlt.Seed = seed
}

// GetAddresses return all addresses in wallet.
func (wlt *walletBase) GetAddresses() []string {
	addrs := []string{}
	for _, e := range wlt.AddressEntries {
		addrs = append(addrs, e.Address)
	}
	return addrs
}

// GetKeypair get pub/sec key pair of specific address
func (wlt walletBase) GetKeypair(addr string) (string, string) {
	for _, e := range wlt.AddressEntries {
		if e.Address == addr {
			return e.Public, e.Secret
		}
	}
	return "", ""
}

// Save save the wallet
func (wlt *walletBase) Save(w io.Writer) error {
	d, err := json.MarshalIndent(wlt, "", "    ")
	if err != nil {
		return err
	}
	_, err = io.Copy(w, bytes.NewBuffer(d))
	return err
}

// Load load wallet from reader.
func (wlt *walletBase) Load(r io.Reader) error {
	return json.NewDecoder(r).Decode(wlt)
}

// Copy return the copy of self, for thread safe.
func (wlt walletBase) Copy() walletBase {
	return walletBase{
		ID:             wlt.ID,
		InitSeed:       wlt.InitSeed,
		Seed:           wlt.Seed,
		AddressEntries: wlt.AddressEntries,
	}
}
