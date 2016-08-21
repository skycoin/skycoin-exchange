package wallet_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/skycoin/skycoin-exchange/src/rpclient/wallet"
	"github.com/skycoin/skycoin-exchange/src/server/coin"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (func(), error) {
	teardown := func() {}
	tmpDir := filepath.Join(os.TempDir(), "./rpclient/wallet")
	if err := os.MkdirAll(tmpDir, 0777); err != nil {
		return teardown, err
	}

	teardown = func() {
		if err := os.RemoveAll(filepath.Dir(tmpDir)); err != nil {
			panic(err)
		}
	}

	return teardown, nil
}

func TestInitDir(t *testing.T) {
	tmpDir := os.TempDir() + "rpclient/wallet"

	wallet.InitDir(tmpDir)
	// check the wallet dir
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Error("rpclient wallet init dir failed")
		return
	}

	// clear the wallet dir
	os.RemoveAll(filepath.Dir(tmpDir))
}

func TestNewWallet(t *testing.T) {
	_, err := setup(t)
	assert.Nil(t, err)
	// defer teardown()

	_, err = wallet.New(coin.Bitcoin, "testseed")
	assert.Nil(t, err)
	// check the existence of wallet file.
	path := filepath.Join(wallet.GetWalletDir(), fmt.Sprintf("btc_testseed.wlt"))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("create wallet failed")
		return
	}
}
