package wallet_test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
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

	teardown = func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			panic(err)
		}
	}
	wallet.InitDir(tmpDir)
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
	wltDir, teardown, err := setup(t)
	assert.Nil(t, err)
	defer teardown()

	testData := []struct {
		Type string
		Seed string
		Path string
	}{
		{"bitcoin", "sd123", filepath.Join(wltDir, "bitcoin_sd123.wlt")},
		{"bitcoin", "sd234", filepath.Join(wltDir, "bitcoin_sd234.wlt")},
		{"skycoin", "sd123", filepath.Join(wltDir, "skycoin_sd123.wlt")},
		{"skycoin", "sd234", filepath.Join(wltDir, "skycoin_sd234.wlt")},
	}

	for _, d := range testData {
		if _, err := wallet.New(d.Type, d.Seed); err != nil {
			t.Errorf("create %s wallet of seed:%s failed, err:%s", d.Type, d.Seed, err)
			return
		}

		// check the existence of wallet file.
		if _, err := os.Stat(d.Path); os.IsNotExist(err) {
			t.Error("create wallet failed")
			return
		}
	}
}

func TestNewAddresses(t *testing.T) {
	wltDir, teardown, err := setup(t)
	// wltDir, _, err := setup(t)
	assert.Nil(t, err)
	defer teardown()
	testData := []struct {
		Type    string
		Seed    string
		Num     int
		Entries []coin.AddressEntry
	}{
		{
			Type: "bitcoin",
			Seed: "sd999",
			Num:  2,
			Entries: []coin.AddressEntry{
				{
					Address: "1FLZTRDS51eiMGu1MwV75VmQPags7UjysZ",
					Public:  "0378c76e20e4f93730e67bb469bc7186681a8c85023088b64c70930e78d4aff690",
					Secret:  "L4fDKYKxMSoZ3EUfKHacykA5cM8h6EXeqQ1w2TrpeQ7f81cR5EhT",
				},
				{
					Address: "1HsUndbHFjRMSXuGyxo1kzVMsQcuhpJcwE",
					Public:  "0270d2d9b6df46e1b22effee8a3dfb42f6c3fe69b4361158b6101b451f6cced51c",
					Secret:  "Kz9vEMVPXTzTEXFrP4Pmnv79UfPRr2HWgZoQt4VAWzbUauF2MrNf",
				},
			},
		},
		{
			Type: "skycoin",
			Seed: "sd888",
			Num:  2,
			Entries: []coin.AddressEntry{
				{
					Address: "fYJPkCTqdChw3sPSGUgze9nuGMNtC5DvPY",
					Public:  "02ba572a03c8471822c308e5d041aba549b35676a0ef1c737b4517eef70c32377e",
					Secret:  "2f4aacc72a6d192e04ec540328689588caf4167d71904bdb870a4a2cee7f29c8",
				},
				{
					Address: "t6t7bJ9Ruxq9z44pYQT5AkEeAjGjgantU",
					Public:  "039f4b6a110a9c5c38da08a0bff133edf07472348a4dc4c9d63b178fe26807606e",
					Secret:  "b720d3c0f67f3c91e23805237f182e78121b90890f483133cc46f9d91232cf4c",
				},
			},
		},
	}

	for _, d := range testData {
		// new wallet
		wlt, err := wallet.New(d.Type, d.Seed)
		assert.Nil(t, err)

		for i := 0; i < d.Num; i++ {
			if _, err := wallet.NewAddresses(wlt.GetID(), 1); err != nil {
				t.Fatal(err)
			}
		}
		path := filepath.Join(wltDir, fmt.Sprintf("%s.%s", wlt.GetID(), wallet.Ext))
		cnt, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		for _, e := range d.Entries {
			if !strings.Contains(string(cnt), e.Address) {
				t.Fatalf("not contains address:%s", e.Address)
			}

			if !strings.Contains(string(cnt), e.Public) {
				t.Fatalf("not cointains pubkey:%s", e.Public)
			}
			if !strings.Contains(string(cnt), e.Secret) {
				t.Fatalf("not cointains seckey:%s", e.Secret)
			}
		}
	}
}

