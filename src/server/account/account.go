package account

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/op/go-logging.v1"

	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/util"
)

var (
	acntDir  string = filepath.Join(util.UserHome(), ".skycoin-exchange/account")
	acntName string = "account.data"
	logger          = logging.MustGetLogger("exchange.account")
)

// type AccountID cipher.PubKey

type Accounter interface {
	GetID() string                                     // return the account id.
	GetBalance(ct wallet.CoinType) uint64              // return the account's Balance.
	AddDepositAddress(ct wallet.CoinType, addr string) // add the deposit address to the account.
	DecreaseBalance(ct wallet.CoinType, amt uint64) error
	IncreaseBalance(ct wallet.CoinType, amt uint64) error
}

// ExchangeAccount maintains the account state
type ExchangeAccount struct {
	ID          string                       // account id
	Balance     map[wallet.CoinType]uint64   // the Balance should not be accessed directly.
	Addresses   map[wallet.CoinType][]string // deposit addresses
	addr_mtx    sync.Mutex
	balance_mtx sync.RWMutex // mutex used to protect the Balance's concurrent read and write.
}

type exchgAcntJson struct {
	ID        string              `json:"id"`
	Balance   map[string]uint64   `json:"balance"`
	Addresses map[string][]string `json:"addresses"`
}

func InitDir(path string) {
	if path == "" {
		path = acntDir
	} else {
		acntDir = path
	}
	// create the account dir if not exist.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			panic(err)
		}
	}
}

// newExchangeAccount helper function for generating and initialize ExchangeAccount
func newExchangeAccount(id string) ExchangeAccount {
	return ExchangeAccount{
		ID: id,
		Balance: map[wallet.CoinType]uint64{
			wallet.Bitcoin: 0,
			wallet.Skycoin: 0,
		},
		Addresses: make(map[wallet.CoinType][]string),
	}
}

func (self ExchangeAccount) GetID() string {
	return self.ID
}

// Get the current recored Balance.
func (self *ExchangeAccount) GetBalance(coinType wallet.CoinType) uint64 {
	self.balance_mtx.RLock()
	defer self.balance_mtx.RUnlock()
	return self.Balance[coinType]
}

func (self *ExchangeAccount) AddDepositAddress(coinType wallet.CoinType, addr string) {
	self.addr_mtx.Lock()
	self.Addresses[coinType] = append(self.Addresses[coinType], addr)
	self.addr_mtx.Unlock()
}

// SetBalance update the balanace of specific coin.
func (self *ExchangeAccount) setBalance(coinType wallet.CoinType, balance uint64) error {
	self.balance_mtx.Lock()
	defer self.balance_mtx.Unlock()
	if _, ok := self.Balance[coinType]; !ok {
		return fmt.Errorf("the account does not have %s", coinType)
	}
	self.Balance[coinType] = balance
	return nil
}

func (self *ExchangeAccount) DecreaseBalance(ct wallet.CoinType, amt uint64) error {
	self.balance_mtx.Lock()
	defer self.balance_mtx.Unlock()
	if _, ok := self.Balance[ct]; !ok {
		return errors.New("unknow coin type")
	}
	if self.Balance[ct] < amt {
		logger.Debug("balance:", self.Balance[ct], " amt:", amt)
		return errors.New("account balance is not sufficient")
	}

	self.Balance[ct] -= amt
	return nil
}

func (self *ExchangeAccount) IncreaseBalance(ct wallet.CoinType, amt uint64) error {
	self.balance_mtx.Lock()
	defer self.balance_mtx.Unlock()
	if _, ok := self.Balance[ct]; !ok {
		return errors.New("unknow coin type")
	}

	self.Balance[ct] += amt
	return nil
}

func (self ExchangeAccount) ToMarshalable() exchgAcntJson {
	eaj := exchgAcntJson{
		ID:        self.ID,
		Balance:   make(map[string]uint64),
		Addresses: make(map[string][]string),
	}

	for ct, bal := range self.Balance {
		eaj.Balance[ct.String()] = bal
	}

	for ct, addrs := range self.Addresses {
		eaj.Addresses[ct.String()] = append(eaj.Addresses[ct.String()], addrs...)
	}
	return eaj
}

func (self exchgAcntJson) ToExchgAcnt() *ExchangeAccount {
	// pk := cipher.PubKey{}
	// copy(pk[:], self.ID[0:33])
	at := ExchangeAccount{
		ID:        self.ID,
		Balance:   make(map[wallet.CoinType]uint64),
		Addresses: make(map[wallet.CoinType][]string),
	}

	// convert balance.
	for ct, bal := range self.Balance {
		t, err := wallet.ConvertCoinType(ct)
		if err != nil {
			panic(err)
		}
		at.Balance[t] = bal
	}

	// convert address
	for ct, addrs := range self.Addresses {
		t, err := wallet.ConvertCoinType(ct)
		if err != nil {
			panic(err)
		}
		at.Addresses[t] = append(at.Addresses[t], addrs...)
	}
	return &at
}
