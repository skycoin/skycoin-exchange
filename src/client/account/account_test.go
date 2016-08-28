package account_test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/skycoin/skycoin-exchange/src/client/account"
	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/stretchr/testify/assert"
)

// set rand seed.
var _ = func() int64 {
	t := time.Now().Unix()
	rand.Seed(t)
	return t
}()

func setup(t *testing.T) (string, func(), error) {
	actName := fmt.Sprintf(".account%d", rand.Int31n(100))
	teardown := func() {}
	tmpDir := filepath.Join(os.TempDir(), actName)
	if err := os.MkdirAll(tmpDir, 0777); err != nil {
		return "", teardown, err
	}

	teardown = func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			panic(err)
		}
	}
	account.InitDir(tmpDir)
	return tmpDir, teardown, nil
}

func TestInitDir(t *testing.T) {
	// dir, teardown, err := setup(t)
	dir, _, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(dir)
	// defer teardown()
	// check the exitence of dir.
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("account init dir failed")
	}

	// store some account.
	p, s := cipher.GenerateKeyPair()
	acnt := account.Account{
		Pubkey: p.Hex(),
		Seckey: s.Hex(),
	}

	v := fmt.Sprintf(`{"active_account":"%s","accounts":[{"pubkey":"%s", "seckey":"%s"}]}`, p.Hex(), p.Hex(), s.Hex())

	if err := ioutil.WriteFile(filepath.Join(dir, "data.act"), []byte(v), 0777); err != nil {
		t.Fatal(err)
	}

	// init dir again.
	account.InitDir(dir)

	// check if the account is loaded.
	a, err := account.Get(p.Hex())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, acnt.Pubkey, a.Pubkey)
	assert.Equal(t, acnt.Seckey, a.Seckey)

	activeAccount := account.GetActive()
	assert.Equal(t, activeAccount.Pubkey, p.Hex())
}

func TestNewAndSet(t *testing.T) {
	dir, teardown, err := setup(t)
	// dir, _, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	a := account.New()
	a.WltIDs[coin.Bitcoin] = "bitcoin_sd110"
	a.WltIDs[coin.Skycoin] = "skycoin_sd110"
	account.Set(a)

	// get account
	d, err := ioutil.ReadFile(filepath.Join(dir, account.FileName()))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(d), a.Pubkey) {
		t.Fatal("new account failed")
	}

	if !strings.Contains(string(d), a.Seckey) {
		t.Fatal("new account failed")
	}

	if !strings.Contains(string(d), "bitcoin_sd110") {
		t.Fatal("new account failed")
	}

	if !strings.Contains(string(d), "skycoin_sd110") {
		t.Fatal("new account failed")
	}
}

func TestGetAccount(t *testing.T) {
	_, teardown, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}

	defer teardown()

	// new account
	a := account.New()
	account.Set(a)

	act, err := account.Get(a.Pubkey)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, act, a)

	// add wallet id
	act.WltIDs[coin.Bitcoin] = "bitcoin_sd19"
	act.WltIDs[coin.Skycoin] = "skycoin_sd19"
	account.Set(act)

	newA, err := account.Get(act.Pubkey)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, newA, act)
}
