package skycoin_exchange

import (
	"errors"
	"sync"

	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

type AccountManager interface {
	CreateAccount() (Accounter, cipher.SecKey, error)
	GetAccount(id AccountID) (Accounter, error)
	Save()
	Load()
}

// AccountManager manage all the accounts in the server.
type ExchangeAccountManager struct {
	Accounts map[AccountID]Accounter
	mtx      sync.RWMutex
	//AccountMap map[cipher.Address]uint64
}

// NewAccountManager
func NewExchangeAccountManager() AccountManager {
	return &ExchangeAccountManager{
		Accounts: make(map[AccountID]Accounter)}
}

// CreateAccount create new account, and bind a new wallet to this account,
// generate pubkey/seckey pair, the pubkey will be stored in the account, and the
// seckey will be returned.
// Notice: for client, this is the only chance to get seckey!
func (self *ExchangeAccountManager) CreateAccount() (Accounter, cipher.SecKey, error) {
	seed := cipher.SumSHA256(cipher.RandByte(1024)).Hex()
	wlt, err := wallet.NewWallet(seed)
	if err != nil {
		return nil, cipher.SecKey{}, err
	}

	p, s := cipher.GenerateDeterministicKeyPair([]byte(seed))
	act := newExchangeAccount(AccountID(p), wlt.GetID())

	self.mtx.Lock()
	// TODO: check duplicate account.

	// add the account.
	self.Accounts[AccountID(p)] = &act
	self.mtx.Unlock()
	return &act, s, nil
}

// GetAccount return the account of specific id.
func (self *ExchangeAccountManager) GetAccount(id AccountID) (Accounter, error) {
	self.mtx.RLock()
	defer self.mtx.RUnlock()
	if account, ok := self.Accounts[id]; ok {
		return account, nil
	} else {
		return nil, errors.New("account does not exist")
	}
}

//persistance to disc. Save as JSON
func (self *ExchangeAccountManager) Save() {

}

func (self *ExchangeAccountManager) Load() {
	//load accounts
}
