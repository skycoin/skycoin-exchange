package api

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/xchacha20"
	"github.com/skycoin/skycoin/src/cipher"
)

// Authorize will decrypt the request, and encrypt the response.
func Authorize(ee engine.Exchange) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			cnt_req pp.EncryptReq
			errRlt  = &pp.EmptyRes{}
		)

		for {
			if c.BindJSON(&cnt_req) == nil {
				cliPubkey := pp.BytesToPubKey(cnt_req.GetPubkey())
				if err := cliPubkey.Verify(); err != nil {
					errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongAccountId)
					break
				}

				data, err := xchacha20.Decrypt(cnt_req.GetEncryptdata(), cliPubkey, ee.GetServPrivKey(), cnt_req.GetNonce())
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
				if rsp, isExist := c.Get("response"); isExist {
					// encrypt the response.
					encryptData, nonce := mustEncryptRes(cliPubkey, ee.GetServPrivKey(), rsp)
					encpt_res := pp.EncryptRes{
						Result:      pp.MakeResultWithCode(pp.ErrCode_Success),
						Encryptdata: encryptData,
						Nonce:       nonce,
					}

					c.JSON(200, encpt_res)
				}
				return
			}
			errRlt = pp.MakeErrRes(errors.New("bad request"))
			break
		}
		c.JSON(200, *errRlt)
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

	nonce = cipher.RandByte(xchacha20.NonceSize)
	encryptData, err = xchacha20.Encrypt(d, pubkey, seckey, nonce)
	if err != nil {
		panic(err)
	}
	return
}

// func mustMakeEncryptRes(pubkey cipher.PubKey, seckey cipher.SecKey, rsp interface{}) *pp.EncryptRes {
// 	encryptData, nonce := mustEncryptRes(pubkey, seckey, rsp)
// 	encpt_res := pp.EncryptRes{
// 		Result:      pp.MakeResultWithCode(pp.ErrCode_Success),
// 		Encryptdata: encryptData,
// 		Nonce:       nonce,
// 	}
// 	return &encpt_res
// }

// reply set the code and response in gin, the gin Security middleware
// will encrypt the response, and send the encryped response to client.
func reply(c *gin.Context, r interface{}) {
	c.Set("response", r)
}
