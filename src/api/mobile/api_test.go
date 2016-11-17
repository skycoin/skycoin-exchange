package mobile_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
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
		{"bitcoin", "1EknG7EauSW4zxFtSrCQSHe5PJenkn55s6", 938000},
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

func TestSendBtc(t *testing.T) {
	_, teardown, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	id, err := api.NewWallet("bitcoin", "adfasda")
	if err != nil {
		t.Fatal(err)
	}

	_, err = api.NewAddress(id, 3)
	if err != nil {
		t.Fatal(err)
	}

	txid, err := api.SendBtc(id, "14NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz", "1000", "1000")
	// if err != nil {
	// 	t.Fatal(err)
	// 	return
	// }
	if !strings.Contains(err.Error(), "insufficient balance") {
		t.Fatal(err)
	}
	fmt.Println(txid)
}

func TestSendSky(t *testing.T) {
	_, teardown, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	id, err := api.NewWallet("skycoin", "adfasda")
	if err != nil {
		t.Fatal(err)
	}

	_, err = api.NewAddress(id, 3)
	if err != nil {
		t.Fatal(err)
	}

	txid, err := api.SendSky(id, "fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B", "1000000")
	// if err != nil {
	// 	t.Fatal(err)
	// 	return
	// }
	if !strings.Contains(err.Error(), "insufficient balance") {
		t.Fatal(err)
	}
	fmt.Println(txid)
}

func TestGetTransactionByID(t *testing.T) {
	type args struct {
		coinType string
		txid     string
	}
	tests := []struct {
		name    string
		args    args
		contain string
		wantErr error
	}{
		// TODO: Add test cases.
		{
			"bitcoin normal",
			args{
				"bitcoin",
				"69be3a3b98541e609f5a4935f94c92012d2b3e3437e9508770ba2257f532142f",
			},
			"0000000000000000021f3adb9ce12e3c70a42cfb6b7095805bee7bdefb392725",
			nil,
		},
		{
			"bitcoin invalid txid len",
			args{
				"bitcoin",
				"69be3a3b98541e6",
			},
			"",
			errors.New("invalid transaction id"),
		},
		{
			"bitcoin invalid txid",
			args{
				"bitcoin",
				"69be3a3b98541e609f5a4935f94c92012d2b3e3437e9508770ba2257f532142d",
			},
			"",
			errors.New("not found"),
		},
		{
			"skycoin normal",
			args{
				"skycoin",
				"b1481d614ffcc27408fe2131198d9d2821c78601a0aa23d8e9965b2a5196edc0",
			},
			"7583587d02bedbeb3c15dde9e13baac36b0eb2b7ba7b2063c323a226d0784619",
			nil,
		},
		{
			"skycoin invalid txid len",
			args{
				"skycoin",
				"b1481d",
			},
			"",
			errors.New("invalid transaction id"),
		},
		{
			"skycoin invalid txid",
			args{
				"skycoin",
				"b1481d614ffcc27408fe2131198d9d2821c78601a0aa23d8e9965b2a5196edc1",
			},
			"",
			errors.New("not found\n"),
		},
		{
			"invalid coin type",
			args{
				"unknow",
				"",
			},
			"",
			errors.New("unknow is not supported"),
		},
	}
	for _, tt := range tests {
		got, err := api.GetTransactionByID(tt.args.coinType, tt.args.txid)
		if !assert.Equal(t, err, tt.wantErr) {
			t.Errorf("%q. GetTransactionByID() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !strings.Contains(got, tt.contain) {
			t.Errorf("%q. GetTransactionByID() = %v, want %v", tt.name, got, tt.contain)
		}
	}
}
