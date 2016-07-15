package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/stretchr/testify/assert"
)

var pubkey string = "02c0a2e523be9234028874a08d001d422a1a191af910b8b4c315ab7fd59223726c"
var errPubkey string = "02c0a2e523be9234028874a08d001d422a1a191af910b8b4c315ab7fd59223726"

var server_pubkey string = "02942e46684114b35fe15218dfdc6e0d74af0446a397b8fcbf8b46fb389f756eb8"
var server_seckey string = "38d010a84c7b9374352468b41b076fa585d7dfac67ac34adabe2bbba4f4f6257"

var cli_pubkey string = "025a8a0807eb20c5f6b18e62bf078ebec5b78383ed98be370d3f427969e32d490a"
var cli_seckey string = "c8f9ab54a22c5cee1c5b76dde7437db4a4f4e5555b190eb70e1c9f328740834d"

// TestCreateAccountSuccess must success.
func TestCreateAccountSuccess(t *testing.T) {
	svr := FakeServer{
		A: &FakeAccount{
			ID:      cli_pubkey,
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: uint64(0),
		},
		Seckey: cipher.MustSecKeyFromHex(server_seckey),
	}

	r := CreateAccountRequest{
		Pubkey: cli_pubkey,
	}

	sp := cipher.MustPubKeyFromHex(server_pubkey)
	cs := cipher.MustSecKeyFromHex(cli_seckey)

	key := cipher.ECDH(sp, cs)
	req := MustToContentRequest(r, cli_pubkey, key).MustJson()
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/accounts", bytes.NewBuffer(req)))
	assert.Equal(t, w.Code, 201)
}

func TestCreateAccountBadRequest(t *testing.T) {
	svr := FakeServer{
		A: &FakeAccount{
			ID:      cli_pubkey,
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: uint64(0),
		},
	}
	jsonStr := fmt.Sprintf(``)
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/accounts", strings.NewReader(jsonStr)))
	assert.Equal(t, w.Code, 400)
}

func TestCreateAccountBadPubkey(t *testing.T) {
	svr := FakeServer{
		A: &FakeAccount{
			ID:      cli_pubkey,
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: uint64(0),
		},
		Seckey: cipher.MustSecKeyFromHex(server_seckey),
	}

	r := CreateAccountRequest{
		Pubkey: errPubkey,
	}
	sp := cipher.MustPubKeyFromHex(server_pubkey)
	cs := cipher.MustSecKeyFromHex(cli_seckey)
	key := cipher.ECDH(sp, cs)

	req := MustToContentRequest(r, cli_pubkey, key).MustJson()
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/accounts", bytes.NewBuffer(req)))
	assert.Equal(t, w.Code, 400)
}

func TestGetDepositAddress(t *testing.T) {
	svr := FakeServer{
		A: &FakeAccount{
			ID:      cli_pubkey,
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: uint64(0),
		},
		Seckey: cipher.MustSecKeyFromHex(server_seckey),
	}

	r := DepositAddressRequest{
		AccountID: cli_pubkey,
		CoinType:  "bitcoin",
	}

	sp := cipher.MustPubKeyFromHex(server_pubkey)
	cs := cipher.MustSecKeyFromHex(cli_seckey)

	key := cipher.ECDH(sp, cs)
	req := MustToContentRequest(r, cli_pubkey, key).MustJson()
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/deposit_address", bytes.NewBuffer(req)))
	assert.Equal(t, w.Code, 201)
}

func TestGetDepositAddressBadCoinType(t *testing.T) {
	svr := FakeServer{
		A: &FakeAccount{
			ID:      cli_pubkey,
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: uint64(0),
		},
		Seckey: cipher.MustSecKeyFromHex(server_seckey),
	}

	r := DepositAddressRequest{
		AccountID: cli_pubkey,
		CoinType:  "abc",
	}

	sp := cipher.MustPubKeyFromHex(server_pubkey)
	cs := cipher.MustSecKeyFromHex(cli_seckey)

	key := cipher.ECDH(sp, cs)
	req := MustToContentRequest(r, cli_pubkey, key).MustJson()
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/deposit_address", bytes.NewBuffer(req)))
	assert.Equal(t, w.Code, 400)
}

func TestGetDepositAddressIDNotExist(t *testing.T) {
	svr := FakeServer{
		A: &FakeAccount{
			ID:      cli_pubkey,
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: uint64(0),
		},
		Seckey: cipher.MustSecKeyFromHex(server_seckey),
	}

	r := DepositAddressRequest{
		AccountID: pubkey,
		CoinType:  "bitcoin",
	}

	sp := cipher.MustPubKeyFromHex(server_pubkey)
	cs := cipher.MustSecKeyFromHex(cli_seckey)

	key := cipher.ECDH(sp, cs)
	req := MustToContentRequest(r, cli_pubkey, key).MustJson()
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/deposit_address", bytes.NewBuffer(req)))
	assert.Equal(t, w.Code, 404)
}

func TestGetDepositAddressBadAccountID(t *testing.T) {
	svr := FakeServer{
		A: &FakeAccount{
			ID:      cli_pubkey,
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: uint64(0),
		},
		Seckey: cipher.MustSecKeyFromHex(server_seckey),
	}

	r := DepositAddressRequest{
		AccountID: errPubkey,
		CoinType:  "bitcoin",
	}

	sp := cipher.MustPubKeyFromHex(server_pubkey)
	cs := cipher.MustSecKeyFromHex(cli_seckey)

	key := cipher.ECDH(sp, cs)
	req := MustToContentRequest(r, cli_pubkey, key).MustJson()
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/deposit_address", bytes.NewBuffer(req)))
	assert.Equal(t, w.Code, 400)
}

func TestGetbalance(t *testing.T) {
	svr := FakeServer{
		A: &FakeAccount{
			ID:      cli_pubkey,
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: uint64(0),
		},
		Seckey: cipher.MustSecKeyFromHex(server_seckey),
	}

	r := GetBalanceRequest{
		AccountID: cli_pubkey,
		CoinType:  "bitcoin",
	}

	sp := cipher.MustPubKeyFromHex(server_pubkey)
	cs := cipher.MustSecKeyFromHex(cli_seckey)

	key := cipher.ECDH(sp, cs)
	req := MustToContentRequest(r, cli_pubkey, key).MustJson()
	w := MockServer(&svr, HttpRequestCase("GET", "/api/v1/account/balance", bytes.NewBuffer(req)))
	assert.Equal(t, w.Code, 200)
}

func PrintResponse(w *httptest.ResponseRecorder) {
	d, _ := ioutil.ReadAll(w.Body)
	fmt.Println(string(d))
}
