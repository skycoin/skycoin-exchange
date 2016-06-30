package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/skycoin/skycoin-exchange/src/server/account"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

type RespJson struct {
	Succress  bool   `json:"success"`
	AccountID string `json:"account_id"`
}

type AddressRespJson struct {
	RespJson
	CoinType string `json:"coin_type"`
	Address  string `json:"address"`
}

func CreateAccount(svr Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		pubkeyStr := c.PostForm("pubkey")
		pubkey, err := cipher.PubKeyFromHex(pubkeyStr)
		if err != nil {
			c.JSON(400, ErrorMsg{Code: 400, Error: "invalide pubkey"})
			return
		}

		_, err = svr.CreateAccountWithPubkey(pubkey)
		if err != nil {
			c.JSON(501, ErrorMsg{Code: 501, Error: "create account failed!"})
			return
		}
		c.JSON(201, RespJson{
			Succress:  true,
			AccountID: fmt.Sprintf("%x", pubkey),
		})
	}
}

// GetNewAddress create new address for specific account.
// POST param: id, coin_type
func GetNewAddress(svr Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.PostForm("id")
		coinType := c.PostForm("coin_type")

		fmt.Printf("account id:%s, len:%d\n", id, len(id))

		// convert to cipher.PubKey
		pubkey, err := cipher.PubKeyFromHex(id)
		if err != nil {
			c.JSON(400, ErrorMsg{Code: 400, Error: "error account id"})
			return
		}
		at, err := svr.GetAccount(account.AccountID(pubkey))
		if err != nil {
			c.JSON(404, ErrorMsg{Code: 404, Error: fmt.Sprintf("account id does not exist")})
			return
		}

		ct, err := wallet.ConvertCoinType(coinType)
		if err != nil {
			c.JSON(400, ErrorMsg{Code: 400, Error: err.Error()})
			return
		}

		addr := at.GetNewAddress(ct)
		c.JSON(201, AddressRespJson{
			RespJson: RespJson{
				Succress:  true,
				AccountID: id},
			CoinType: ct.String(),
			Address:  addr})
	}
}
