package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
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
	Balance account.Balance
}

// type FakeAccountManager struct {
//
// }

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

//
// GetAccount(id AccountID) (Accounter, error)
// Save()
// Load()

// func (fm *FakeAccountManager) CreateAccount() (Accounter, cipher.SecKey, error) {
// 	return nil, cipher.SecKey{}, nil
// }
//
// func (fm *FakeAccountManager) CreateAccountWithPubkey(pk cipher.PubKey) (Accounter, error) {
//
// }

func (fs *FakeServer) CreateAccountWithPubkey(pk cipher.PubKey) (account.Accounter, error) {
	if fs.A.GetWalletID() == "" {
		return nil, fmt.Errorf("create wallet failed")
	}
	return fs.A, nil
}

func (fs *FakeServer) GetAccount(id account.AccountID) (account.Accounter, error) {
	return fs.A, nil
}

func (fs *FakeServer) Run() {

}

type CaseHandler func() (*httptest.ResponseRecorder, *http.Request)

func MockServer(svr Server, fs CaseHandler) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := NewRouter(svr)
	w, r := fs()
	router.ServeHTTP(w, r)
	return w
}

func HttpRequestCase(method string, url string, body io.Reader) CaseHandler {
	return func() (*httptest.ResponseRecorder, *http.Request) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			panic(err)
		}
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return w, r
	}
}

var pubkey string = "02c0a2e523be9234028874a08d001d422a1a191af910b8b4c315ab7fd59223726c"
var errPubkey string = "02c0a2e523be9234028874a08d001d422a1a191af910b8b4c315ab7fd59223726"

// TestCreateAccountSuccess must success.
func TestCreateAccountSuccess(t *testing.T) {
	svr := FakeServer{
		A: FakeAccount{
			ID:      pubkey,
			WltID:   "test.wlt",
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: account.Balance(0),
		},
	}

	form := url.Values{}
	form.Add("pubkey", pubkey)
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/account", strings.NewReader(form.Encode())))
	assert.Equal(t, w.Code, 201)
}

// TestCreateAccountInvalidPubkey invalid pubkey.
func TestCreateAccountInvalidPubkey(t *testing.T) {
	svr := FakeServer{
		A: FakeAccount{
			ID:      pubkey,
			WltID:   "test.wlt",
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: account.Balance(0),
		},
	}

	form := url.Values{}
	form.Add("pubkey", errPubkey) // invalid pubkey
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/account", strings.NewReader(form.Encode())))
	assert.Equal(t, w.Code, 400)
}

// TestCreateAccountFaildBindWallet test case of creating wallet faild.
func TestCreateAccountFaildBindWallet(t *testing.T) {
	svr := FakeServer{
		A: FakeAccount{
			ID:      "02c0a2e523be9234028874a08d001d422a1a191af910b8b4c315ab7fd59223726c",
			WltID:   "",
			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
			Balance: account.Balance(0),
		},
	}

	form := url.Values{}
	form.Add("pubkey", pubkey) // invalid pubkey
	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/account", strings.NewReader(form.Encode())))
	assert.Equal(t, w.Code, 501)
}
