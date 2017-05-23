package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/skycoin/skycoin-exchange/src/coin"
)

// Wallet wallet struct
type Wallet struct {
	ID             string              `json:"id"`                // wallet id
	InitSeed       string              `json:"init_seed"`         // Init seed, used to recover the wallet.
	Seed           string              `json:"seed"`              // used to track the latset seed
	AddressEntries []coin.AddressEntry `json:"entries,omitempty"` // address entries.
	Type           string              `json:"type"`              // wallet type
}

// GetID return wallet id.
func (wlt Wallet) GetID() string {
	return wlt.ID
}

// SetID set wallet id
func (wlt *Wallet) SetID(id string) {
	wlt.ID = id
}

// SetSeed initialize the wallet seed.
func (wlt *Wallet) SetSeed(seed string) {
	wlt.InitSeed = seed
	wlt.Seed = seed
}

// GetAddresses return all addresses in wallet.
func (wlt *Wallet) GetAddresses() []string {
	addrs := []string{}
	for _, e := range wlt.AddressEntries {
		addrs = append(addrs, e.Address)
	}
	return addrs
}

// GetKeypair get pub/sec key pair of specific address
func (wlt Wallet) GetKeypair(addr string) (string, string, error) {
	for _, e := range wlt.AddressEntries {
		if e.Address == addr {
			return e.Public, e.Secret, nil
		}
	}
	return "", "", fmt.Errorf("%s addr does not exist in wallet", addr)
}

// Save save the wallet
func (wlt *Wallet) Save(w io.Writer) error {
	d, err := json.MarshalIndent(wlt, "", "    ")
	if err != nil {
		return err
	}
	_, err = io.Copy(w, bytes.NewBuffer(d))
	return err
}

// Load load wallet from reader.
func (wlt *Wallet) Load(r io.Reader) error {
	return json.NewDecoder(r).Decode(wlt)
}

// GetType returns the wallet type
func (wlt *Wallet) GetType() string {
	return wlt.Type
}

// Copy return the copy of self, for thread safe.
func (wlt Wallet) Copy() Wallet {
	return Wallet{
		ID:             wlt.ID,
		InitSeed:       wlt.InitSeed,
		Seed:           wlt.Seed,
		AddressEntries: wlt.AddressEntries,
	}
}
