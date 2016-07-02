package server

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/codahale/chacha20"
	"github.com/skycoin/skycoin-exchange/src/server/account"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/stretchr/testify/assert"
)

// FakeAccount for mocking various account state.
type FakeAccount struct {
	ID      string
	WltID   string
	Addr    string
	Nk      account.NonceKey
	Balance account.Balance
}

// FakeServer for mocking various server status.
type FakeServer struct {
	A account.Accounter
}

func (fa FakeAccount) GetWalletID() string {
	return fa.WltID
}

func (fa FakeAccount) GetAccountID() account.AccountID {
	d, err := cipher.PubKeyFromHex(fa.ID)
	if err != nil {
		panic(err)
	}
	return account.AccountID(d)
}

func (fa FakeAccount) GetNewAddress(ct wallet.CoinType) string {
	return fa.Addr
}

func (fa FakeAccount) GetBalance(ct wallet.CoinType) (account.Balance, error) {
	return fa.Balance, nil
}

func (fa *FakeAccount) SetNonceKey(nk account.NonceKey) {
	fa.Nk = nk
}

func (fa FakeAccount) GetNonceKey() account.NonceKey {
	return fa.Nk
}

func (fa FakeAccount) Encrypt(r io.Reader) ([]byte, error) {
	d, err := ioutil.ReadAll(r)
	if err != nil {
		return []byte{}, err
	}

	data := make([]byte, len(d))
	c, err := chacha20.New(fa.Nk.Key, fa.Nk.Nonce)
	if err != nil {
		return []byte{}, err
	}
	c.XORKeyStream(data, d)
	return data, nil
}

func (fa FakeAccount) Decrypt(r io.Reader) ([]byte, error) {
	d, err := ioutil.ReadAll(r)
	if err != nil {
		return []byte{}, err
	}

	data := make([]byte, len(d))
	c, err := chacha20.New(fa.Nk.Key, fa.Nk.Nonce)
	if err != nil {
		return []byte{}, err
	}
	c.XORKeyStream(data, d)
	return data, nil
}

func (fa FakeAccount) IsExpired() bool {
	d := time.Now().Unix() - fa.Nk.Expire_at.Unix()
	return d >= 0
}

func (fs *FakeServer) CreateAccountWithPubkey(pk cipher.PubKey) (account.Accounter, error) {
	if fs.A.GetWalletID() == "" {
		return nil, fmt.Errorf("create wallet failed")
	}
	return fs.A, nil
}

func (fs *FakeServer) GetAccount(id account.AccountID) (account.Accounter, error) {
	if fs.A != nil && fs.A.GetAccountID() == id {
		return fs.A, nil
	}
	return nil, errors.New("account not found")
}

func (fs *FakeServer) Run() {

}

func (fs FakeServer) GetNonceKeyLifetime() time.Duration {
	return time.Second * time.Duration(10*60)
}

var pubkey string = "02c0a2e523be9234028874a08d001d422a1a191af910b8b4c315ab7fd59223726c"
var errPubkey string = "02c0a2e523be9234028874a08d001d422a1a191af910b8b4c315ab7fd59223726"

// TestCreateAccountSuccess must success.
func TestCreateAccountSuccess(t *testing.T) {
	svr := FakeServer{
		A: &FakeAccount{
			ID:      pubkey,
			WltID:   "test.wlt",
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: account.Balance(0),
		},
	}
	jsonStr := fmt.Sprintf(`{"pubkey": "%s"}`, pubkey)
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/accounts", strings.NewReader(jsonStr)))
	assert.Equal(t, w.Code, 201)
}

func TestCreateAccountBadRequest(t *testing.T) {
	svr := FakeServer{
		A: &FakeAccount{
			ID:      pubkey,
			WltID:   "test.wlt",
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: account.Balance(0),
		},
	}
	jsonStr := fmt.Sprintf(``)
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/accounts", strings.NewReader(jsonStr)))
	assert.Equal(t, w.Code, 400)
}

// TestCreateAccountInvalidPubkey invalid pubkey.
func TestCreateAccountInvalidPubkey(t *testing.T) {
	svr := FakeServer{
		A: &FakeAccount{
			ID:      pubkey,
			WltID:   "test.wlt",
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: account.Balance(0),
		},
	}

	jsonStr := fmt.Sprintf(`{"pubkey": "%s"}`, errPubkey)
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/accounts", strings.NewReader(jsonStr)))
	assert.Equal(t, w.Code, 400)
}

