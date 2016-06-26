package skycoin_exchange

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/stretchr/testify/assert"
)

type funHandler func(ema AccountManager)

func PrepareFunc(f funHandler) {
	// create
	ema := NewExchangeAccountManager()
	f(ema)
	// clean wallet data.
	removeContents(wallet.GetWalletDatadir())
	// clear GWallet
	wallet.Reload()
}

func TestCreateAccountConcurrent(t *testing.T) {
	PrepareFunc(func(eam AccountManager) {
		wg := sync.WaitGroup{}
		var count int = 10
		ac := make(chan Accounter, count)
		for i := 0; i < count; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				a, err := eam.CreateAccount()
				assert.Nil(t, err)
				ac <- a
				wg.Done()
			}(&wg)
		}

		wg.Wait()
		close(ac)
		actMap := make(map[AccountID]bool, count)
		for a := range ac {
			actMap[a.GetAccountID()] = true
		}
		assert.Equal(t, len(actMap), count)
	})
}

// TestCreateNewBtcAddress create bitcoin address concurrently.
func TestCreateNewBtcAddress(t *testing.T) {
	PrepareFunc(func(eam AccountManager) {
		// ema.CreateAccount()
		a, err := eam.CreateAccount()
		id := a.GetAccountID()
		assert.Nil(t, err)
		wg := sync.WaitGroup{}
		var count int = 10
		addrC := make(chan string, count)
		for i := 0; i < count; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				an, err := eam.GetAccount(id)
				assert.Nil(t, err)
				addr := an.GetNewAddress(wallet.Bitcoin)
				addrC <- addr
				wg.Done()
			}(&wg)
		}
		wg.Wait()
		close(addrC)
		addrMap := make(map[string]bool, count)
		for addr := range addrC {
			addrMap[addr] = true
		}
		assert.Equal(t, count, len(addrMap))
	})
}

// func TestSetBalance(t *testing.T) {
// 	ah := newAccountHelper(t)
// 	ah.Account.SetBalance(wallet.Bitcoin, Balance(10))
// 	ah.Account.SetBalance(wallet.Skycoin, Balance(20))
// 	assert.Equal(t, ah.Account.balance[wallet.Bitcoin], Balance(10))
// 	assert.Equal(t, ah.Account.balance[wallet.Skycoin], Balance(20))
// }
//
// func TestGetBalance(t *testing.T) {
// 	ah := newAccountHelper(t)
// 	ah.Account.SetBalance(wallet.Bitcoin, Balance(10))
// 	ah.Account.SetBalance(wallet.Skycoin, Balance(20))
//
// 	bb, err := ah.Account.GetBalance(wallet.Bitcoin)
// 	assert.Nil(t, err)
// 	assert.Equal(t, bb, Balance(10))
// 	sb, err := ah.Account.GetBalance(wallet.Skycoin)
// 	assert.Nil(t, err)
// 	assert.Equal(t, sb, Balance(20))
// }

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
