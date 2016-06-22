package skycoin_exchange

import (
	"errors"
	"fmt"
	"sync"

	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

type AccountID cipher.Address
type Balance uint64

// set of all users
type AccountManager struct {
	Accounts map[AccountID]*AccountState
	mtx      sync.Mutex
	//AccountMap map[cipher.Address]uint64
}

func NewAccountManager() *AccountManager {
	return &AccountManager{
		Accounts: make(map[AccountID]*AccountState)}
}

//store state of user on server
type AccountState struct {
	ID      AccountID
	balance map[wallet.CoinType]Balance // the balance should not be accessed directly.
	mtx     sync.Mutex
	//Bitcoin balance in satoshis
	//Skycoin balance in drops

	//Inc1 uint64 //inc every write? Associated with local change
	//Inc2 uint64 //set to last change. Associatd with global event id
}

// SetBalance update the balanace of specific coin.
func (self *AccountState) SetBalance(coinType wallet.CoinType, balance Balance) error {
	self.mtx.Lock()
	defer self.mtx.Unlock()
	if _, ok := self.balance[coinType]; !ok {
		return fmt.Errorf("the account does not have %s", coinType)
	}
	self.balance[coinType] = balance
	return nil
}

func (self *AccountState) GetBalance(coinType wallet.CoinType) (Balance, error) {
	self.mtx.Lock()
	defer self.mtx.Unlock()
	if b, ok := self.balance[coinType]; ok {
		return b, nil
	}
	return 0, fmt.Errorf("the account does not have %s", coinType)
}

// GetBalanceMap return the balance map.
func (self AccountState) GetBalanceMap() map[wallet.CoinType]Balance {
	return self.balance
}

// GetAccount, return the copy value of speicific account.
func (self *AccountManager) GetAccount(addr AccountID) (AccountState, error) {
	self.mtx.Lock()
	defer self.mtx.Unlock()
	if account, ok := self.Accounts[addr]; ok {
		return *account, nil
	} else {
		return AccountState{}, errors.New("Account does not exist")
	}
}

func (self *AccountManager) CreateAccount(addr AccountID) (AccountState, error) {
	self.mtx.Lock()
	defer self.mtx.Unlock()
	if _, ok := self.Accounts[addr]; ok == true {
		return AccountState{}, errors.New("Account already exists")
	}

	act := newAccount(addr)
	self.Accounts[addr] = &act
	return act, nil
}

func newAccount(id AccountID) AccountState {
	return AccountState{
		ID: id,
		balance: map[wallet.CoinType]Balance{
			wallet.Bitcoin: 0,
			wallet.Skycoin: 0,
		}}
}

//persistance to disc. Save as JSON
func (self *AccountManager) Save() {

}

func (self *AccountManager) Load() {
	//load accounts
}