// TestCreateAccountFaildBindWallet test case of creating wallet faild.
func TestCreateAccountFaildBindWallet(t *testing.T) {
	svr := FakeServer{
		A: &FakeAccount{
			ID:      "02c0a2e523be9234028874a08d001d422a1a191af910b8b4c315ab7fd59223726c",
			WltID:   "",
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: account.Balance(0),
		},
	}

	jsonStr := fmt.Sprintf(`{"pubkey": "%s"}`, pubkey)
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/accounts", strings.NewReader(jsonStr)))
	assert.Equal(t, w.Code, 501)
}

func TestAuth(t *testing.T) {
	p, s := cipher.GenerateKeyPair()
	pubkey := fmt.Sprintf("%x", p)
	svr := FakeServer{
		A: &FakeAccount{
			ID:      pubkey,
			WltID:   "",
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: account.Balance(0),
		},
	}

	jsonStr := fmt.Sprintf(`{"pubkey": "%s"}`, pubkey)
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/authorization", strings.NewReader(jsonStr)))
	assert.Equal(t, w.Code, 200)
	d, err := ioutil.ReadAll(w.Body)
	assert.Nil(t, err)
	ar := AuthResponse{}
	json.Unmarshal(d, &ar)
	pk, err := cipher.PubKeyFromHex(ar.Pubkey)
	assert.Nil(t, err)
	key := cipher.ECDH(pk, s)
	assert.Equal(t, bytes.Equal(svr.A.GetNonceKey().Key, key), true)
}

func TestAuthInvalidPubkey(t *testing.T) {
	svr := FakeServer{
		A: &FakeAccount{
			ID:      pubkey,
			WltID:   "",
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: account.Balance(0),
		},
	}

	jsonStr := fmt.Sprintf(`{"pubkey": "%s"}`, errPubkey)
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/authorization", strings.NewReader(jsonStr)))
	assert.Equal(t, w.Code, 400)
}

func TestAuthUnknowID(t *testing.T) {
	p, _ := cipher.GenerateKeyPair()
	client_pubkey := fmt.Sprintf("%x", p)
	svr := FakeServer{
		A: &FakeAccount{
			ID:      pubkey,
			WltID:   "",
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: account.Balance(0),
		},
	}

	jsonStr := fmt.Sprintf(`{"pubkey": "%s"}`, client_pubkey)
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/authorization", strings.NewReader(jsonStr)))
	assert.Equal(t, w.Code, 404)
}

func TestAuthReqRequest(t *testing.T) {
	p, _ := cipher.GenerateKeyPair()
	client_pubkey := fmt.Sprintf("%x", p)
	svr := FakeServer{
		A: &FakeAccount{
			ID:      pubkey,
			WltID:   "",
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: account.Balance(0),
		},
	}

	jsonStr := fmt.Sprintf(`{"pubkey": "%s"}`, client_pubkey)
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/authorization", strings.NewReader(jsonStr[:10])))
	assert.Equal(t, w.Code, 400)
}

func TestAuthReqExpire(t *testing.T) {
	p, _ := cipher.GenerateKeyPair()
	pubkey := fmt.Sprintf("%x", p)
	svr := FakeServer{
		A: &FakeAccount{
			ID:    pubkey,
			WltID: "test.wlt",
			Addr:  "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			// the account will expire after 1 second
			Nk: account.NonceKey{Key: key, Nonce: nonce, Expire_at: time.Now().Add(time.Millisecond * time.Duration(200))},
		},
	}

	// send get deposite address request
	dar := DepositAddressRequest{
		AccountID: pubkey, // NOTE: different key
		CoinType:  "bitcoin",
	}

	ct := dar.MustToContentRequest(key, nonce)
	ctd, err := json.Marshal(ct)
	assert.Nil(t, err)

	time.Sleep(time.Millisecond * time.Duration(300))
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/account/deposit_address", bytes.NewBuffer(ctd)))
	assert.Equal(t, w.Code, 401)
}

func TestUnAuth(t *testing.T) {
	svr := FakeServer{
		A: &FakeAccount{
			ID:      pubkey,
			WltID:   "",
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: account.Balance(0),
		},
	}

	dar := DepositAddressRequest{
		AccountID: pubkey,
		CoinType:  "bitcoin",
	}
	cr := dar.MustToContentRequest(make([]byte, 32), make([]byte, 8))
	d, _ := json.Marshal(cr)
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/account/deposit_address", bytes.NewBuffer(d)))
	assert.Equal(t, w.Code, 401)
}

var key []byte = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
var nonce []byte = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

