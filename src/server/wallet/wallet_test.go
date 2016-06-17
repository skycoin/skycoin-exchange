package wallet

import (
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
	// defer func() {
	// 	removeContents(dataDir)
	// }()
	//
	wlt, err := NewWallet("")
	if err != nil {
		t.Error(err)
		return
	}
	addrs := wlt.NewAddresses(Bitcoin, 10)
	// using the api from blockchain.info to validate the addresses.
	addrList := []string{}
	for _, addr := range addrs {
		// fmt.Println(addr.Address)
		addrList = append(addrList, addr.Address)
	}
	data, err := getDataOfUrl(fmt.Sprintf("https://blockchain.info/multiaddr?active=%s", strings.Join(addrList, "|")))
	if err != nil {
		t.Error(err)
		return
	}
	errAddr := "Invalid Bitcoin Address"
	if string(data) == errAddr {
		t.Error(errAddr)
	}
}

// func TestGetBalance(t *testing) {
//
// }

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
