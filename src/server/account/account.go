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
	GetBalance(ct wallet.CoinType) uint64              // return the account's balance.
	AddDepositAddress(ct wallet.CoinType, addr string) // add the deposit address to the account.
	DecreaseBalance(amt uint64) error
}

// ExchangeAccount maintains the account state
type ExchangeAccount struct {
	ID          AccountID                  // account id
	balance     map[wallet.CoinType]uint64 // the balance should not be accessed directly.
	balance_mtx sync.RWMutex               // mutex used to protect the balance's concurrent read and write.
	addresses   map[wallet.CoinType][]string
	addr_mtx    sync.Mutex
}

type addrBalance struct {
	Addr    string
	Balance uint64
}

type byBalance []addrBalance

func (bb byBalance) Len() int           { return len(bb) }
func (bb byBalance) Swap(i, j int)      { bb[i], bb[j] = bb[j], bb[i] }
func (bb byBalance) Less(i, j int) bool { return bb[i].Balance < bb[j].Balance }

// newExchangeAccount helper function for generating and initialize ExchangeAccount
func newExchangeAccount(id AccountID) ExchangeAccount {
	return ExchangeAccount{
		ID: id,
		balance: map[wallet.CoinType]uint64{
			wallet.Bitcoin: 0,
			wallet.Skycoin: 0,
		}}
}

func (self ExchangeAccount) GetID() AccountID {
	return self.ID
}

// GetNewAddress generate new address for this account.
// func (self *ExchangeAccount) GetNewAddress(ct wallet.CoinType) string {
// 	// get the wallet.
// 	wlt, err := wallet.GetWallet(self.wltID)
// 	if err != nil {
// 		panic(fmt.Sprintf("account get wallet faild, wallet id:%s", self.wltID))
// 	}
//
// 	self.wlt_mtx.Lock()
// 	defer self.wlt_mtx.Unlock()
// 	addr, err := wlt.NewAddresses(ct, 1)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return addr[0].Address
// }

// Get the current recored balance.
func (self *ExchangeAccount) GetBalance(coinType wallet.CoinType) uint64 {
	self.balance_mtx.RLock()
	defer self.balance_mtx.RUnlock()
	return self.balance[coinType]
}

func (self *ExchangeAccount) AddDepositAddress(coinType wallet.CoinType, addr string) {
	self.addr_mtx.Lock()
	self.addresses[coinType] = append(self.addresses[coinType], addr)
	self.addr_mtx.Unlock()
}

// func (self ExchangeAccount) GetAddressBalance(addr string) (uint64, error) {
// 	return bitcoin.GetBalance([]string{addr})
// }

// SetBalance update the balanace of specific coin.
func (self *ExchangeAccount) setBalance(coinType wallet.CoinType, balance uint64) error {
	self.balance_mtx.Lock()
	defer self.balance_mtx.Unlock()
	if _, ok := self.balance[coinType]; !ok {
		return fmt.Errorf("the account does not have %s", coinType)
	}
	self.balance[coinType] = balance
	return nil
}

func (self *ExchangeAccount) DecreaseBalance(ct wallet.CoinType, amt uint64) error {
	self.balance_mtx.Lock()
	defer self.balance_mtx.Unlock()
	if _, ok := self.balance[ct]; !ok {
		return errors.New("unknow coin type")
	}
	if self.balance[ct] < amt {
		return errors.New("account balance is not sufficient")
	}

	self.balance[ct] = self.balance[ct] - amt
	return nil
}
