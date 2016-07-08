package wallet

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWallet(t *testing.T) {
	wlt, err := New("server.wlt", Deterministic, "test")
	assert.Nil(t, err)
	path := filepath.Join(WltDir, "server.wlt")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fail()
	}

	// create address
	_, err = wlt.NewAddresses(Bitcoin, 1)
	assert.Nil(t, err)
}

func TestLoadWallet(t *testing.T) {
	wlt, err := New("server.wlt", Deterministic, "test")
	assert.Nil(t, err)
	path := filepath.Join(WltDir, "server.wlt")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fail()
	}

	// create address
	e, err := wlt.NewAddresses(Bitcoin, 1)
	assert.Nil(t, err)

	wlt1, err := Load(path)
	assert.Nil(t, err)

	addrs := wlt1.GetAddressEntries(Bitcoin)
	assert.Nil(t, err)
	assert.Equal(t, e[0].Address, addrs[0].Address)
	assert.Equal(t, e[0].Public, addrs[0].Public)
	assert.Equal(t, e[0].Secret, addrs[0].Secret)
	assert.Equal(t, "server.wlt", wlt1.GetID())
}

// test create wallets concurrently.
// func TestNewWallet(t *testing.T) {
// 	defer func() {
// 		// clear all the wallet files.
// 		removeContents(dataDir)
// 	}()
//
// 	num := 10
// 	wltChan := make(chan Wallet, num)
// 	wg := sync.WaitGroup{}
// 	for i := 0; i < num; i++ {
// 		wg.Add(1)
// 		go func(wg *sync.WaitGroup) {
// 			wlt, err := New("")
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			wltChan <- wlt
// 			// fmt.Printf("%+v\n", wlt)
// 			wg.Done()
// 		}(&wg)
// 	}
// 	wg.Wait()
// 	close(wltChan)
// 	wlts := make(map[string]bool)
// 	for wlt := range wltChan {
// 		wlts[wlt.GetID()] = true
// 	}
// 	assert.Equal(t, len(wlts), num)
// }

// func TestNewAddress(t *testing.T) {
// 	wlt, err := New("test")
// 	assert.Equal(t, err, nil)
// 	addrList := []string{}
// 	for i := 0; i < 10; i++ {
// 		var ct CoinType
// 		if i%2 == 0 {
// 			ct = Bitcoin
// 		} else {
// 			ct = Skycoin
// 		}
// 		addrs, err := NewAddresses(wlt.GetID(), ct, 1)
// 		addrList = append(addrList, addrs...)
// 		assert.Equal(t, err, nil)
// 	}
//
// 	// assert.Equal(t, validateAddress(addrList), nil)
// }
//
// func BenchmarkNewAddress(b *testing.B) {
// 	wlt, err := New("test")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
//
// 	for n := 0; n < b.N; n++ {
// 		wlt.NewAddresses(Bitcoin, 1)
// 	}
// }

// func TestStringBitcoinToAddress(t *testing.T) {
// 	btcAddr := "156ua4hbst6TQ4x47yZPqLMRWQSkE8puvs"
// 	addr := cipher.BitcoinMustDecodeBase58Address(btcAddr)
// 	assert.Equal(t, validateAddress([]string{addr.BitcoinString()}), nil)
// }

func validateAddress(addrs []string) error {
	data, err := getDataOfUrl(fmt.Sprintf("https://blockchain.info/multiaddr?active=%s", strings.Join(addrs, "|")))
	if err != nil {
		return err
	}
	errAddr := "Invalid Bitcoin Address"
	if string(data) == errAddr {
		return errors.New(errAddr)
	}
	return nil
}

func getDataOfUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	resp.Body.Close()
	return data, nil
}

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
