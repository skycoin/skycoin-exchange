package account_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/server/account"
)

func TestInitDir(t *testing.T) {
	tmpDir := os.TempDir()
	dir := tmpDir + "/.skycoin-exchange/account"
	account.InitDir(dir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("InitDir faild")
		return
	}

	// clear the dir
	os.RemoveAll(filepath.Dir(dir))
}

func TestGetID(t *testing.T) {
	a := account.ExchangeAccount{
		ID: "1234",
	}

	if a.GetID() != "1234" {
		t.Error("get account id failed")
		return
	}
}

func TestGetBalance(t *testing.T) {
	a := account.ExchangeAccount{
		Balance: map[coin.Type]uint64{
			coin.Bitcoin: 90000,
			coin.Skycoin: 450000,
		},
	}

	if a.GetBalance(coin.Bitcoin) != 90000 {
		t.Error("get bitcoin balance failed")
		return
	}

	if a.GetBalance(coin.Skycoin) != 450000 {
		t.Error("get skycoin balance failed")
		return
	}
}

func TestIncreaseBalance(t *testing.T) {
	var btcInit uint64 = 90000
	var skyInit uint64 = 450000
	testData := map[coin.Type][]struct {
		V      uint64
		Expect uint64
	}{
		coin.Bitcoin: {
			{10000, 100000},
			{20000, 110000},
			{1000, 91000},
			{100, 90100},
		},
		coin.Skycoin: {
			{10000, 460000},
			{30000, 480000},
			{50000, 500000},
		},
	}

	for cp, tds := range testData {
		for _, d := range tds {
			a := account.ExchangeAccount{
				Balance: map[coin.Type]uint64{
					coin.Bitcoin: btcInit,
					coin.Skycoin: skyInit,
				},
			}
			if err := a.IncreaseBalance(cp, d.V); err != nil {
				t.Error(err)
				return
			}

			if a.GetBalance(cp) != d.Expect {
				t.Errorf("decrease %s balance failed, v:%d, expect:%d", cp, a.GetBalance(cp), d.Expect)
				return
			}
		}
	}
}

func TestDecreaseBalance(t *testing.T) {
	var btcInit uint64 = 90000
	var skyInit uint64 = 450000
	testData := map[coin.Type][]struct {
		V      uint64
		Expect uint64
	}{
		coin.Bitcoin: {
			{10000, 80000},
			{20000, 70000},
			{1000, 89000},
			{100, 89900},
		},
		coin.Skycoin: {
			{10000, 440000},
			{30000, 420000},
			{50000, 400000},
		},
	}

	for cp, tds := range testData {
		for _, d := range tds {
			a := account.ExchangeAccount{
				Balance: map[coin.Type]uint64{
					coin.Bitcoin: btcInit,
					coin.Skycoin: skyInit,
				},
			}
			if err := a.DecreaseBalance(cp, d.V); err != nil {
				t.Error(err)
				return
			}
			b := a.GetBalance(cp)
			if b != d.Expect {
				t.Errorf("decrease %s balance failed, v:%d, expect:%d", cp, b, d.Expect)
				return
			}
		}
	}

}