func TestAuthReqBadRequest(t *testing.T) {
	p, _ := cipher.GenerateKeyPair()
	pubkey := fmt.Sprintf("%x", p)
	svr := FakeServer{
		A: &FakeAccount{
			ID:    pubkey,
			WltID: "test.wlt",
			Addr:  "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Nk:    account.NonceKey{Key: key, Nonce: nonce, Expire_at: time.Now().Add(time.Second * time.Duration(10))},
		},
	}

	// send get deposite address request
	dar := DepositAddressRequest{
		AccountID: pubkey,
		CoinType:  "bitcoin",
	}

	ct := dar.MustToContentRequest(key, nonce)
	ctd, err := json.Marshal(ct)
	assert.Nil(t, err)
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/account/deposit_address", bytes.NewBuffer(ctd[:10]))) // NOTE: Bad request
	assert.Equal(t, w.Code, 400)
}

func TestAuthReqBadPubkey(t *testing.T) {
	p, _ := cipher.GenerateKeyPair()
	pubkey := fmt.Sprintf("%x", p)
	svr := FakeServer{
		A: &FakeAccount{
			ID:    pubkey,
			WltID: "test.wlt",
			Addr:  "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Nk:    account.NonceKey{Key: key, Nonce: nonce, Expire_at: time.Now().Add(time.Second * time.Duration(10))},
		},
	}

	// send get deposite address request
	dar := DepositAddressRequest{
		AccountID: errPubkey, // NOTE: bad key
		CoinType:  "bitcoin",
	}

	ct := dar.MustToContentRequest(key, nonce)
	ctd, err := json.Marshal(ct)
	assert.Nil(t, err)
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/account/deposit_address", bytes.NewBuffer(ctd)))
	assert.Equal(t, w.Code, 400)
}

func TestAuthReqBadID(t *testing.T) {
	p, _ := cipher.GenerateKeyPair()
	cli_pubkey := fmt.Sprintf("%x", p)
	svr := FakeServer{
		A: &FakeAccount{
			ID:    pubkey,
			WltID: "test.wlt",
			Addr:  "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
		},
	}

	// send get deposite address request
	dar := DepositAddressRequest{
		AccountID: cli_pubkey, // NOTE: different key
		CoinType:  "bitcoin",
	}

	ct := dar.MustToContentRequest(key, nonce)
	ctd, err := json.Marshal(ct)
	assert.Nil(t, err)
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/account/deposit_address", bytes.NewBuffer(ctd)))
	assert.Equal(t, w.Code, 404)
}

func TestAuthReqBadNonceKey(t *testing.T) {
	p, _ := cipher.GenerateKeyPair()
	pubkey := fmt.Sprintf("%x", p)
	svr := FakeServer{
		A: &FakeAccount{
			ID:    pubkey,
			WltID: "test.wlt",
			Addr:  "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Nk:    account.NonceKey{Key: key, Nonce: nonce, Expire_at: time.Now().Add(time.Second * time.Duration(10))},
		},
	}

	// send get deposite address request
	dar := DepositAddressRequest{
		AccountID: pubkey,
		CoinType:  "bitcoin",
	}
	errKey := make([]byte, 32)
	copy(errKey, key)
	errKey[0] = 0x98
	ct := dar.MustToContentRequest(errKey, nonce)
	ctd, err := json.Marshal(ct)
	assert.Nil(t, err)
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/account/deposit_address", bytes.NewBuffer(ctd)))
	assert.Equal(t, w.Code, 401)
}

func TestCreateAddress(t *testing.T) {
	p, s := cipher.GenerateKeyPair()
	pubkey := fmt.Sprintf("%x", p)
	svr := FakeServer{
		A: &FakeAccount{
			ID:      pubkey,
			WltID:   "test.wlt",
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: account.Balance(0),
		},
	}

	// auth first.
	jsonStr := fmt.Sprintf(`{"pubkey": "%s"}`, pubkey)
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/authorization", strings.NewReader(jsonStr)))
	assert.Equal(t, w.Code, 200)
	d, err := ioutil.ReadAll(w.Body)
	assert.Nil(t, err)
	// get the key.
	ar := AuthResponse{}
	err = json.Unmarshal(d, &ar)
	assert.Nil(t, err)

	// generate nonce key
	spk, err := cipher.PubKeyFromHex(ar.Pubkey)
	assert.Nil(t, err)
	key := cipher.ECDH(spk, s)

	// send get deposite address request
	dar := DepositAddressRequest{
		AccountID: pubkey,
		CoinType:  "bitcoin",
	}

	nonce, err := hex.DecodeString(ar.Nonce)
	assert.Nil(t, err)
	ct := dar.MustToContentRequest(key, nonce)
	ctd, err := json.Marshal(ct)
	assert.Nil(t, err)

	w = MockServer(&svr, HttpRequestCase("POST", "/api/v1/account/deposit_address", bytes.NewBuffer(ctd)))
	assert.Equal(t, w.Code, 201)
}

func PrintResponse(w *httptest.ResponseRecorder) {
	d, _ := ioutil.ReadAll(w.Body)
	fmt.Println(string(d))
}
