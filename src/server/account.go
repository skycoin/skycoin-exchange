package skycoin_exchange

import (
	"fmt"
	"sync"

	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

type AccountID cipher.PubKey
type Balance uint64 // satoshis

type Accounter interface {
	GetWalletID() string                            // return the wallet id
	GetAccountID() AccountID                        // return the account id
	GetNewAddress(ct wallet.CoinType) string        // return new address for receiveing coins
	GetBalance(ct wallet.CoinType) (Balance, error) // return the current balance.
	// GetUnspentOutput(ct wallet.CoinType, minConf int) // return all unspent output of this account that confirms minConf times.
}

// ExchangeAccount maintains the account state
type ExchangeAccount struct {
	ID          AccountID                   // account id
	balance     map[wallet.CoinType]Balance // the balance should not be accessed directly.
	wltID       string                      // wallet used to maintain the address, UTXOs, balance, etc.
	balance_mtx sync.RWMutex                // mutex used to protect the balance's concurrent read and write.
	wlt_mtx     sync.Mutex                  // mutex used to protect the wallet's conncurrent read and write.
}

// newExchangeAccount helper function for generating and initialize ExchangeAccount
func newExchangeAccount(id AccountID, wltID string) ExchangeAccount {
	return ExchangeAccount{
		ID:    id,
		wltID: wltID,
		balance: map[wallet.CoinType]Balance{
			wallet.Bitcoin: 0,
			wallet.Skycoin: 0,
		}}
}

func (self ExchangeAccount) GetAccountID() AccountID {
	return self.ID
}

func (self ExchangeAccount) GetWalletID() string {
	return self.wltID
}

// GetNewAddress generate new address for this account.
func (self *ExchangeAccount) GetNewAddress(ct wallet.CoinType) string {
	// get the wallet.
	wlt, err := wallet.GetWallet(self.wltID)
	if err != nil {
		panic(fmt.Sprintf("account get wallet faild, wallet id:%s", self.wltID))
	}

	self.wlt_mtx.Lock()
	defer self.wlt_mtx.Unlock()
	addr, err := wlt.NewAddresses(ct, 1)
	if err != nil {
		panic(err)
	}
	return addr[0].Address
}

// Get the current recored balance.
func (self *ExchangeAccount) GetBalance(coinType wallet.CoinType) (Balance, error) {
	self.balance_mtx.RLock()
	defer self.balance_mtx.RUnlock()
	if b, ok := self.balance[coinType]; ok {
		return b, nil
	}
	return 0, fmt.Errorf("the account does not have %s", coinType)
}

// SetBalance update the balanace of specific coin.
func (self *ExchangeAccount) setBalance(coinType wallet.CoinType, balance Balance) error {
	self.balance_mtx.Lock()
	defer self.balance_mtx.Unlock()
	if _, ok := self.balance[coinType]; !ok {
		return fmt.Errorf("the account does not have %s", coinType)
	}
	self.balance[coinType] = balance
	return nil
}
