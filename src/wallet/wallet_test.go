package wallet_test

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/wallet"
	"github.com/stretchr/testify/assert"
)

// set rand seed.
var _ = func() int64 {
	t := time.Now().Unix()
	rand.Seed(t)
	return t
}()

func setup(t *testing.T) (string, func(), error) {
	wltName := fmt.Sprintf(".wallet%d", rand.Int31n(100))
	teardown := func() {}
	tmpDir := filepath.Join(os.TempDir(), wltName)
	if err := os.MkdirAll(tmpDir, 0777); err != nil {
		return "", teardown, err
	}

	// register bitcoin walleter creator.
	if err := wallet.RegisterCreator(coin.Bitcoin, wallet.NewBtcWltCreator()); err != nil {
		return "", teardown, err
	}

	teardown = func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			panic(err)
		}
	}

	return tmpDir, teardown, nil
}

func TestInitDir(t *testing.T) {
	wltName := fmt.Sprintf(".wallet%d", rand.Int31n(100))
	tmpDir := filepath.Join(os.TempDir(), wltName)
	wallet.InitDir(tmpDir)
	// check if the dir is created.
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Error("init dir failed")
		return
	}

	if wallet.GetWalletDir() != tmpDir {
		t.Error("GetWalletDir function failed")
		return
	}

	// remove the created wallet dir.
	err := os.RemoveAll(tmpDir)
	assert.Nil(t, err)
}

func TestNewWallet(t *testing.T) {
	tmpDir, teardown, err := setup(t)
	assert.Nil(t, err)
	defer teardown()

	wallet.InitDir(tmpDir)

	btcWltName := fmt.Sprintf("xde_%d", rand.Int31n(100))
	btcWlt, err := wallet.New(coin.Bitcoin, btcWltName)
	assert.Nil(t, err)

	// check the existence of wallet file.
	path := filepath.Join(wallet.GetWalletDir(), fmt.Sprintf("%s.%s", btcWlt.GetID(), wallet.Ext))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("create wallet failed")
		return
	}
}
