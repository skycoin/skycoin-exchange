package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skycoin/skycoin-exchange/src/server/account"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

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
		irawdata, _ := c.Get("rawdata")
		rawdata := irawdata.([]byte)

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
		d, err := json.Marshal(ds)
		if err != nil {
			panic(err)
		}

		resp, err := at.Encrypt(bytes.NewBuffer(d))
		if err != nil {
			panic(err)
		}

		c.JSON(201, ContentResponse{
			Success: true,
			Data:    resp,
		})
	}
}

func Encrypt(key []byte, data []byte) []byte {
	return data
}

func (dar DepositAddressRequest) MustToContentRequest(key []byte) ContentRequest {
	d, err := json.Marshal(dar)
	if err != nil {
		panic(err)
	}

	return ContentRequest{
		AccountID: dar.AccountID,
		Data:      d,
	}

}
