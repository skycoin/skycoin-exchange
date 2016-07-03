package account

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"time"

	"github.com/codahale/chacha20"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

type AccountID cipher.PubKey

type Accounter interface {
	GetWalletID() string                                                 // return the wallet id.
	GetAccountID() AccountID                                             // return the account id.
	GetNewAddress(ct wallet.CoinType) string                             // return new address for receiveing coins.
	GetBalance(ct wallet.CoinType) uint64                                // return the current balance.
	SetNonceKey(nk NonceKey)                                             // set the account's nonce key
	GetNonceKey() NonceKey                                               // get the account's nonce key
	Encrypt(r io.Reader) ([]byte, error)                                 // encrypt data
	Decrypt(r io.Reader) ([]byte, error)                                 // decrypt data
	IsExpired() bool                                                     // check if the nonce key is expired.
	GenerateWithdrawTx(coins uint64, ct wallet.CoinType) ([]byte, error) // account generate withdraw transaction.
}

// ExchangeAccount maintains the account state
type ExchangeAccount struct {
	ID          AccountID                  // account id
	balance     map[wallet.CoinType]uint64 // the balance should not be accessed directly.
	wltID       string                     // wallet used to maintain the address, UTXOs, balance, etc.
	nonceKey    NonceKey                   // key used to encrypt and decrypt data.
	balance_mtx sync.RWMutex               // mutex used to protect the balance's concurrent read and write.
	wlt_mtx     sync.Mutex                 // mutex used to protect the wallet's conncurrent read and write.
	key_mtx     sync.RWMutex               // mutex used to protect the nonce key's concurrent read and write.
}

// NonceKey used to encrypt and decrypt data,
// key will become invalid when time exceed the specific time.
type NonceKey struct {
	Nonce     []byte
	Key       []byte
	Expire_at time.Time
}

// newExchangeAccount helper function for generating and initialize ExchangeAccount
func newExchangeAccount(id AccountID, wltID string) ExchangeAccount {
	return ExchangeAccount{
		ID:    id,
		wltID: wltID,
		balance: map[wallet.CoinType]uint64{
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
func (self *ExchangeAccount) GetBalance(coinType wallet.CoinType) uint64 {
	self.balance_mtx.RLock()
	defer self.balance_mtx.RUnlock()
	return self.balance[coinType]
}

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

func (self *ExchangeAccount) SetNonceKey(nk NonceKey) {
	self.key_mtx.Lock()
	self.nonceKey = nk
	self.key_mtx.Unlock()
}

func (self ExchangeAccount) GetNonceKey() NonceKey {
	self.key_mtx.RLock()
	defer self.key_mtx.RUnlock()
	return self.nonceKey
}

// Encrypt encrypt the data from io.Reader with the local key.
func (self ExchangeAccount) Encrypt(r io.Reader) ([]byte, error) {
	return self.decOrEnc(r)
}

// Decrypt decrypt the data from io.Reader with the local key.
func (self *ExchangeAccount) Decrypt(r io.Reader) ([]byte, error) {
	return self.decOrEnc(r)
}

func (self ExchangeAccount) decOrEnc(r io.Reader) ([]byte, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return []byte{}, err
	}
	d := make([]byte, len(data))
	c, err := chacha20.New(self.nonceKey.Key, self.nonceKey.Nonce)
	if err != nil {
		return []byte{}, err
	}
	c.XORKeyStream(d, data)
	return d, nil
}

func (self ExchangeAccount) IsExpired() bool {
	return time.Now().Unix() >= self.nonceKey.Expire_at.Unix()
}

func (self ExchangeAccount) GenerateWithdrawTx(coins uint64, coinType wallet.CoinType) ([]byte, error) {
	// check if balance sufficient
	bla := self.GetBalance(coinType)
	if bla < coins {
		return []byte{}, errors.New("balance is not sufficient")
	}

	// get all unspent outputs of this account.

	return []byte{}, nil
}
