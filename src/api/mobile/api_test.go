package mobile

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

	"github.com/skycoin/skycoin-exchange/src/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

//go:generate goautomock -template=testify Coiner

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

	Init(&Config{
		WalletDirPath: tmpDir,
	})

	return tmpDir, teardown, nil
}

func TestGetBalance(t *testing.T) {
	skyM := NewCoinerMock()
	skyM.On("Name").Return("skycoin")
	skyM.On("ValidateAddr", "cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW").Return(nil)
	skyM.On("GetBalance", []string{"cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW"}).Return(uint64(6e6), nil)

	mzM := NewCoinerMock()
	mzM.On("Name").Return("mzcoin")
	mzM.On("ValidateAddr", "2BMHv3PEyat9K9snsnDyRv7UBuRuycMPyWH").Return(nil)
	mzM.On("GetBalance", []string{"2BMHv3PEyat9K9snsnDyRv7UBuRuycMPyWH"}).Return(uint64(998e6), nil)

	btcM := NewCoinerMock()
	btcM.On("Name").Return("bitcoin")
	btcM.On("ValidateAddr", "1EknG7EauSW4zxFtSrCQSHe5PJenkn55s6").Return(nil)
	btcM.On("GetBalance", []string{"1EknG7EauSW4zxFtSrCQSHe5PJenkn55s6"}).Return(uint64(936000), nil)

	shellM := NewCoinerMock()
	shellM.On("Name").Return("shellcoin")
	shellM.On("ValidateAddr", "1EknG7EauSW4zxFtSrCQSHe5PJenkn55s6").Return(nil)
	shellM.On("GetBalance", []string{"1EknG7EauSW4zxFtSrCQSHe5PJenkn55s6"}).Return(uint64(10e6), nil)

	initConfig(&Config{}, skyM, mzM, btcM, shellM)

	var testData = []struct {
		coinType string
		address  string
		expect   uint64
	}{
		{"skycoin", "cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW", 6000000},
		{"bitcoin", "1EknG7EauSW4zxFtSrCQSHe5PJenkn55s6", 936000},
		{"mzcoin", "2BMHv3PEyat9K9snsnDyRv7UBuRuycMPyWH", 998000000},
		{"shellcoin", "1EknG7EauSW4zxFtSrCQSHe5PJenkn55s6", 10e6},
	}
	for _, td := range testData {
		b, err := GetBalance(td.coinType, td.address)
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

// func TestGetWalletBalance(t *testing.T) {
// 	_, teardown, err := setup()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer teardown()

// 	coinType := "skycoin"

// 	id, err := api.NewWallet(coinType, "test123")
// 	assert.Nil(t, err)
// 	api.NewAddress()

// 	api.GetWalletBalance("skycoin")
// }

func TestNewWallet(t *testing.T) {
	testCases := []struct {
		name        string
		coinType    string
		seed        string
		expectWltID string
		expectErr   error
	}{
		{
			"create skycoin wallet",
			"skycoin",
			"abc",
			"skycoin_abc",
			nil,
		},
		{
			"create mzcoin wallet",
			"mzcoin",
			"abcd",
			"mzcoin_abcd",
			nil,
		},
		{
			"create shellcoin wallet",
			"shellcoin",
			"abcde",
			"shellcoin_abcde",
			nil,
		},
		{
			"create suncoin wallet",
			"suncoin",
			"abcde",
			"suncoin_abcde",
			nil,
		},
		{
			"create aynrandcoin wallet",
			"aynrandcoin",
			"abcde",
			"aynrandcoin_abcde",
			nil,
		},
		{
			"create metalicoin wallet",
			"metalicoin",
			"abcde",
			"metalicoin_abcde",
			nil,
		},
		{
			"create bitcoin wallet",
			"bitcoin",
			"abcde",
			"bitcoin_abcde",
			nil,
		},
		{
			"create unknow wallet",
			"unknow",
			"abcde",
			"",
			errors.New("unknow wallet not regestered"),
		},
	}

	_, teardown, err := setup()
	if err != nil {
		t.Fatal(err)
		return
	}
	defer teardown()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := NewWallet(tc.coinType, tc.seed)
			require.Equal(t, tc.expectErr, err)
			require.Equal(t, tc.expectWltID, id)
		})
	}
}

