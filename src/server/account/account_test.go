package account

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/stretchr/testify/assert"
)

type funcHandler func(ema AccountManager)

func DataMaintainer(f funcHandler) {
	// create
	ema := NewExchangeAccountManager()
	// do test work
	f(ema)
	// clean wallet data.
	removeContents(wallet.GetWalletDatadir())
	// clear GWallet
	wallet.Reload()
}

// func TestCreateAccountConcurrent(t *testing.T) {
// 	DataMaintainer(func(eam AccountManager) {
// 		wg := sync.WaitGroup{}
// 		var count int = 10
// 		ac := make(chan Accounter, count)
// 		for i := 0; i < count; i++ {
// 			wg.Add(1)
// 			go func(wg *sync.WaitGroup) {
// 				a, _, err := eam.CreateAccount()
// 				assert.Nil(t, err)
// 				ac <- a
// 				wg.Done()
// 			}(&wg)
// 		}
//
// 		wg.Wait()
// 		close(ac)
// 		actMap := make(map[AccountID]bool, count)
// 		for a := range ac {
// 			actMap[a.GetAccountID()] = true
// 		}
// 		assert.Equal(t, len(actMap), count)
// 	})
// }

// TestCreateNewBtcAddress create bitcoin address concurrently.
func TestCreateNewBtcAddress(t *testing.T) {
	DataMaintainer(func(eam AccountManager) {
		// ema.CreateAccount()
		a, _, err := eam.CreateAccount()
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

func TestMsgAuth(t *testing.T) {
	// DataMaintainer(func(am AccountManager) {
	// 	_, s, err := am.CreateAccount()
	// 	assert.Nil(t, err)
	// 	addr := cipher.AddressFromSecKey(s)
	//
	// 	ma := CreateMsgAuth(s, MsgAuth{Msg: []byte{"hello world"}})
	//
	// 	CheckMsgAuth(ma)
	//
	//
	// 	// cipher.a.GetAccountID()
	// })
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
