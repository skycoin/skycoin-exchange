package skycoin_exchange

import (
	"sync"
	"testing"

	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/stretchr/testify/assert"
)

type AccountHelper struct {
	AntManager *AccountManager
	Account    *AccountState
}

func newAccountHelper(t *testing.T) AccountHelper {
	ah := AccountHelper{
		AntManager: NewAccountManager()}
	ant, err := ah.AntManager.CreateAccount(addressForTest())
	assert.Nil(t, err)
	ah.Account = &ant
	return ah
}

// test for account_state.
func addressForTest() AccountID {
	addr := "1Cd86LYDEicDq4esruFEFuM5b4wLFGU6Xc"
	a := cipher.BitcoinMustDecodeBase58Address(addr)
	return AccountID(a)
}

func TestCreateAccount(t *testing.T) {
	ah := newAccountHelper(t)
	assert.Equal(t, ah.Account.balance[wallet.Bitcoin], Balance(0))
	assert.Equal(t, ah.Account.balance[wallet.Skycoin], Balance(0))
}

func TestCreateAccountConcurrent(t *testing.T) {
	btcAddrs := []string{
		"18Kc6k8CqEM3b9Uc4KkKDdofRi4e6uJ1aG",
		"1ArpWzE5nfnE7WJ9RBJ5vSkV3w4pdBYPoo",
		"1JxDpnNLbW5SyhoG4H6CxcZVgKciTQPUHx",
		"18djdctuWsybFG1cCWej2XJE4eJcw5guya",
		"1HFfX1oHkuXLxvvPgTFtiRdGP85PXEF8kk",
		"13JqGcbD5yEYRWW1vpg86zPZWAqEVZW1F9",
		"17xD32p5dDfPyUfV3SMQyR6Cfj2ntaDWPa",
		"15hYoHsPdbFhua6C7QhvRViRfGUsYUpotZ",
		"15Krv5RrgXD4TvkedMovxRaYDmeAPxF4qP"}

	ah := newAccountHelper(t)
	wg := sync.WaitGroup{}
	for _, addr := range btcAddrs {
		wg.Add(1)
		go func(wg *sync.WaitGroup, addr string) {
			defer wg.Done()
			_, err := ah.AntManager.CreateAccount(AccountID(cipher.BitcoinMustDecodeBase58Address(addr)))
			assert.Nil(t, err)
		}(&wg, addr)
	}
	wg.Wait()
}

func TestCreateDupAccount(t *testing.T) {
	ah := newAccountHelper(t)
	_, err := ah.AntManager.CreateAccount(addressForTest())
	assert.NotNil(t, err)
}

func TestGetAccount(t *testing.T) {
	ah := newAccountHelper(t)
	ant, err := ah.AntManager.GetAccount(addressForTest())
	assert.Nil(t, err)
	assert.Equal(t, ant.balance[wallet.Bitcoin], Balance(0))
	assert.Equal(t, ant.balance[wallet.Skycoin], Balance(0))
}

func TestSetBalance(t *testing.T) {
	ah := newAccountHelper(t)
	ah.Account.SetBalance(wallet.Bitcoin, Balance(10))
	ah.Account.SetBalance(wallet.Skycoin, Balance(20))
	assert.Equal(t, ah.Account.balance[wallet.Bitcoin], Balance(10))
	assert.Equal(t, ah.Account.balance[wallet.Skycoin], Balance(20))
}

func TestGetBalance(t *testing.T) {
	ah := newAccountHelper(t)
	ah.Account.SetBalance(wallet.Bitcoin, Balance(10))
	ah.Account.SetBalance(wallet.Skycoin, Balance(20))

	bb, err := ah.Account.GetBalance(wallet.Bitcoin)
	assert.Nil(t, err)
	assert.Equal(t, bb, Balance(10))
	sb, err := ah.Account.GetBalance(wallet.Skycoin)
	assert.Nil(t, err)
	assert.Equal(t, sb, Balance(20))
}
