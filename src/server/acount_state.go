package skycoin_exchange

import (
	"errors"
	"sync"

	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

type AccountID cipher.Address

// set of all users
type AccountManager struct {
	Accounts map[AccountID]*AccountState
	lck      sync.Mutex
	//AccountMap map[cipher.Address]uint64
}

//store state of user on server
type AccountState struct {
	ID      AccountID
	Balance map[string]uint64
	//Bitcoin balance in satoshis
	//Skycoin balance in drops

	//Inc1 uint64 //inc every write? Associated with local change
	//Inc2 uint64 //set to last change. Associatd with global event id
}

// GetAccount, return the copy value of speicific account.
func (self *AccountManager) GetAccount(addr AccountID) (AccountState, error) {
	self.lck.Lock()
	self.lck.Unlock()
	if account, ok := self.Accounts[addr]; ok {
		return *account, nil
	} else {
		return nil, errors.New("Account does not exist")
	}
}

func (self *AccountManager) CreateAccount(addr AccountID) (AccountState, error) {
	self.lck.Lock()
	self.lck.Unlock()
	if _, ok := self.Accounts[addr]; ok == true {
		return nil, errors.New("Account already exists")
	}

	act := newAccount(addr)
	self.Accounts[addr] = act
	return act, nil
}

func newAccount(addr AccountID) *AccountState {
	return &AccountState{
		AccountID: addr,
		Balance: map[string]uint64{
			wallet.Bitcoin.String(): 0,
			wallet.Skycoin.String(): 0,
		}}
}

//persistance to disc. Save as JSON
func (self *AccountManager) Save() {

}

func (self *AccountManager) Load() {
	//load accounts
}
