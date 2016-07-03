package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skycoin/skycoin-exchange/src/server/account"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

// CreateAccount create account with specific pubkey,
func CreateAccount(svr Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := CreateAccountRequest{}
		if c.BindJSON(&req) == nil {
			pubkey, err := cipher.PubKeyFromHex(req.Pubkey)
			if err != nil {
				c.JSON(400, ErrorMsg{Code: 400, Error: "invalide pubkey"})
				return
			}

			// create account with pubkey.
			_, err = svr.CreateAccountWithPubkey(pubkey)
			if err != nil {
				c.JSON(501, ErrorMsg{Code: 501, Error: "create account failed!"})
				return
			}

			r := CreateAccountResponse{
				Succress:  true,
				AccountID: req.Pubkey,
				CreatedAt: time.Now(),
			}
			c.JSON(201, r)
			return
		}
		c.JSON(400, ErrorMsg{Code: 400, Error: "error request"})
	}
}

// GetNewAddress account create new address for depositing.
func GetNewAddress(svr Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawdata := c.MustGet("rawdata").([]byte)

		// unmarshal rawdata
		dar := DepositAddressRequest{}
		err := json.Unmarshal(rawdata, &dar)
		if err != nil {
			c.JSON(400, ErrorMsg{Code: 400, Error: "bad deposit address request"})
			return
		}

		// convert to cipher.PubKey
		pubkey, err := cipher.PubKeyFromHex(dar.AccountID)
		if err != nil {
			c.JSON(400, ErrorMsg{Code: 400, Error: "error account id"})
			return
		}
		at, err := svr.GetAccount(account.AccountID(pubkey))
		if err != nil {
			c.JSON(404, ErrorMsg{Code: 404, Error: fmt.Sprintf("account id does not exist")})
			return
		}

		ct, err := wallet.ConvertCoinType(dar.CoinType)
		if err != nil {
			c.JSON(400, ErrorMsg{Code: 400, Error: err.Error()})
			return
		}

		addr := at.GetNewAddress(ct)
		ds := DepositAddressResponse{
			AccountID:   dar.AccountID,
			DepositAddr: addr,
		}
		AuthReply(c, 201, ds)
	}
}

func Withdraw(svr Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		wr, err := newWithdrawRequest(c)
		if err != nil {
			c.JSON(400, ErrorMsg{Code: 400, Error: err.Error()})
			return
		}

		// convert to cipher.PubKey
		pubkey, err := cipher.PubKeyFromHex(wr.AccountID)
		if err != nil {
			c.JSON(400, ErrorMsg{Code: 400, Error: "error account id"})
			return
		}
		a, err := svr.GetAccount(account.AccountID(pubkey))
		if err != nil {
			c.JSON(404, ErrorMsg{Code: 404, Error: fmt.Sprintf("account id does not exist")})
			return
		}

		ct, err := wallet.ConvertCoinType(wr.CoinType)
		if err != nil {
			c.JSON(400, ErrorMsg{Code: 400, Error: err.Error()})
			return
		}
		tx, err := at.GenerateWithdrawTx(wr.Coins, ct)
		if err != nil {
			c.JSON(400, ErrorMsg{Code: 400, Error: err.Error()})
			return
		}

		resp := WithdrawResponse{
			AccountID: wr.AccountID,
			Tx:        tx,
		}
		AuthReply(c, 200, resp)
	}
}

func newWithdrawRequest(c *gin.Context) (WithdrawRequest, error) {
	rawdata := c.MustGet("rawdata").([]byte)
	// unmarshal rawdata
	wr := WithdrawRequest{}
	err := json.Unmarshal(rawdata, &wr)
	if err != nil {
		return WithdrawRequest{}, errors.New("bad withdraw request")
	}

	return wr, nil
}

func AuthReply(c *gin.Context, code int, r interface{}) {
	c.Set("code", code)
	c.Set("response", r)
}

func (wr WithdrawResponse) MustToContentResponse(a account.Accounter) ContentResponse {
	d, err := json.Marshal(wr)
	if err != nil {
		panic(err)
	}

	data, _ := a.Encrypt(bytes.NewBuffer(d))

	return ContentResponse{
		Success: true,
		Data:    data,
	}
}

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