func TestGetAddresses(t *testing.T) {
	_, teardown, err := setup(t)
	assert.Nil(t, err)
	defer teardown()
	testData := []struct {
		Type    string
		Seed    string
		Num     int
		Entries []coin.AddressEntry
	}{
		{
			Type: "bitcoin",
			Seed: "sd999",
			Num:  2,
			Entries: []coin.AddressEntry{
				{
					Address: "1FLZTRDS51eiMGu1MwV75VmQPags7UjysZ",
					Public:  "0378c76e20e4f93730e67bb469bc7186681a8c85023088b64c70930e78d4aff690",
					Secret:  "L4fDKYKxMSoZ3EUfKHacykA5cM8h6EXeqQ1w2TrpeQ7f81cR5EhT",
				},
				{
					Address: "1HsUndbHFjRMSXuGyxo1kzVMsQcuhpJcwE",
					Public:  "0270d2d9b6df46e1b22effee8a3dfb42f6c3fe69b4361158b6101b451f6cced51c",
					Secret:  "Kz9vEMVPXTzTEXFrP4Pmnv79UfPRr2HWgZoQt4VAWzbUauF2MrNf",
				},
			},
		},
		{
			Type: "skycoin",
			Seed: "sd888",
			Num:  2,
			Entries: []coin.AddressEntry{
				{
					Address: "fYJPkCTqdChw3sPSGUgze9nuGMNtC5DvPY",
					Public:  "02ba572a03c8471822c308e5d041aba549b35676a0ef1c737b4517eef70c32377e",
					Secret:  "2f4aacc72a6d192e04ec540328689588caf4167d71904bdb870a4a2cee7f29c8",
				},
				{
					Address: "t6t7bJ9Ruxq9z44pYQT5AkEeAjGjgantU",
					Public:  "039f4b6a110a9c5c38da08a0bff133edf07472348a4dc4c9d63b178fe26807606e",
					Secret:  "b720d3c0f67f3c91e23805237f182e78121b90890f483133cc46f9d91232cf4c",
				},
			},
		},
	}

	for _, d := range testData {
		// new wallet
		wlt, err := wallet.New(d.Type, d.Seed)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := wallet.NewAddresses(wlt.GetID(), d.Num); err != nil {
			t.Fatal(err)
		}

		addrs, err := wallet.GetAddresses(wlt.GetID())
		if err != nil {
			t.Fatal(err)
		}

		for _, e := range d.Entries {
			find := func(addr string) bool {
				for _, a := range addrs {
					if a == addr {
						return true
					}
				}
				return false
			}
			if !find(e.Address) {
				t.Fatal("GetAddresses failed")
			}
		}
	}
}

func TestGetKeypair(t *testing.T) {
	_, teardown, err := setup(t)
	assert.Nil(t, err)
	defer teardown()
	testData := []struct {
		Type    string
		Seed    string
		Num     int
		Entries []coin.AddressEntry
	}{
		{
			Type: "bitcoin",
			Seed: "sd999",
			Num:  2,
			Entries: []coin.AddressEntry{
				{
					Address: "1FLZTRDS51eiMGu1MwV75VmQPags7UjysZ",
					Public:  "0378c76e20e4f93730e67bb469bc7186681a8c85023088b64c70930e78d4aff690",
					Secret:  "L4fDKYKxMSoZ3EUfKHacykA5cM8h6EXeqQ1w2TrpeQ7f81cR5EhT",
				},
				{
					Address: "1HsUndbHFjRMSXuGyxo1kzVMsQcuhpJcwE",
					Public:  "0270d2d9b6df46e1b22effee8a3dfb42f6c3fe69b4361158b6101b451f6cced51c",
					Secret:  "Kz9vEMVPXTzTEXFrP4Pmnv79UfPRr2HWgZoQt4VAWzbUauF2MrNf",
				},
			},
		},
		{
			Type: "skycoin",
			Seed: "sd888",
			Num:  2,
			Entries: []coin.AddressEntry{
				{
					Address: "fYJPkCTqdChw3sPSGUgze9nuGMNtC5DvPY",
					Public:  "02ba572a03c8471822c308e5d041aba549b35676a0ef1c737b4517eef70c32377e",
					Secret:  "2f4aacc72a6d192e04ec540328689588caf4167d71904bdb870a4a2cee7f29c8",
				},
				{
					Address: "t6t7bJ9Ruxq9z44pYQT5AkEeAjGjgantU",
					Public:  "039f4b6a110a9c5c38da08a0bff133edf07472348a4dc4c9d63b178fe26807606e",
					Secret:  "b720d3c0f67f3c91e23805237f182e78121b90890f483133cc46f9d91232cf4c",
				},
			},
		},
	}

	for _, d := range testData {
		// new wallet
		wlt, err := wallet.New(d.Type, d.Seed)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := wallet.NewAddresses(wlt.GetID(), d.Num); err != nil {
			t.Fatal(err)
		}

		for _, e := range d.Entries {
			p, s, err := wallet.GetKeypair(wlt.GetID(), e.Address)
			if err != nil {
				t.Fatal(err)
			}
			if p != e.Public || s != e.Secret {
				t.Fatal("get key pair failed")
			}
		}
	}
}

func TestRemove(t *testing.T) {
	wltDir, teardown, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	// create wallet
	testData := []struct {
		Type string
		Seed string
		ID   string
	}{
		{"bitcoin", "sd777", "bitcoin_sd777"},
		{"skycoin", "sd777", "skycoin_sd777"},
	}

	for _, d := range testData {
		wlt, err := wallet.New(d.Type, d.Seed)
		if err != nil {
			t.Fatal(err)
		}

		// remove this wlt.
		if err := wallet.Remove(wlt.GetID()); err != nil {
			t.Fatal(err)
		}

		// check if the wlt file is already removed.
		path := filepath.Join(wltDir, fmt.Sprintf("%s.%s", d.ID, wallet.Ext))
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatal("remove wallet failed")
		}
	}
}

func TestIsExist(t *testing.T) {
	_, teardown, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}

	defer teardown()

	testData := []struct {
		Type string
		Seed string
	}{
		{"bitcoin", "sd666"},
		{"bitcoin", "sd667"},
		{"skycoin", "sd666"},
		{"skycoin", "sd667"},
	}

	for _, d := range testData {
		id := wallet.MakeWltID(d.Type, d.Seed)
		if wallet.IsExist(id) {
			t.Fatalf("wallet:%s should not exist", id)
		}

		_, err := wallet.New(d.Type, d.Seed)
		if err != nil {
			t.Fatalf("creat wallet :%s failed", id)
		}

		if !wallet.IsExist(id) {
			t.Fatalf("wallet:%s should exist", id)
		}
	}
}
