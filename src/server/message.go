package server

import "time"

type AuthRequest struct {
	Pubkey string `json:"pubkey"`
}

type AuthResponse struct {
	Pubkey string `json:"pubkey"`
	Nonce  string `json:"nonce"`
}

type ContentRequest struct {
	AccountID string `json:"account_id"`
	Data      []byte `json:"data"`
}

type ContentResponse struct {
	Success bool   `json:"success"`
	Nonce   string `json:"nonce"`
	Data    []byte `json:"data"`
}

type CreateAccountRequest struct {
	Pubkey string `json:"pubkey"`
	Data   []byte `json:"data"`
}

type CreateAccountResponse struct {
	Succress  bool      `json:"success"`
	AccountID string    `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
}

type DepositAddressRequest struct {
	AccountID string `json:"account_id"`
	CoinType  string `json:"coin_type"`
}

type DepositAddressResponse struct {
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
	Tx        string `json:"tx"` // signed transaction
}
