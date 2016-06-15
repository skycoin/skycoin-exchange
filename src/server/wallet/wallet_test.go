package wallet

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test new wallet concurrently.
func TestNewWallet(t *testing.T) {
	num := 10000
	wltChan := make(chan Wallet, num)
	wg := sync.WaitGroup{}
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			wlt := NewWallet("")
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
	wlt := NewWallet("test")
	addrs := wlt.NewAddresses(Bitcoin, 2)
	s, err := json.MarshalIndent(addrs, "", " ")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(s))
}
