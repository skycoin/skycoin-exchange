package mobile_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
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

func setup() (string, func(), error) {
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
		// ServerAddr:    "121.41.103.148:8080",
		ServerAddr: "localhost:8080",
	})

	return tmpDir, teardown, nil
}

func TestGetBalance(t *testing.T) {
	_, teardown, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()
	var testData = []struct {
		coinType string
		address  string
		expect   uint64
	}{
		{"skycoin", "cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW", 6000000},
		{"bitcoin", "1EknG7EauSW4zxFtSrCQSHe5PJenkn55s6", 936000},
		{"mzcoin", "2BMHv3PEyat9K9snsnDyRv7UBuRuycMPyWH", 998000000},
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
	_, teardown, err := setup()
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
		{
			"mzcoin",
			"2345",
			2,
			2,
			map[string]bool{
				"5v6aHMp7dxwFzZmnF1c2XjcJSKnaMixXmr":  true,
				"2CiP9MXLy6KpUyUzU3Y8nTRTuffYwreqgqa": true,
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
	_, teardown, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	id, err := api.NewWallet("bitcoin", "asdfasdf")
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
	_, teardown, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	id, err := api.NewWallet("skycoin", "qwerqer")
	if err != nil {
		t.Fatal(err)
	}

	_, err = api.NewAddress(id, 3)
	if err != nil {
		t.Fatal(err)
	}

	txid, err := api.SendSky(id, "UsS43vk2yRqjXvgbwq12Dkjr8cHVTBxYoj", "1000000")
	// if err != nil {
	// 	t.Fatal(err)
	// 	return
	// }
	if !strings.Contains(err.Error(), "insufficient balance") {
		t.Fatal(err)
	}
	fmt.Println(txid)
}

func TestSendMzc(t *testing.T) {
	_, teardown, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	id, err := api.NewWallet("mzcoin", "99999")
	if err != nil {
		t.Fatal(err)
	}

	_, err = api.NewAddress(id, 3)
	if err != nil {
		t.Fatal(err)
	}

	txid, err := api.SendMzc(id, "ntXEnuc6JoDie9eV4jEFFvALiRMpadhyGS", "1000000")
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
	_, teardown, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

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
		{
			"mzcoin normal",
			args{
				"mzcoin",
				"ae27da319436e0397dbfc1d596b7dccb71238dec77120137f9347426afd668a2",
			},
			"16d2e0a200ecdb7a363debfe5883b1e8f0c902aabb831f7e5e8908ccbd037388",
			nil,
		},
		{
			"mzcoin invalid txid len",
			args{
				"mzcoin",
				"b1481d",
			},
			"",
			errors.New("invalid transaction id"),
		},
		{
			"mzcoin invalid txid",
			args{
				"mzcoin",
				"b1481d614ffcc27408fe2131198d9d2821c78601a0aa23d8e9965b2a5196edc1",
			},
			"",
			errors.New("not found\n"),
		},
	}
	for _, tt := range tests {
		got, err := api.GetTransactionByID(tt.args.coinType, tt.args.txid)
		if !reflect.DeepEqual(err, tt.wantErr) {
			t.Errorf("%q. GetTransactionByID() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !strings.Contains(got, tt.contain) {
			t.Errorf("%q. GetTransactionByID() = %v, want %v", tt.name, got, tt.contain)
		}
	}
}

func TestGetOutputByID(t *testing.T) {
	_, teardown, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	type args struct {
		hash string
	}
	tests := []struct {
		name     string
		coinType string
		hash     string
		contain  string
		wantErr  error
	}{
		// TODO: Add test cases.
		{
			"skycin normal",
			"skycoin",
			"a57c038591f862b8fada57e496ef948183b153348d7932921f865a8541a477c5",
			"cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW",
			nil,
		},
		{
			"invalid hash len",
			"skycoin",
			"a57c038",
			"",
			errors.New("invalid output hash, encoding/hex: odd length hex string"),
		},
		{
			"invalid hash",
			"skycoin",
			"a57c038591f862b8fada57e496ef948183b153348d7932921f865a8541a477c9",
			"",
			errors.New("not found\n"),
		},
		{
			"mzcoin normal",
			"mzcoin",
			"0c22889e2d76512aea33063eae9e1a05891e4e6f55514c72ceabee0f813cca0b",
			"2BMHv3PEyat9K9snsnDyRv7UBuRuycMPyWH",
			nil,
		},
	}
	for _, tt := range tests {
		got, err := api.GetOutputByID(tt.coinType, tt.hash)
		if !reflect.DeepEqual(err, tt.wantErr) {
			t.Errorf("%q. GetOutputByHash() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !strings.Contains(got, tt.contain) {
			t.Errorf("%q. GetOutputByHash() = %v, want contains %v", tt.name, got, tt.contain)
		}
	}
}

// func BenchmarkGetBalance(b *testing.B) {
// 	_, teardown, err := setup()
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer teardown()
// 	for i := 0; i < b.N; i++ {
// 		var testData = []struct {
// 			coinType string
// 			address  string
// 			expect   uint64
// 		}{
// 			{"skycoin", "cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW", 4000000},
// 			{"bitcoin", "1EknG7EauSW4zxFtSrCQSHe5PJenkn55s6", 938000},
// 		}

// 		var err error
// 		for _, td := range testData {
// 			_, err = api.GetBalance(td.coinType, td.address)
// 			if err != nil {
// 				panic(err)
// 			}
// 		}
// 	}
// }

// func BenchmarkGetOutByID(b *testing.B) {
// 	_, teardown, err := setup()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer teardown()
// 	for i := 0; i < b.N; i++ {
// 		api.GetSkyOutputByID("a57c038591f862b8fada57e496ef948183b153348d7932921f865a8541a477c5")
// 	}
// }

// func BenchmarkGetTx(b *testing.B) {
// 	_, teardown, err := setup()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer teardown()
// 	for i := 0; i < b.N; i++ {
// 		go func() {
// 			api.GetTransactionByID("bitcoin", "69be3a3b98541e609f5a4935f94c92012d2b3e3437e9508770ba2257f532142f")
// 			api.GetTransactionByID("skycoin", "b1481d614ffcc27408fe2131198d9d2821c78601a0aa23d8e9965b2a5196edc0")

// 		}()
// 	}
// }

func TestValidateAddress(t *testing.T) {
	_, teardown, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()
	type args struct {
		coinType string
		addr     string
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr error
	}{
		// TODO: Add test cases.
		{
			"normal bitcoin",
			args{
				"bitcoin",
				"14NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz",
			},
			true,
			nil,
		},
		{
			"invalid bitcoin address length",
			args{
				"bitcoin",
				"14NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz980",
			},
			false,
			errors.New("Invalid address length"),
		},
		{
			"invalid bitcoin address version",
			args{
				"bitcoin",
				"24NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz",
			},
			false,
			errors.New("Invalid version"),
		},
		{
			"normal skycoin address",
			args{
				"skycoin",
				"cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW",
			},
			true,
			nil,
		},
		{
			"invalid skycoin address length",
			args{
				"skycoin",
				"cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW100",
			},
			false,
			errors.New("Invalid address length"),
		},
		{
			"invalid skycoin address version",
			args{
				"skycoin",
				"BBbbbbbvv12dovBmjQKTtfE4rbjMmf3fzW",
			},
			false,
			errors.New("Invalid version"),
		},
	}
	for _, tt := range tests {
		got, err := api.ValidateAddress(tt.args.coinType, tt.args.addr)
		if !reflect.DeepEqual(err, tt.wantErr) {
			t.Errorf("%q. ValidateAddress() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. ValidateAddress() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
