package api

import (
	"errors"
	"regexp"

	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/server/engine"
	"github.com/skycoin/skycoin-exchange/src/sknet"
	"github.com/skycoin/skycoin/src/cipher"
)

// Authorize will decrypt the request, and encrypt the response.
func Authorize(ee engine.Exchange) sknet.HandlerFunc {
	return func(c *sknet.Context) {
		var (
			cnt_req pp.EncryptReq
			errRlt  = &pp.EmptyRes{}
		)

		for {
			if c.BindJSON(&cnt_req) == nil {
				cliPubkey, err := cipher.PubKeyFromHex(cnt_req.GetPubkey())
				if err != nil {
					logger.Error(err.Error())
					errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongAccountId)
					break
				}
				key := cipher.ECDH(cliPubkey, ee.GetServPrivKey())
				data, err := cipher.Chacha20Decrypt(cnt_req.GetEncryptdata(), key, cnt_req.GetNonce())
				if err != nil {
					logger.Error(err.Error())
					errRlt = pp.MakeErrResWithCode(pp.ErrCode_UnAuthorized)
					break
				}

				ok, err := regexp.MatchString(`^\{.*\}$`, string(data))
				if err != nil || !ok {
					logger.Error(err.Error())
					errRlt = pp.MakeErrResWithCode(pp.ErrCode_UnAuthorized)
					break
				}

				c.Set("rawdata", data)

				c.Next()

				rsp, exist := c.Get("response")
				if exist {
					// encrypt the response.
					encData, nonce, err := pp.Encrypt(rsp, cliPubkey.Hex(), ee.GetServPrivKey().Hex())
					if err != nil {
						panic(err)
					}

					// encryptData, nonce := mustEncryptRes(cliPubkey, ee.GetServPrivKey(), rsp)
					res := pp.EncryptRes{
						Result:      pp.MakeResultWithCode(pp.ErrCode_Success),
						Encryptdata: encData,
						Nonce:       nonce,
					}
					c.JSON(res)
				}
				return
			}
			errRlt = pp.MakeErrRes(errors.New("bad request"))
			break
		}
		c.JSON(errRlt)
	}
}

// mustEncryptRes marshal and encrypt the response object,
// return the encrypted data and the random nonce.
// func mustEncryptRes(pubkey cipher.PubKey, seckey cipher.SecKey, rsp interface{}) (encryptData []byte, nonce []byte) {
// 	var (
// 		d   []byte
// 		err error
// 	)
// 	d, err = json.Marshal(rsp)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	nonce = cipher.RandByte(chacha20.NonceSize)
// 	key := cipher.ECDH(pubkey, seckey)
// 	encryptData, err = cipher.Chacha20Encrypt(d, key, nonce)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return
// }

// reply set the code and response in gin, the gin Security middleware
// will encrypt the response, and send the encryped response to client.
func reply(c *sknet.Context, r interface{}) {
	c.Set("response", r)
}
