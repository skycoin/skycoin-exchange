package server

import (
	"encoding/json"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
)

type EncryptedData struct {
	Data []byte `json:"data"`
}

// type AuthRequest struct {
// 	Pubkey string `json:"pubkey"`
// }
//
// type AuthResponse struct {
// 	Pubkey string `json:"pubkey"`
// 	Nonce  string `json:"nonce"`
// }

type ContentRequest struct {
	AccountID string `json:"account_id"`
	Nonce     []byte `json:"nonce"`
	Data      []byte `json:"data"`
}

type ContentResponse struct {
	Success bool   `json:"success"`
	Nonce   []byte `json:"nonce"`
	Data    []byte `json:"data"`
}

type CreateAccountRequest struct {
	Pubkey string `json:"pubkey"`
}

type CreateAccountResponse struct {
	Success   bool      `json:"success"`
	AccountID string    `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
}

type DepositAddressRequest struct {
	AccountID string `json:"account_id"`
	CoinType  string `json:"coin_type"`
}

type DepositAddressResponse struct {
	Success     bool   `json:"success"`
	AccountID   string `json:"account_id"`
	DepositAddr string `json:"deposit_address"`
}

type WithdrawRequest struct {
	AccountID     string `json:"account_id"`
	CoinType      string `json:"coin_type"`
	Coins         uint64 `json:"coins"`
	OutputAddress string `json:"output_address"`
}

type WithdrawResponse struct {
	Success   bool   `json:"success"`
	AccountID string `json:"account_id"`
	NewTxid   string `json:"new_txid"` // signed transaction
}

type GetBalanceRequest struct {
	AccountID string `json:"account_id"`
	CoinType  string `json:"coin_type"`
}

type GetBalanceResponse struct {
	Success   string `json:"success"`
	AccountID string `json:"account_id"`
	CoinType  string `json:"coin_type"`
	Balance   int64  `json:"balance"`
}

func MustToContentRequest(r interface{}, id string, key []byte) ContentRequest {
	d, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	nonce := cipher.RandByte(8)
	data, err := encOrDec(d, key, nonce)
	if err != nil {
		panic(err)
	}
	return ContentRequest{
		AccountID: id,
		Data:      data,
		Nonce:     nonce,
	}
}

func (cr ContentRequest) MustJson() []byte {
	d, err := json.Marshal(cr)
	if err != nil {
		panic(err)
	}
	return d
}

// MustToContentRequest convert DepositAddressRequest to ContentRequest
func (dar DepositAddressRequest) MustToContentRequest(key []byte, nonce []byte) ContentRequest {
	d, err := json.Marshal(dar)
	if err != nil {
		panic(err)
	}

	data, err := encOrDec(d, key, nonce)
	if err != nil {
		panic(err)
	}

	return ContentRequest{
		AccountID: dar.AccountID,
		Data:      data,
	}
}
