package server

import (
	"bytes"
	"fmt"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skycoin/skycoin-exchange/src/server/account"
	"github.com/skycoin/skycoin/src/cipher"
)

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

// Authorize generate a nonce pubkey/seckey pairs, do ECDH to generate
// NonceKey, store the key into the account and return the pubkey
// to client.
func Authorize(svr Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := AuthRequest{}
		if c.BindJSON(&r) == nil {
			pubkey, err := cipher.PubKeyFromHex(r.Pubkey)
			if err != nil {
				c.JSON(400, ErrorMsg{Code: 400, Error: "invalide pubkey"})
				return
			}

			a, err := svr.GetAccount(account.AccountID(pubkey))
			if err != nil {
				c.JSON(404, ErrorMsg{Code: 404, Error: err.Error()})
				return
			}

			p, s := cipher.GenerateKeyPair()
			nk := account.NonceKey{
				Key:       cipher.ECDH(pubkey, s),
				Nonce:     cipher.RandByte(8),
				Expire_at: time.Now().Add(svr.GetNonceKeyLifetime()),
			}

			// set the nonce key of the account.
			a.SetNonceKey(nk)

			c.JSON(200, AuthResponse{Pubkey: fmt.Sprintf("%x", p), Nonce: fmt.Sprintf("%x", nk.Nonce)})
			return
		}
		c.JSON(400, ErrorMsg{Code: 400, Error: "bad request"})
	}
}

// AuthRequired middleware for check the authorization of client, and
// decrypt the data, set the data in rawdata.
func AuthRequired(svr Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := ContentRequest{}
		if c.BindJSON(&r) == nil {
			// check the existence of this account.
			pubkey, err := cipher.PubKeyFromHex(r.AccountID)
			if err != nil {
				c.JSON(400, ErrorMsg{Code: 400, Error: "invalide account id"})
				c.Abort()
				return
			}

			// find the account.
			id := account.AccountID(pubkey)
			a, err := svr.GetAccount(id)
			if err != nil {
				c.JSON(404, ErrorMsg{Code: 404, Error: err.Error()})
				c.Abort()
				return
			}

			// check the existence of the nonce key.
			nk := a.GetNonceKey()
			if len(nk.Key) == 0 {
				c.JSON(401, ErrorMsg{Code: 401, Error: "unauthorized"})
				c.Abort()
				return
			}

			// check if the nonce key is expired.
			if a.IsExpired() {
				c.JSON(401, ErrorMsg{Code: 401, Error: "key is expired"})
				c.Abort()
				return
			}

			// decrypt the data.
			d, err := a.Decrypt(bytes.NewBuffer(r.Data))
			if err != nil {
				c.JSON(400, ErrorMsg{Code: 400, Error: err.Error()})
				c.Abort()
				return
			}

			// start with {" and end with }.
			match, _ := regexp.MatchString(`^{".*}$`, string(d))
			if !match {
				c.JSON(401, ErrorMsg{Code: 401, Error: "bad encrypt key"})
				c.Abort()
				return
			}

			c.Set("id", r.AccountID)
			c.Set("rawdata", d)

			// update the key expire time, and nonce value.
			t := time.Now().Add(svr.GetNonceKeyLifetime())
			nk.Expire_at = t
			nk.Nonce = cipher.RandByte(8)
			a.SetNonceKey(nk)
			return
		}
		c.JSON(400, ErrorMsg{Code: 400, Error: "bad request"})
		c.Abort()
	}
}
