package api

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/skycoin/skycoin-exchange/src/chacha20"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin/src/cipher"
)

// Authorize will decrypt the request, and encrypt the response.
func Authorize(ee engine.Exchange) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			cnt_req   pp.ContentReq
			rsp       interface{}
			errRlt    *pp.EmptyRes
			accountid string
			cliPubkey cipher.PubKey
			err       error
		)

		for {
			if c.BindJSON(&cnt_req) == nil {
				accountid = cnt_req.GetAccountId()
				cliPubkey, err = cipher.PubKeyFromHex(cnt_req.GetAccountId())
				if err != nil {
					errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongAccountId)
					c.Abort()
					break
				}

				data, err := chacha20.Decrypt(cnt_req.GetEncryptdata(), cliPubkey, ee.GetServPrivKey(), cnt_req.GetNonce())
				if err != nil {
					errRlt = pp.MakeErrResWithCode(pp.ErrCode_UnAuthorized)
					c.Abort()
					break
				}

				ok, err := regexp.MatchString(`^\{.*\}$`, string(data))
				if err != nil || !ok {
					errRlt = pp.MakeErrResWithCode(pp.ErrCode_UnAuthorized)
					c.Abort()
					break
				}

				c.Set("rawdata", data)

				c.Next()

				// get response code
				code := c.MustGet("code").(int)
				rsp = c.MustGet("response")
				break
			}
			errRlt = pp.MakeErrRes(errors.New("bad request"))
			break
		}

		if errRlt != nil {
			c.JSON(200, *errRlt)
			return
		}

		// encrypt the response.
		encryptData, nonce := mustEncryptRes(cliPubkey, ee.GetServPrivKey(), rsp)
		cnt_res := pp.ContentRes{
			AccountId:   &accountid,
			Encryptdata: encryptData,
			Nonce:       nonce,
		}

		c.JSON(200, cnt_res)
	}
}

// mustEncryptRes marshal and encrypt the response object,
// return the encrypted data and the random nonce.
func mustEncryptRes(pubkey cipher.PubKey, seckey cipher.SecKey, rsp interface{}) (encryptData []byte, nonce []byte) {
	var (
		d   []byte
		err error
	)
	d, err = json.Marshal(rsp)
	if err != nil {
		panic(err)
	}

	nonce = cipher.RandByte(chacha20.NonceSize)
	encryptData, err = chacha20.Encrypt(d, pubkey, seckey, nonce)
	if err != nil {
		panic(err)
	}
	return
}
