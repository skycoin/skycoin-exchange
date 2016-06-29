package server

import (
	"encoding/hex"
	"fmt"
	"log"

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

type AccountRespJson struct {
	RespJson
	Seckey string `json:"seckey"`
}

func CreateAccount(svr *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		a, s, err := svr.AccountManager.CreateAccount()
		if err != nil {
			log.Println(err)
			c.JSON(501, ErrorMsg{Code: 501, Error: "Create Account Failed!"})
		}
		pubkey := a.GetAccountID()
		c.JSON(201, AccountRespJson{
			RespJson: RespJson{
				Succress:  true,
				AccountID: fmt.Sprintf("%x", pubkey),
			},
			Seckey: fmt.Sprintf("%x", s),
		})
	}
}

func GetNewAddress(svr *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.PostForm("id")
		coinType := c.PostForm("coin_type")

		fmt.Printf("account id:%s, len:%d\n", id, len(id))

		// convert to cipher.PubKey
		d, err := hex.DecodeString(id)
		if err != nil {
			c.JSON(400, ErrorMsg{Code: 400, Error: "error account id"})
			return
		}
		pk := cipher.PubKey{}
		copy(pk[:], d[:])

		at, err := svr.GetAccount(account.AccountID(pk))
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
