package mobile_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	api "github.com/skycoin/skycoin-exchange/src/api/mobile"
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
	api.Init(&api.Config{
		WalletDirPath: tmpDir,
		ServerAddr:    "localhost:8080",
	})

	return tmpDir, teardown, nil
}

func TestGetBalance(t *testing.T) {
	_, teardown, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()
	var testData = []struct {
		coinType string
		address  string
		expect   uint64
	}{
		{"skycoin", "cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW", 4000000},
		{"bitcoin", "1EknG7EauSW4zxFtSrCQSHe5PJenkn55s6", 994000},
	}
	for _, td := range testData {
		b, err := api.GetBalance(td.coinType, td.address)
		if err != nil {
			t.Fatal(err)
		}
		var res struct {
			Balance uint64 `json:"balance"`
		}

		if err := json.Unmarshal([]byte(b), &res); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, res.Balance, td.expect)
	}
}

func TestGetAddresses(t *testing.T) {
	_, teardown, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	var testData = []struct {
		coinType   string
		seed       string
		num        int
		expectNum  int
		expectAddr map[string]bool
	}{
		{
			"bitcoin",
			"12345",
			2,
			2,
			map[string]bool{
				"17wGH9K5sE5qq7iKiWqYd89gpz15frtYVA": true,
				"1Kn8pqcX71VooZRQrEBPr9xhAP24jSZ95n": true,
			},
		},
		{
			"skycoin",
			"12345",
			2,
			2,
			map[string]bool{
				"VYPSLGumCu1BgPUFW9yo9Y9Wm6L8v1qpZt": true,
				"Ays2XnjRKFxLR5Z5cN4VFwUoctAiftLuwB": true,
			},
		},
	}

	for _, td := range testData {
		id, err := api.NewWallet(td.coinType, td.seed)
		if err != nil {
			t.Fatal(err)
		}

		_, err = api.NewAddress(id, td.num)
		if err != nil {
			t.Fatal(err)
		}

		addrJSON, err := api.GetAddresses(id)
		if err != nil {
			t.Fatal(err)
		}
		var addrs struct {
			Addresses []string `json:"addresses"`
		}

		if err := json.Unmarshal([]byte(addrJSON), &addrs); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, len(addrs.Addresses), td.expectNum)
		for _, addr := range addrs.Addresses {
			if _, ok := td.expectAddr[addr]; !ok {
				t.Fatalf("addr: %s not expected", addr)
			}
		}
	}
}
