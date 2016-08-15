package rpclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/rpclient/api"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/stretchr/testify/assert"
)

// mux.Handle("/api/v1/account/balance", GetBalance(se))
// mux.Handle("/api/v1/account/withdrawal", Withdraw(se))
//
// // order handlers
// mux.Handle("/api/v1/account/order/bid", CreateBidOrder(se))
// mux.Handle("/api/v1/account/order/ask", CreateAskOrder(se))
// mux.Handle("/api/v1/orders/bid", GetBidOrders(se))
// mux.Handle("/api/v1/orders/ask", GetAskOrders(se))

var (
	se = New(Config{
		ApiRoot:    "localhost:8080",
		ServPubkey: cipher.MustPubKeyFromHex("02942e46684114b35fe15218dfdc6e0d74af0446a397b8fcbf8b46fb389f756eb8"),
	})
)

func TestGetCoins(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost/api/v1/coins", nil)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	api.GetCoins(se)(w, req)
	res := pp.EmptyRes{}
	err = json.NewDecoder(w.Body).Decode(&res)
	assert.Nil(t, err)
	assert.Equal(t, res.Result.GetSuccess(), true)
}

func createAccount() (string, string, error) {
	req, err := http.NewRequest("POST", "http://localhost/api/v1/accounts", nil)
	if err != nil {
		return "", "", err
	}

	w := httptest.NewRecorder()
	api.CreateAccount(se)(w, req)
	res := struct {
		Result    pp.Result `json:"result"`
		ID        string    `json:"account_id"`
		Key       string    `json:"key"`
		CreatedAt int64     `json:"created_at"`
	}{}
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		return "", "", err
	}
	if !res.Result.GetSuccess() {
		return "", "", errors.New(res.Result.GetReason())
	}

	return res.ID, res.Key, nil
}

func TestCreatAccount(t *testing.T) {
	_, _, err := createAccount()
	assert.Nil(t, err)
}

func TestDeposit(t *testing.T) {
	id, key, err := createAccount()
	assert.Nil(t, err)
	url := fmt.Sprintf("http://localhost/api/v1/account/deposit_address?coin_type=%s&id=%s&key=%s",
		"bitcoin", id, key)
	req, err := http.NewRequest("POST", url, nil)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	api.GetNewAddress(se)(w, req)
	res := pp.EmptyRes{}
	err = json.NewDecoder(w.Body).Decode(&res)
	assert.Nil(t, err)
	assert.Equal(t, res.Result.GetSuccess(), true)
	log.Println(res.Result)
}

func TestGetBalance(t *testing.T) {
	id, key, err := createAccount()
	assert.Nil(t, err)
	url := fmt.Sprintf("http://localhost/api/v1/account/balance?coin_type=%s&id=%s&key=%s",
		"bitcoin", id, key)
	req, err := http.NewRequest("GET", url, nil)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	api.GetBalance(se)(w, req)
	res := pp.EmptyRes{}
	err = json.NewDecoder(w.Body).Decode(&res)
	assert.Nil(t, err)
	assert.Equal(t, res.Result.GetSuccess(), true)
}

func TestGetBidOrders(t *testing.T) {
	url := fmt.Sprintf("http://localhost/api/v1/account/order/bid?coin_pair=bitcoin/skycoin&start=0&end=10")
	req, err := http.NewRequest("GET", url, nil)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	api.GetBidOrders(se)(w, req)
	res := pp.EmptyRes{}
	err = json.NewDecoder(w.Body).Decode(&res)
	assert.Nil(t, err)
	assert.Equal(t, res.Result.GetSuccess(), true)
}

func TestGetAskOrders(t *testing.T) {
	url := fmt.Sprintf("http://localhost/api/v1/account/order/ask?coin_pair=bitcoin/skycoin&start=0&end=10")
	req, err := http.NewRequest("GET", url, nil)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	api.GetAskOrders(se)(w, req)
	res := pp.EmptyRes{}
	err = json.NewDecoder(w.Body).Decode(&res)
	assert.Nil(t, err)
	assert.Equal(t, res.Result.GetSuccess(), true)
}
