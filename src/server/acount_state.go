package skycoin_exchange

import (
	"errors"

	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

type AccountIdentifier cipher.Address

// set of all users
type AccountManager struct {
	Accounts map[AccountIdentifier]*AccountState
	//AccountMap map[cipher.Address]uint64
}

//store state of user on server
type AccountState struct {
	//AccountId  uint64
	AccountID AccountIdentifier

	Balance map[string]uint64
	//Bitcoin balance in satoshis
	//Skycoin balance in drops

	//Inc1 uint64 //inc every write? Associated with local change
	//Inc2 uint64 //set to last change. Associatd with global event id
}

func (self *AccountManager) GetAccount(addr AccountIdentifier) (*AccountState, error) {
	if account, ok := self.Accounts[addr]; ok {
		return account, nil
	} else {
		return nil, errors.New("Account does not exist")
	}
}

func (self *AccountManager) CreateAccount(addr AccountIdentifier) (*AccountState, error) {
	if _, ok := self.Accounts[addr]; ok == true {
		return nil, errors.New("Account already exists")
	}

	act := newAccount(addr)
	self.Accounts[addr] = act
	return act, nil
}

func newAccount(addr AccountIdentifier) *AccountState {
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