func TestGetAddresses(t *testing.T) {
	_, teardown, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	var testData = []struct {
		name       string
		coinType   string
		seed       string
		num        int
		expectNum  int
		expectAddr map[string]bool
	}{
		{
			"get from bitcoin wallet",
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
			"get from skycoin wallet",
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
			"get from mzcoin wallet",
			"mzcoin",
			"2345",
			2,
			2,
			map[string]bool{
				"5v6aHMp7dxwFzZmnF1c2XjcJSKnaMixXmr":  true,
				"2CiP9MXLy6KpUyUzU3Y8nTRTuffYwreqgqa": true,
			},
		},
		{
			"get from shellcoin wallet",
			"shellcoin",
			"2345",
			2,
			2,
			map[string]bool{
				"5v6aHMp7dxwFzZmnF1c2XjcJSKnaMixXmr":  true,
				"2CiP9MXLy6KpUyUzU3Y8nTRTuffYwreqgqa": true,
			},
		},
		{
			"get from suncoin wallet",
			"suncoin",
			"2345",
			2,
			2,
			map[string]bool{
				"5v6aHMp7dxwFzZmnF1c2XjcJSKnaMixXmr":  true,
				"2CiP9MXLy6KpUyUzU3Y8nTRTuffYwreqgqa": true,
			},
		},
		{
			"get from aynrandcoin wallet",
			"aynrandcoin",
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
		t.Run(td.name, func(t *testing.T) {
			id, err := NewWallet(td.coinType, td.seed)
			if err != nil {
				t.Fatal(err)
			}

			_, err = NewAddress(id, td.num)
			if err != nil {
				t.Fatal(err)
			}

			addrJSON, err := GetAddresses(id)
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
		})
	}
}

func TestSend(t *testing.T) {
	txid := "32444c08568cf03f4be5bb1110124d6a00bb94bc5338abddc9fb2497f3825a91"
	btc := NewCoinerMock()
	btc.On("Name").Return("bitcoin")
	btc.On("Send", "bitcoin_abc", "14NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz", "10000", mock.AnythingOfType("[]mobile.Option")).
		Return(fmt.Sprintf(`{"txid":"%s"}`, txid), nil)
	btc.On("Send", "bitcoin_abc", "14NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz", "100ss", mock.AnythingOfType("[]mobile.Option")).
		Return("",
			fmt.Errorf(`parse amount string to uint64 failed: strconv.ParseUint: parsing "100ss": invalid syntax"`))

	sky := NewCoinerMock()
	sky.On("Name").Return("skycoin")

	sky.On("Send", "skycoin_abc", "14NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz", "1e6", []Option(nil)).
		Return(fmt.Sprintf(`{"txid":"%s"}`, txid), nil)
	sky.On("Send", "skycoin_abc", "14NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz", "100ss", []Option(nil)).
		Return("", errors.New(`parse amount string to uint64 failed: strconv.ParseUint: parsing "100ss": invalid syntax"`))

	initConfig(&Config{}, btc, sky)

	type args struct {
		walletID string
		toAddr   string
		amount   string
		opt      *SendOption
	}
	tests := []struct {
		name     string
		coinType string
		args     args
		want     string
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			"bitcoin send normal",
			"bitcoin",
			args{
				"bitcoin_abc",
				"14NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz",
				"10000",
				&SendOption{Fee: "1000"},
			},
			fmt.Sprintf(`{"txid":"%s"}`, txid),
			false,
		},
		{
			"bitcoin invalid amount",
			"bitcoin",
			args{
				"bitcoin_abc",
				"14NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz",
				"100ss",
				&SendOption{Fee: "1000"},
			},
			"",
			true,
		},
		{
			"skycoin normal",
			"skycoin",
			args{
				"skycoin_abc",
				"14NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz",
				"1e6",
				nil,
			},
			fmt.Sprintf(`{"txid":"%s"}`, txid),
			false,
		},
		{
			"invalid amount",
			"skycoin",
			args{
				"skycoin_abc",
				"14NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz",
				"100ss",
				nil,
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		got, err := Send(tt.coinType, tt.args.walletID, tt.args.toAddr, tt.args.amount, tt.args.opt)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. SendBtc() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. SendBtc() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

// func TestSendSky(t *testing.T) {
// 	txid := "32444c08568cf03f4be5bb1110124d6a00bb94bc5338abddc9fb2497f3825a91"
// 	m := NewCoinerMock()
// 	m.On("Name").Return("skycoin")

// 	m.On("Send", "skycoin_abc", "14NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz", "1e6", []Option(nil)).
// 		Return(fmt.Sprintf(`{"txid":"%s"}`, txid), nil)
// 	m.On("Send", "skycoin_abc", "14NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz", "100ss", []Option(nil)).
// 		Return("", errors.New(`parse amount string to uint64 failed: strconv.ParseUint: parsing "100ss": invalid syntax"`))

// 	initConfig(&Config{}, m)

// 	type args struct {
// 		walletID string
// 		toAddr   string
// 		amount   string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    string
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			"normal",
// 			args{
// 				"skycoin_abc",
// 				"14NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz",
// 				"1e6",
// 			},
// 			fmt.Sprintf(`{"txid":"%s"}`, txid),
// 			false,
// 		},
// 		{
// 			"invalid amount",
// 			args{
// 				"skycoin_abc",
// 				"14NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz",
// 				"100ss",
// 			},
// 			"",
// 			true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		got, err := SendSky(tt.args.walletID, tt.args.toAddr, tt.args.amount)
// 		if (err != nil) != tt.wantErr {
// 			t.Errorf("%q. SendBtc() error = %v, wantErr %v", tt.name, err, tt.wantErr)
// 			continue
// 		}
// 		if got != tt.want {
// 			t.Errorf("%q. SendBtc() = %v, want %v", tt.name, got, tt.want)
// 		}
// 	}
// }

// func TestSendMzc(t *testing.T) {
// 	txid := "32444c08568cf03f4be5bb1110124d6a00bb94bc5338abddc9fb2497f3825a91"
// 	m := NewCoinerMock()
// 	m.On("Name").Return("mzcoin")

// 	m.On("Send", "mzcoin_abc", "14NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz", "1e6", []Option(nil)).
// 		Return(fmt.Sprintf(`{"txid":"%s"}`, txid), nil)
// 	m.On("Send", "mzcoin_abc", "14NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz", "100ss", []Option(nil)).
// 		Return("", errors.New(`parse amount string to uint64 failed: strconv.ParseUint: parsing "100ss": invalid syntax"`))

// 	initConfig(&Config{}, m)

// 	type args struct {
// 		walletID string
// 		toAddr   string
// 		amount   string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    string
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			"normal",
// 			args{
// 				"mzcoin_abc",
// 				"14NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz",
// 				"1e6",
// 			},
// 			fmt.Sprintf(`{"txid":"%s"}`, txid),
// 			false,
// 		},
// 		{
// 			"invalid amount",
// 			args{
// 				"mzcoin_abc",
// 				"14NAt8DhxMYKUwP5ZyH1yu7m1psYsn9Wqz",
// 				"100ss",
// 			},
// 			"",
// 			true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		got, err := SendMzc(tt.args.walletID, tt.args.toAddr, tt.args.amount)
// 		if (err != nil) != tt.wantErr {
// 			t.Errorf("%q. SendBtc() error = %v, wantErr %v", tt.name, err, tt.wantErr)
// 			continue
// 		}
// 		if got != tt.want {
// 			t.Errorf("%q. SendBtc() = %v, want %v", tt.name, got, tt.want)
// 		}
// 	}
// }

var skyTxStr = `{
    "status": {
        "confirmed": true,
        "unconfirmed": false,
        "height": 9,
        "unknown": false
    },
    "txn": {
        "length": 183,
        "type": 0,
        "txid": "367fc68cd78adc5ed5361f9cd982289f4815da6db5a9f0bdb6c59cf463018b00",
        "inner_hash": "b8c519e34942ffaf1aa5a36e7df5b5cc6387cf5f055aced8d039c4db5216288e",
        "timestamp": 1483714555,
        "sigs": [
            "2a3874d06e0627eb7b99c725957ee697f9b862562e7b347aa9afc680ca7801cc41c4d1db5c8b5696341e19ade406d9f580c7fba126d284707b1e327f1bfcd07901"
        ],
        "inputs": [
            "aced4e58f22774056d2419d41f52c71920211af72c596bb5f8fd222baa41b586"
        ],
        "outputs": [
            {
                "uxid": "140f81cdbac057e1559e94a070dd25f14b0212e3cb16389d750507c7f42e5406",
                "dst": "fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B",
                "coins": "1",
                "hours": 1
            }
        ]
    }
}
`

var btcTxStr = `
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "coin_type": "bitcoin",
  "tx": {
    "btc": {
      "txid": "69be3a3b98541e609f5a4935f94c92012d2b3e3437e9508770ba2257f532142f",
      "version": 1,
      "locktime": 0,
      "vin": [
        {
          "txid": "069f1968925c437c9fca2e567afd36d36ba2e8d0e55b25b18bc6b2c49438ea32",
          "vout": 2,
          "scriptSig": {
            "asm": "3045022100dd4e1b960726e3d3d205cb5ef4d92b3e04f3839757606800ed662069a841ffdc02203f68723bbbf9800d16555ace1ef2f46e65c2a6341643f3c5bf84158b108e6d5d[ALL] 03eb8b81f8ebc988c61d3cc4c4ac3d546b02a4994d612725e91d8d69a72045fb18",
            "hex": "483045022100dd4e1b960726e3d3d205cb5ef4d92b3e04f3839757606800ed662069a841ffdc02203f68723bbbf9800d16555ace1ef2f46e65c2a6341643f3c5bf84158b108e6d5d012103eb8b81f8ebc988c61d3cc4c4ac3d546b02a4994d612725e91d8d69a72045fb18"
          },
          "sequence": 4294967295
        }
      ],
      "vout": [
        {
          "value": "0.35601309",
          "n": 0,
          "scriptPubkey": {
            "asm": "OP_HASH160 bfc03379d17dd1e918a026b76cde472bea7ac726 OP_EQUAL",
            "hex": "a914bfc03379d17dd1e918a026b76cde472bea7ac72687",
            "type": "scripthash",
            "addresses": [
              "3KAuEYkuJQw1Ad2GzWjfC7V5XoL2fCqjGN"
            ]
          }
        }
      ],
      "blockhash": "0000000000000000021f3adb9ce12e3c70a42cfb6b7095805bee7bdefb392725",
      "confirmations": 21421,
      "time": 1471532832,
      "blocktime": 1471532832
    }
  }
}
`

func TestGetTransactionByID(t *testing.T) {
	// new bitcoin mocker
	btcM := NewCoinerMock()
	btcM.On("Name").Return("bitcoin")
	btcM.On("GetTransactionByID", "69be3a3b98541e609f5a4935f94c92012d2b3e3437e9508770ba2257f532142f").
		Return(btcTxStr, nil)
	btcM.On("GetTransactionByID", "69be3a3b98541e6").
		Return("", errors.New("invalid transaction id"))
	btcM.On("GetTransactionByID", "69be3a3b98541e609f5a4935f94c92012d2b3e3437e9508770ba2257f532142d").
		Return("", errors.New("not found"))

	// new skycoin mocker
	skyM := NewCoinerMock()
	skyM.On("Name").Return("skycoin")
	skyM.On("GetTransactionByID", "367fc68cd78adc5ed5361f9cd982289f4815da6db5a9f0bdb6c59cf463018b00").
		Return(skyTxStr, nil)
	skyM.On("GetTransactionByID", "b1481d").Return("", errors.New("invalid transaction id"))
	skyM.On("GetTransactionByID", "b1481d614ffcc27408fe2131198d9d2821c78601a0aa23d8e9965b2a5196edc1").
		Return("", errors.New("not found\n"))

	// new mzcoin mocker
	mzM := NewCoinerMock()
	mzM.On("Name").Return("mzcoin")
	mzM.On("GetTransactionByID", "367fc68cd78adc5ed5361f9cd982289f4815da6db5a9f0bdb6c59cf463018b00").
		Return(skyTxStr, nil)
	mzM.On("GetTransactionByID", "b1481d").Return("", errors.New("invalid transaction id"))
	mzM.On("GetTransactionByID", "b1481d614ffcc27408fe2131198d9d2821c78601a0aa23d8e9965b2a5196edc1").
		Return("", errors.New("not found\n"))

	initConfig(&Config{}, btcM, skyM, mzM)

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
				"367fc68cd78adc5ed5361f9cd982289f4815da6db5a9f0bdb6c59cf463018b00",
			},
			"aced4e58f22774056d2419d41f52c71920211af72c596bb5f8fd222baa41b586",
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
				"367fc68cd78adc5ed5361f9cd982289f4815da6db5a9f0bdb6c59cf463018b00",
			},
			"aced4e58f22774056d2419d41f52c71920211af72c596bb5f8fd222baa41b586",
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
		got, err := GetTransactionByID(tt.args.coinType, tt.args.txid)
		if !reflect.DeepEqual(err, tt.wantErr) {
			t.Errorf("%q. GetTransactionByID() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !strings.Contains(got, tt.contain) {
			t.Errorf("%q. GetTransactionByID() = %v, want %v", tt.name, got, tt.contain)
		}
	}
}

var outStr = `{
    "uxid": "a57c038591f862b8fada57e496ef948183b153348d7932921f865a8541a477c5",
    "time": 1477037552,
    "src_block_seq": 443,
    "src_tx": "b8ca61c0788bd711c89563f9bc60add172ee01b543ea5dcb1955c51bbfcbbaa2",
    "owner_address": "cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW",
    "coins": 1000000,
    "hours": 7,
    "spent_block_seq": 450,
    "spent_tx": "b1481d614ffcc27408fe2131198d9d2821c78601a0aa23d8e9965b2a5196edc0"
}`

func TestGetOutputByID(t *testing.T) {
	skyM := NewCoinerMock()
	skyM.On("Name").Return("skycoin")
	skyM.On("GetOutputByID", "a57c038591f862b8fada57e496ef948183b153348d7932921f865a8541a477c5").
		Return(outStr, nil)
	skyM.On("GetOutputByID", "a57c038").Return("", errors.New("invalid output hash, encoding/hex: odd length hex string"))
	skyM.On("GetOutputByID", "a57c038591f862b8fada57e496ef948183b153348d7932921f865a8541a477c9").
		Return("", errors.New("not found\n"))

	initConfig(&Config{}, skyM)

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
	}
	for _, tt := range tests {
		got, err := GetOutputByID(tt.coinType, tt.hash)
		if !reflect.DeepEqual(err, tt.wantErr) {
			t.Errorf("%q. GetOutputByHash() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !strings.Contains(got, tt.contain) {
			t.Errorf("%q. GetOutputByHash() = %v, want contains %v", tt.name, got, tt.contain)
		}
	}
}

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
		got, err := ValidateAddress(tt.args.coinType, tt.args.addr)
		if !reflect.DeepEqual(err, tt.wantErr) {
			t.Errorf("%q. ValidateAddress() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. ValidateAddress() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestNewSeed(t *testing.T) {
	sd := NewSeed()
	ss := strings.Split(sd, " ")
	if len(ss) != 12 {
		t.Fatal("error seed")
	}
}

func TestGetWalletBalance(t *testing.T) {
	tmpDir, teardown, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	skyAddressSet := []string{"2YyLVUMwjNCRZT5mBGmF13wS8yXe79eqEtu", "rRudryiBMr9zMXhb1mhZ9VwKsNVdJPGUHP"}

	skyM := NewCoinerMock()
	skyM.On("Name").Return("skycoin")
	skyM.On("GetBalance", skyAddressSet).Return(uint64(10e6), nil)

	initConfig(&Config{WalletDirPath: tmpDir}, skyM)

	id, err := NewWallet("skycoin", "123")
	if err != nil {
		t.Fatal(err)
	}

	NewAddress(id, 2)

	_, err = wallet.GetAddresses(id)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		coinType string
		wltID    string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"normal",
			args{
				"skycoin",
				"skycoin_123",
			},
			`{"balance":10000000}`,
			false,
		},
	}
	for _, tt := range tests {
		got, err := GetWalletBalance(tt.args.coinType, tt.args.wltID)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. GetWalletBalance() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. GetWalletBalance() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
