package wallet

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
)

type CoinType int8
type WalletType int8
type WalletID string

var (
	// GWallets Wallets
	wltDir string = filepath.Join(util.UserHome(), ".skycoin-exchange/wallets")
)

const (
	Bitcoin CoinType = iota
	Skycoin
	// Shellcoin
	// Ethereum
	// other coins...
)

const (
	Deterministic WalletType = iota // default wallet type
)

var walletTypeStr = []string{
	Deterministic: "deterministic",
}

var coinStr = []string{
	Bitcoin: "bitcoin",
	Skycoin: "skycoin",
}

func (c CoinType) String() string {
	switch c {
	case Bitcoin:
		return coinStr[c]
	case Skycoin:
		return coinStr[c]
	default:
		// return fmt.Sprintf("unknow coin type:%d", c)
		panic(fmt.Sprintf("unknow coin type:%d", c))
	}
}

func ConvertCoinType(ct string) (CoinType, error) {
	switch ct {
	case "bitcoin":
		return Bitcoin, nil
	case "skycoin":
		return Skycoin, nil
	default:
		return -1, fmt.Errorf("unknow coin type:%s", ct)
	}
}

func (w WalletType) String() string {
	switch w {
	case Deterministic:
		return walletTypeStr[w]
	default:
		// return fmt.Sprintf("unknow wallet type:%d", w)
		panic(fmt.Sprintf("unknow wallet type:%d", w))
	}
}

type Wallet interface {
	SetID(id string)
	GetID() string
	NewAddresses(coinType CoinType, num int) ([]AddressEntry, error)
	GetAddressEntries(coinType CoinType) []AddressEntry
	GetAddressEntry(coinType CoinType, addr string) (AddressEntry, error)
}

// WalletBase, used to serialise wallet into json, and unserialise wallet from json.
type WalletBase struct {
	Meta           map[string]string         `json:"meta"`
	AddressEntries map[string][]AddressEntry `json:"addresses"` // key is coin type, value is address entry list.
}

type AddressEntry struct {
	Address string `json:"address"`
	Public  string `json:"pubkey"`
	Secret  string `json:"seckey"`
}

// New create a new wallet, and save it to local disk.
func New(name string, wltType WalletType, seed string) (Wallet, error) {
	if seed == "" {
		seed = cipher.SumSHA256(cipher.RandByte(1024)).Hex()
	}

	switch wltType {
	case Deterministic:
		wlt := &DeterministicWallet{
			ID:             name,
			Seed:           map[CoinType]string{Bitcoin: seed, Skycoin: seed},
			InitSeed:       seed,
			AddressEntries: make(map[string][]AddressEntry)}
		if err := wlt.save(); err != nil {
			return nil, err
		}
		return wlt, nil
	default:
		return nil, fmt.Errorf("newWallet fail, unknow wallet type:%d", wltType)
	}
}

// Load load wallet from specific local disk.
func Load(name string) (Wallet, error) {
	p := filepath.Join(wltDir, name)
	glog.Info("load wallet")
	w, err := loadWalletFromFile(p)
	if err != nil {
		return nil, err
	}

	concretWlt, err := w.newConcretWallet()
	if err != nil {
		return nil, err
	}
	concretWlt.SetID(name)
	return concretWlt, nil
}

func IsExist(wltName string) bool {
	p := filepath.Join(wltDir, wltName)
	_, err := os.Stat(p)
	return !os.IsNotExist(err)
}

func InitDir(dir string) {
	// check if the wallet dir is exist.
	wltDir = dir
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// create the dir
		if err := os.MkdirAll(dir, 0700); err != nil {
			panic(err)
		}
	}
}

func loadWalletFromFile(filename string) (WalletBase, error) {
	w := WalletBase{}
	err := util.LoadJSON(filename, &w)
	if err != nil {
		return WalletBase{}, err
	}
	return w, nil
}

// newConcretWallet, create concret wallet base on the wallet type.
func (self *WalletBase) newConcretWallet() (Wallet, error) {
	if wltType, ok := self.Meta["wallet_type"]; ok {
		switch wltType {
		case Deterministic.String():
			wlt, err := newDeterministicWalletFromBase(self)
			if err != nil {
				return nil, err
			}
			return wlt, nil
		default:
			return nil, fmt.Errorf("unknow wallet type:%s", self.Meta["wallet_type"])
		}
	}
	return nil, errors.New("invalide wallet meta info from wallet file")
}
