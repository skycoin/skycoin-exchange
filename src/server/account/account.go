package account

import (
	"errors"
	"fmt"
	"sync"
	"time"

	bitcoin "github.com/skycoin/skycoin-exchange/src/server/coin_interface/bitcoin"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

type AccountID cipher.PubKey

type Accounter interface {
	GetWalletID() string                                                               // return the wallet id.
	GetID() AccountID                                                                  // return the account id.
	GetNewAddress(ct wallet.CoinType) string                                           // return new address for receiveing coins.
	GetBalance(ct wallet.CoinType) uint64                                              // return the current balance.
	GenerateWithdrawTx(ct wallet.CoinType, outAddrs []bitcoin.OutAddr) ([]byte, error) // account generate withdraw transaction.
}

// ExchangeAccount maintains the account state
type ExchangeAccount struct {
	ID          AccountID                  // account id
	balance     map[wallet.CoinType]uint64 // the balance should not be accessed directly.
	wltID       string                     // wallet used to maintain the address, UTXOs, balance, etc.
	balance_mtx sync.RWMutex               // mutex used to protect the balance's concurrent read and write.
	wlt_mtx     sync.Mutex                 // mutex used to protect the wallet's conncurrent read and write.
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

func (self ExchangeAccount) GetID() AccountID {
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

// GenerateWithdrawTx
func (self ExchangeAccount) GenerateWithdrawTx(coinType wallet.CoinType, outAddrs []bitcoin.OutAddr) ([]byte, error) {
	// check if balance sufficient
	bla := self.GetBalance(coinType)
	var allCoins int64
	for _, out := range outAddrs {
		allCoins += out.Value
	}

	if bla < uint64(allCoins) {
		return []byte{}, errors.New("balance is not sufficient")
	}

	addrs, err := self.getAddressEntries(coinType)
	if err != nil {
		return []byte{}, err
	}

	// get utxos
	utxos := []bitcoin.UtxoWithkey{}
	switch coinType {
	case wallet.Bitcoin:
		for _, addrEntry := range addrs {
			us := bitcoin.GetUnspentOutputsBlkChnInfo(addrEntry.Address)
			for _, u := range us {
				usk := bitcoin.BlkchnUtxoWithkey{
					BlkChnUtxo: u,
					Privkey:    addrEntry.Secret,
				}
				utxos = append(utxos, usk)
			}
		}
		return bitcoin.NewTransaction(utxos, outAddrs)
	default:
		return []byte{}, errors.New("unknow coin type")
	}
}

func (self ExchangeAccount) getAddressEntries(coinType wallet.CoinType) ([]wallet.AddressEntry, error) {
	// get address list of this account
	wlt, err := wallet.GetWallet(self.wltID)
	if err != nil {
		return []wallet.AddressEntry{}, fmt.Errorf("account get wallet faild, wallet id:%s", self.wltID)
	}
	addresses := wlt.GetAddressEntries(coinType)
	return addresses, nil
}
