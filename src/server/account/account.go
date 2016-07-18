package account

import (
	"errors"
	"fmt"
	"sync"

	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

type AccountID cipher.PubKey

type Accounter interface {
	GetID() AccountID                                  // return the account id.
	GetBalance(ct wallet.CoinType) uint64              // return the account's Balance.
	AddDepositAddress(ct wallet.CoinType, addr string) // add the deposit address to the account.
	DecreaseBalance(ct wallet.CoinType, amt uint64) error
	IncreaseBalance(ct wallet.CoinType, amt uint64) error
}

// ExchangeAccount maintains the account state
type ExchangeAccount struct {
	ID          AccountID                    `json:"id"`        // account id
	Balance     map[wallet.CoinType]uint64   `json:"balance"`   // the Balance should not be accessed directly.
	Addresses   map[wallet.CoinType][]string `json:"addresses"` //
	addr_mtx    sync.Mutex
	balance_mtx sync.RWMutex // mutex used to protect the Balance's concurrent read and write.
}

// newExchangeAccount helper function for generating and initialize ExchangeAccount
func newExchangeAccount(id AccountID) ExchangeAccount {
	return ExchangeAccount{
		ID: id,
		Balance: map[wallet.CoinType]uint64{
			wallet.Bitcoin: 0,
			wallet.Skycoin: 0,
		},
		Addresses: make(map[wallet.CoinType][]string),
	}
}

func (self ExchangeAccount) GetID() AccountID {
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
		return errors.New("account Balance is not sufficient")
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

func init() {
	// init the account dir of server.
}
