package skycoin_exchange

import (
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
