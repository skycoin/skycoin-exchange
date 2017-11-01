package account

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	logging "github.com/op/go-logging"
	"github.com/skycoin/skycoin/src/util/file"
)

var (
	acntDir  = filepath.Join(file.UserHome(), ".skycoin-exchange/account")
	acntName = "account.data"
	logger   = logging.MustGetLogger("exchange.account")
)

type Accounter interface {
	GetID() string                            // return the account id.
	GetBalance(ct string) uint64              // return the account's Balance.
	AddDepositAddress(ct string, addr string) // add the deposit address to the account.
	DecreaseBalance(ct string, amt uint64) error
	IncreaseBalance(ct string, amt uint64) error
	SetBalance(cp string, amt uint64) error
}

// ExchangeAccount maintains the account state
type ExchangeAccount struct {
	ID          string              // account id
	Balance     map[string]uint64   // the Balance should not be accessed directly.
	Addresses   map[string][]string // deposit addresses
	addr_mtx    sync.Mutex
	balance_mtx sync.RWMutex // mutex used to protect the Balance's concurrent read and write.
}

type exchgAcntJson struct {
	ID        string              `json:"id"`
	Balance   map[string]uint64   `json:"balance"`
	Addresses map[string][]string `json:"addresses"`
}

// InitDir init the account storage file path.
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
		Balance: map[string]uint64{
			"skycoin": 0,
			"bitcoin": 0,
		},
		Addresses: make(map[string][]string),
	}
}

func (self ExchangeAccount) GetID() string {
	return self.ID
}

// Get the current recored Balance.
func (self *ExchangeAccount) GetBalance(coinType string) uint64 {
	self.balance_mtx.RLock()
	defer self.balance_mtx.RUnlock()
	return self.Balance[coinType]
}

func (self *ExchangeAccount) AddDepositAddress(coinType string, addr string) {
	self.addr_mtx.Lock()
	self.Addresses[coinType] = append(self.Addresses[coinType], addr)
	self.addr_mtx.Unlock()
}

// SetBalance update the balanace of specific coin.
func (self *ExchangeAccount) SetBalance(cp string, amt uint64) error {
	self.balance_mtx.Lock()
	defer self.balance_mtx.Unlock()
	if _, ok := self.Balance[cp]; !ok {
		return fmt.Errorf("the account does not have %s", cp)
	}
	self.Balance[cp] = amt
	return nil
}

func (self *ExchangeAccount) DecreaseBalance(ct string, amt uint64) error {
	self.balance_mtx.Lock()
	defer self.balance_mtx.Unlock()
	if _, ok := self.Balance[ct]; !ok {
		return errors.New("unknow coin type")
	}
	if self.Balance[ct] < amt {
		logger.Debug("balance:%d require:%d", self.Balance[ct], amt)
		return errors.New("account balance is not sufficient")
	}

	self.Balance[ct] -= amt
	return nil
}

func (self *ExchangeAccount) IncreaseBalance(ct string, amt uint64) error {
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
		eaj.Balance[ct] = bal
	}

	for ct, addrs := range self.Addresses {
		eaj.Addresses[ct] = append(eaj.Addresses[ct], addrs...)
	}
	return eaj
}

func (self exchgAcntJson) ToExchgAcnt() *ExchangeAccount {
	// pk := cipher.PubKey{}
	// copy(pk[:], self.ID[0:33])
	at := ExchangeAccount{
		ID:        self.ID,
		Balance:   make(map[string]uint64),
		Addresses: make(map[string][]string),
	}

	// convert balance.
	for ct, bal := range self.Balance {
		at.Balance[ct] = bal
	}

	// convert address
	for ct, addrs := range self.Addresses {
		at.Addresses[ct] = append(at.Addresses[ct], addrs...)
	}
	return &at
}
