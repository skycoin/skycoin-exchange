package wallet

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/util"
)

type CoinType int8
type WalletType int8

// type MetaInfo int8

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
