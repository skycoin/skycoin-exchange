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
			WltID:   "test.wlt",
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
			WltID:   "test.wlt",
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
			WltID:   "test.wlt",
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: uint64(0),
		},
		Seckey: cipher.MustSecKeyFromHex(server_seckey),
	}

	car := CreateAccountRequest{
		Pubkey: errPubkey,
	}
	sp := cipher.MustPubKeyFromHex(server_pubkey)
	cs := cipher.MustSecKeyFromHex(cli_seckey)
	key := cipher.ECDH(sp, cs)

	cr := MustToContentRequest(car, cli_pubkey, key)
	cr.AccountID = cli_pubkey
	req := cr.MustJson()

	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/accounts", bytes.NewBuffer(req)))
	assert.Equal(t, w.Code, 400)
}

func TestCreateAccountServerFail(t *testing.T) {
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
	assert.Equal(t, w.Code, 501)
}

func TestGetDepositAddress(t *testing.T) {
	svr := FakeServer{
		A: &FakeAccount{
			ID:      cli_pubkey,
			WltID:   "test.wlt",
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
			WltID:   "test.wlt",
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
			WltID:   "test.wlt",
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
			WltID:   "test.wlt",
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

// TestCreateAccountInvalidPubkey invalid pubkey.
// func TestCreateAccountInvalidPubkey(t *testing.T) {
// 	svr := FakeServer{
// 		A: &FakeAccount{
// 			ID:      cli_pubkey,
// 			WltID:   "test.wlt",
// 			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
// 			Balance: uint64(0),
// 		},
// 		Seckey: cipher.MustSecKeyFromHex(server_seckey),
// 	}
//
// 	sp := cipher.MustPubKeyFromHex(server_pubkey)
// 	cs := cipher.MustSecKeyFromHex(cli_seckey)
//
// 	r := CreateAccountRequest{
// 		Pubkey: cli_pubkey,
// 	}
//
// 	key := cipher.ECDH(sp, cs)
// 	req := r.MustToContentRequest(cli_pubkey, key).MustJson()
// 	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/accounts", bytes.NewBuffer(req)))
// 	assert.Equal(t, w.Code, 201)
// }

//
// // TestCreateAccountFaildBindWallet test case of creating wallet faild.
// func TestCreateAccountFaildBindWallet(t *testing.T) {
// 	svr := FakeServer{
// 		A: &FakeAccount{
// 			ID:      "02c0a2e523be9234028874a08d001d422a1a191af910b8b4c315ab7fd59223726c",
// 			WltID:   "",
// 			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
// 			Balance: uint64(0),
// 		},
// 	}
//
// 	jsonStr := fmt.Sprintf(`{"pubkey": "%s"}`, pubkey)
// 	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/accounts", strings.NewReader(jsonStr)))
// 	assert.Equal(t, w.Code, 501)
// }
//

//
// func TestAuthInvalidPubkey(t *testing.T) {
// 	svr := FakeServer{
// 		A: &FakeAccount{
// 			ID:      pubkey,
// 			WltID:   "",
// 			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
// 			Balance: uint64(0),
// 		},
// 	}
//
// 	jsonStr := fmt.Sprintf(`{"pubkey": "%s"}`, errPubkey)
// 	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/authorization", strings.NewReader(jsonStr)))
// 	assert.Equal(t, w.Code, 400)
// }
//
// func TestAuthUnknowID(t *testing.T) {
// 	p, _ := cipher.GenerateKeyPair()
// 	client_pubkey := fmt.Sprintf("%x", p)
// 	svr := FakeServer{
// 		A: &FakeAccount{
// 			ID:      pubkey,
// 			WltID:   "",
// 			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
// 			Balance: uint64(0),
// 		},
// 	}
//
// 	jsonStr := fmt.Sprintf(`{"pubkey": "%s"}`, client_pubkey)
// 	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/authorization", strings.NewReader(jsonStr)))
// 	assert.Equal(t, w.Code, 404)
// }
//
// func TestAuthReqRequest(t *testing.T) {
// 	p, _ := cipher.GenerateKeyPair()
// 	client_pubkey := fmt.Sprintf("%x", p)
// 	svr := FakeServer{
// 		A: &FakeAccount{
// 			ID:      pubkey,
// 			WltID:   "",
// 			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
// 			Balance: uint64(0),
// 		},
// 	}
//
// 	jsonStr := fmt.Sprintf(`{"pubkey": "%s"}`, client_pubkey)
// 	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/authorization", strings.NewReader(jsonStr[:10])))
// 	assert.Equal(t, w.Code, 400)
// }
//
// func TestAuthReqExpire(t *testing.T) {
// 	p, _ := cipher.GenerateKeyPair()
// 	pubkey := fmt.Sprintf("%x", p)
// 	svr := FakeServer{
// 		A: &FakeAccount{
// 			ID:    pubkey,
// 			WltID: "test.wlt",
// 			Addr:  "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
// 			// the account will expire after 1 second
// 			Nk: account.NonceKey{Key: key, Nonce: nonce, Expire_at: time.Now().Add(time.Millisecond * time.Duration(200))},
// 		},
// 	}
//
// 	// send get deposite address request
// 	dar := DepositAddressRequest{
// 		AccountID: pubkey, // NOTE: different key
// 		CoinType:  "bitcoin",
// 	}
//
// 	ct := dar.MustToContentRequest(key, nonce)
// 	ctd, err := json.Marshal(ct)
// 	assert.Nil(t, err)
//
// 	time.Sleep(time.Millisecond * time.Duration(300))
// 	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/account/deposit_address", bytes.NewBuffer(ctd)))
// 	assert.Equal(t, w.Code, 401)
// }
//
// func TestUnAuth(t *testing.T) {
// 	svr := FakeServer{
// 		A: &FakeAccount{
// 			ID:      pubkey,
// 			WltID:   "",
// 			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
// 			Balance: uint64(0),
// 		},
// 	}
//
// 	dar := DepositAddressRequest{
// 		AccountID: pubkey,
// 		CoinType:  "bitcoin",
// 	}
// 	cr := dar.MustToContentRequest(make([]byte, 32), make([]byte, 8))
// 	d, _ := json.Marshal(cr)
// 	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/account/deposit_address", bytes.NewBuffer(d)))
// 	assert.Equal(t, w.Code, 401)
// }
//
// var key []byte = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
// var nonce []byte = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
//
// func TestAuthReqBadRequest(t *testing.T) {
// 	p, _ := cipher.GenerateKeyPair()
// 	pubkey := fmt.Sprintf("%x", p)
// 	svr := FakeServer{
// 		A: &FakeAccount{
// 			ID:    pubkey,
// 			WltID: "test.wlt",
// 			Addr:  "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
// 			Nk:    account.NonceKey{Key: key, Nonce: nonce, Expire_at: time.Now().Add(time.Second * time.Duration(10))},
// 		},
// 	}
//
// 	// send get deposite address request
// 	dar := DepositAddressRequest{
// 		AccountID: pubkey,
// 		CoinType:  "bitcoin",
// 	}
//
// 	ct := dar.MustToContentRequest(key, nonce)
// 	ctd, err := json.Marshal(ct)
// 	assert.Nil(t, err)
// 	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/account/deposit_address", bytes.NewBuffer(ctd[:10]))) // NOTE: Bad request
// 	assert.Equal(t, w.Code, 400)
// }
//
// func TestAuthReqBadPubkey(t *testing.T) {
// 	p, _ := cipher.GenerateKeyPair()
// 	pubkey := fmt.Sprintf("%x", p)
// 	svr := FakeServer{
// 		A: &FakeAccount{
// 			ID:    pubkey,
// 			WltID: "test.wlt",
// 			Addr:  "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
// 			Nk:    account.NonceKey{Key: key, Nonce: nonce, Expire_at: time.Now().Add(time.Second * time.Duration(10))},
// 		},
// 	}
//
// 	// send get deposite address request
// 	dar := DepositAddressRequest{
// 		AccountID: errPubkey, // NOTE: bad key
// 		CoinType:  "bitcoin",
// 	}
//
// 	ct := dar.MustToContentRequest(key, nonce)
// 	ctd, err := json.Marshal(ct)
// 	assert.Nil(t, err)
// 	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/account/deposit_address", bytes.NewBuffer(ctd)))
// 	assert.Equal(t, w.Code, 400)
// }
//
// func TestAuthReqBadID(t *testing.T) {
// 	p, _ := cipher.GenerateKeyPair()
// 	cli_pubkey := fmt.Sprintf("%x", p)
// 	svr := FakeServer{
// 		A: &FakeAccount{
// 			ID:    pubkey,
// 			WltID: "test.wlt",
// 			Addr:  "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
// 		},
// 	}
//
// 	// send get deposite address request
// 	dar := DepositAddressRequest{
// 		AccountID: cli_pubkey, // NOTE: different key
// 		CoinType:  "bitcoin",
// 	}
//
// 	ct := dar.MustToContentRequest(key, nonce)
// 	ctd, err := json.Marshal(ct)
// 	assert.Nil(t, err)
// 	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/account/deposit_address", bytes.NewBuffer(ctd)))
// 	assert.Equal(t, w.Code, 404)
// }
//
// func TestAuthReqBadNonceKey(t *testing.T) {
// 	p, _ := cipher.GenerateKeyPair()
// 	pubkey := fmt.Sprintf("%x", p)
// 	svr := FakeServer{
// 		A: &FakeAccount{
// 			ID:    pubkey,
// 			WltID: "test.wlt",
// 			Addr:  "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
// 			Nk:    account.NonceKey{Key: key, Nonce: nonce, Expire_at: time.Now().Add(time.Second * time.Duration(10))},
// 		},
// 	}
//
// 	// send get deposite address request
// 	dar := DepositAddressRequest{
// 		AccountID: pubkey,
// 		CoinType:  "bitcoin",
// 	}
// 	errKey := make([]byte, 32)
// 	copy(errKey, key)
// 	errKey[0] = 0x98
// 	ct := dar.MustToContentRequest(errKey, nonce)
// 	ctd, err := json.Marshal(ct)
// 	assert.Nil(t, err)
// 	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/account/deposit_address", bytes.NewBuffer(ctd)))
// 	assert.Equal(t, w.Code, 401)
// }
//
// func TestCreateAddress(t *testing.T) {
// 	p, s := cipher.GenerateKeyPair()
// 	pubkey := fmt.Sprintf("%x", p)
// 	svr := FakeServer{
// 		A: &FakeAccount{
// 			ID:      pubkey,
// 			WltID:   "test.wlt",
// 			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
// 			Balance: uint64(0),
// 		},
// 	}
//
// 	// auth first.
// 	jsonStr := fmt.Sprintf(`{"pubkey": "%s"}`, pubkey)
// 	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/authorization", strings.NewReader(jsonStr)))
// 	assert.Equal(t, w.Code, 200)
// 	d, err := ioutil.ReadAll(w.Body)
// 	assert.Nil(t, err)
// 	// get the key.
// 	ar := AuthResponse{}
// 	err = json.Unmarshal(d, &ar)
// 	assert.Nil(t, err)
//
// 	// generate nonce key
// 	spk, err := cipher.PubKeyFromHex(ar.Pubkey)
// 	assert.Nil(t, err)
// 	key := cipher.ECDH(spk, s)
//
// 	// send get deposite address request
// 	dar := DepositAddressRequest{
// 		AccountID: pubkey,
// 		CoinType:  "bitcoin",
// 	}
//
// 	nonce, err := hex.DecodeString(ar.Nonce)
// 	assert.Nil(t, err)
// 	ct := dar.MustToContentRequest(key, nonce)
// 	ctd, err := json.Marshal(ct)
// 	assert.Nil(t, err)
//
// 	w = MockServer(&svr, HttpRequestCase("POST", "/api/v1/account/deposit_address", bytes.NewBuffer(ctd)))
// 	assert.Equal(t, w.Code, 201)
// }
//
func PrintResponse(w *httptest.ResponseRecorder) {
	d, _ := ioutil.ReadAll(w.Body)
	fmt.Println(string(d))
}
