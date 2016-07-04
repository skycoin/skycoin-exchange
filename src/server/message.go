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
	Nonce   string `json:"nonce"`
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
	AccountID string `json:"account_id"`
	Tx        []byte `json:"tx"` // signed transaction
}

func MustToContentRequest(r interface{}, id string, key []byte) ContentRequest {
	d, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	nonce := cipher.RandByte(8)
	data, err := Encrypt(d, key, nonce)
	if err != nil {
		panic(err)
	}
	return ContentRequest{
		AccountID: id,
		Data:      data,
		Nonce:     nonce,
	}
}

// func (car CreateAccountRequest) MustToContentRequest(id string, key []byte) ContentRequest {
// }
//
// func (dar DepositAddressRequest) MustToContentResponse(id string, key []byte) ContentRequest {
// 	d, err := json.Marshal(dar)
// 	if err != nil {
// 		panic(err)
// 	}
// 	nonce := cipher.RandByte(8)
// 	data, err := Encrypt(d, key, nonce)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return ContentRequest{
// 		AccountID: id,
// 		Data:      data,
// 		Nonce:     nonce,
// 	}
// }

func (cr ContentRequest) MustJson() []byte {
	d, err := json.Marshal(cr)
	if err != nil {
		panic(err)
	}
	return d
}

// func (wr WithdrawResponse) MustToContentResponse(a account.Accounter) ContentResponse {
// 	d, err := json.Marshal(wr)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	data, _ := a.Encrypt(bytes.NewBuffer(d))
//
// 	return ContentResponse{
// 		Success: true,
// 		Data:    data,
// 	}
// }

// MustToContentRequest convert DepositAddressRequest to ContentRequest
func (dar DepositAddressRequest) MustToContentRequest(key []byte, nonce []byte) ContentRequest {
	d, err := json.Marshal(dar)
	if err != nil {
		panic(err)
	}

	data, err := Encrypt(d, key, nonce)
	if err != nil {
		panic(err)
	}

	return ContentRequest{
		AccountID: dar.AccountID,
		Data:      data,
	}
}
