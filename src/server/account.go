package server

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skycoin/skycoin-exchange/src/server/account"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

func getRequest(c *gin.Context, out interface{}) error {
	d := c.MustGet("rawdata").([]byte)
	return json.Unmarshal(d, out)
}

// CreateAccount create account with specific pubkey,
func CreateAccount(svr Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := CreateAccountRequest{}
		if err := getRequest(c, &req); err != nil {
			Reply(c, 400, ErrorMsg{Code: 400, Error: err.Error()})
			return
		}

		pubkey, err := cipher.PubKeyFromHex(req.Pubkey)
		if err != nil {
			Reply(c, 400, ErrorMsg{Code: 400, Error: "invalide pubkey"})
			return
		}

		// create account with pubkey.
		_, err = svr.CreateAccountWithPubkey(pubkey)
		if err != nil {
			Reply(c, 501, ErrorMsg{Code: 501, Error: "create account failed!"})
			return
		}

		r := CreateAccountResponse{
			Success:   true,
			AccountID: req.Pubkey,
			CreatedAt: time.Now(),
		}
		Reply(c, 201, r)
		return
	}
}

// GetNewAddress account create new address for depositing.
func GetNewAddress(svr Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		dar := DepositAddressRequest{}
		err := getRequest(c, &dar)
		if err != nil {
			Reply(c, 400, ErrorMsg{Code: 400, Error: err.Error()})
			return
		}

		// convert to cipher.PubKey
		pubkey, err := cipher.PubKeyFromHex(dar.AccountID)
		if err != nil {
			Reply(c, 400, ErrorMsg{Code: 400, Error: "error account id"})
			return
		}

		at, err := svr.GetAccount(account.AccountID(pubkey))
		if err != nil {
			Reply(c, 404, ErrorMsg{Code: 404, Error: fmt.Sprintf("account id does not exist")})
			return
		}

		ct, err := wallet.ConvertCoinType(dar.CoinType)
		if err != nil {
			Reply(c, 400, ErrorMsg{Code: 400, Error: err.Error()})
			return
		}

		addr := svr.GetNewAddress(ct)
		at.AddDepositAddress(ct, addr)
		ds := DepositAddressResponse{
			Success:     true,
			AccountID:   dar.AccountID,
			DepositAddr: addr,
		}
		Reply(c, 201, ds)
	}
}

// Withdraw api handler for generating withdraw transaction.
func Withdraw(svr Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		wr := WithdrawRequest{}
		if err := getRequest(c, &wr); err != nil {
			Reply(c, 400, ErrorMsg{Code: 400, Error: err.Error()})
			return
		}

		// convert to cipher.PubKey
		pubkey, err := cipher.PubKeyFromHex(wr.AccountID)
		if err != nil {
			Reply(c, 400, ErrorMsg{Code: 400, Error: "error account id"})
			return
		}

		a, err := svr.GetAccount(account.AccountID(pubkey))
		if err != nil {
			Reply(c, 404, ErrorMsg{Code: 404, Error: "account id does not exist"})
			return
		}

		ct, err := wallet.ConvertCoinType(wr.CoinType)
		if err != nil {
			Reply(c, 400, ErrorMsg{Code: 400, Error: err.Error()})
			return
		}

		tx, err := GenerateWithdrawlTx(svr, a, ct, wr.Coins, wr.OutputAddress)
		if err != nil {
			Reply(c, 400, ErrorMsg{Code: 400, Error: err.Error()})
			return
		}

		resp := WithdrawResponse{
			AccountID: wr.AccountID,
			Tx:        tx,
		}
		Reply(c, 200, resp)
	}
}

// Reply set the code and response in gin, the gin Security middleware
// will encrypt the response, and send the encryped response to client.
func Reply(c *gin.Context, code int, r interface{}) {
	c.Set("code", code)
	c.Set("response", r)
}
