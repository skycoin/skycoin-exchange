package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skycoin/skycoin-exchange/src/server/account"
	"github.com/skycoin/skycoin/src/cipher"
)

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
// encrypt the response data, send it back to client.
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

			c.Set("id", r.AccountID)

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

// Security
func Security(svr Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := ContentRequest{}
		c.BindJSON(&r)

		// find the account.
		id := c.MustGet("id").(account.AccountID)
		a, _ := svr.GetAccount(id)

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

		c.Set("rawdata", d)

		c.Next()

		// get the response and encrypt the data.
		rsp := c.MustGet("response")
		code := c.MustGet("code")
		MustToContentResponse(a, c, code.(int), rsp)
	}
}

// MustToContentResponse marshal and encrypt the response object,
// generate the ContentResponse object.
func MustToContentResponse(a account.Accounter, c *gin.Context, code int, rsp interface{}) {
	d, err := json.Marshal(rsp)
	if err != nil {
		panic(err)
	}

	resp, err := a.Encrypt(bytes.NewBuffer(d))
	if err != nil {
		panic(err)
	}

	c.JSON(code, ContentResponse{
		Success: true,
		Nonce:   fmt.Sprintf("%x", a.GetNonceKey().Nonce),
		Data:    resp,
	})
}
