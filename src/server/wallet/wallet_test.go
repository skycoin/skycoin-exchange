package wallet

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test create wallets concurrently.
func TestNewWallet(t *testing.T) {
	defer func() {
		// clear all the wallet files.
		removeContents(dataDir)
	}()

	num := 10
	wltChan := make(chan Wallet, num)
	wg := sync.WaitGroup{}
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			wlt, err := NewWallet("")
			if err != nil {
				t.Error(err)
				return
			}
			wltChan <- wlt
			// fmt.Printf("%+v\n", wlt)
			wg.Done()
		}(&wg)
	}
	wg.Wait()
	close(wltChan)
	wlts := make(map[string]bool)
	for wlt := range wltChan {
		wlts[wlt.GetID()] = true
	}
	assert.Equal(t, len(wlts), num)
}

func TestNewAddress(t *testing.T) {
	wlt, err := NewWallet("test")
	assert.Equal(t, err, nil)
	addrList := []string{}
	for i := 0; i < 10; i++ {
		var ct CoinType
		if i%2 == 0 {
			ct = Bitcoin
		} else {
			ct = Skycoin
		}
		addrs, err := NewAddresses(wlt.GetID(), ct, 1)
		addrList = append(addrList, addrs...)
		assert.Equal(t, err, nil)
	}

	// assert.Equal(t, validateAddress(addrList), nil)
}

func BenchmarkNewAddress(b *testing.B) {
	wlt, err := NewWallet("test")
	if err != nil {
		fmt.Println(err)
		return
	}

	for n := 0; n < b.N; n++ {
		wlt.NewAddresses(Bitcoin, 1)
	}
}

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
