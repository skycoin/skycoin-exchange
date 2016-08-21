package wallet

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/skycoin/skycoin-exchange/src/server/coin"
	"github.com/skycoin/skycoin/src/util"
)

type walletBase struct {
	ID             string              `json:"id"` // wallet id
	Type           string              `json:"type"`
	InitSeed       string              `json:"init_seed"`         // Init seed, used to recover the wallet.
	Seed           string              `json:"seed"`              // used to track the latset seed
	AddressEntries []coin.AddressEntry `json:"entries,omitempty"` // address entries.
}

// GetID return wallet id.
func (wlt walletBase) GetID() string {
	return wlt.ID
}

// NewAddresses generate bitcoin addresses.
func (wlt *walletBase) NewAddresses(num int) ([]coin.AddressEntry, error) {
	return []coin.AddressEntry{}, nil
}

// GetAddresses return all wallet addresses.
func (wlt *walletBase) GetAddresses() []string {
	return []string{}
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
func (wlt *walletBase) Save() error {
	fileName := wlt.ID + "." + wltExt
	return util.SaveJSON(filepath.Join(wltDir, fileName), wlt, 0600)
}

// Load load wallet from reader.
func (wlt *walletBase) Load(r io.Reader) error {
	d, err := ioutil.ReadAll(r)
	if err != nil {
		return nil
	}
	return json.Unmarshal(d, wlt)
}

// Clear remove wallet file from local disk.
func (wlt *walletBase) Clear() error {
	path := filepath.Join(wltDir, wlt.ID)
	return os.RemoveAll(path)
}
