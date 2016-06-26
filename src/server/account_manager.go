package skycoin_exchange

import (
	"errors"
	"sync"

	"github.com/skycoin/skycoin-exchange/src/server/account"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
)

type AccountManager interface {
	CreateAccount() (Account, error)
	GetAccount(id AccountID) (Account, error)
	Save()
	Load()
}

// AccountManager manage all the accounts in the server.
type ExchangeAccountManager struct {
	Accounts map[AccountID]account.Account
	mtx      sync.RWMutex
	//AccountMap map[cipher.Address]uint64
}

// NewAccountManager
func NewAccountManager() AccountManager {
	return &ExchangeAccountManager{
		Accounts: make(map[AccountID]Account)}
}

func (self *ExchangeAccountManager) CreateAccount() (Account, error) {
	seed := cipher.SumSHA256(cipher.RandByte(1024)).Hex()
	p, _ := cipher.GenerateDeterministicKeyPair([]byte(seed))
	wlt, err := wallet.NewWallet(seed)
	if err != nil {
		return nil, err
	}

	act := newAccountState(p, wlt.GetID)

	self.mtx.Lock()
	// TODO: check duplicate account.

	// add the account.
	self.Accounts[p] = &act
	self.mtx.Unlock()
	return &act, nil
}

// GetAccount return the account of specific id.
func (self *ExchangeAccountManager) GetAccount(id AccountID) (Account, error) {
	self.mtx.RLock()
	defer self.mtx.RUnlock()
	if account, ok := self.Accounts[id]; ok {
		return account, nil
	} else {
		return nil{}, errors.New("account does not exist")
	}
}

//persistance to disc. Save as JSON
func (self *ExchangeAccountManager) Save() {

}

func (self *ExchangeAccountManager) Load() {
	//load accounts
}
