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
			req pp.EncryptReq
			rlt *pp.EmptyRes
		)

		for {
			if c.BindJSON(&req) == nil {
				// validate pubkey.
				if err := validatePubkey(req.GetPubkey()); err != nil {
					logger.Error(err.Error())
					rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongPubkey)
					break
				}

				pubkey, err := cipher.PubKeyFromHex(req.GetPubkey())
				if err != nil {
					logger.Error(err.Error())
					rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongPubkey)
					break
				}

				key := cipher.ECDH(pubkey, ee.GetServPrivKey())
				data, err := cipher.Chacha20Decrypt(req.GetEncryptdata(), key, req.GetNonce())
				if err != nil {
					logger.Error(err.Error())
					rlt = pp.MakeErrResWithCode(pp.ErrCode_UnAuthorized)
					break
				}

				ok, err := regexp.MatchString(`^\{.*\}$`, string(data))
				if err != nil || !ok {
					logger.Error(err.Error())
					rlt = pp.MakeErrResWithCode(pp.ErrCode_UnAuthorized)
					break
				}

				c.Set("rawdata", data)

				c.Next()

				rsp, exist := c.Get("response")
				if exist {
					// encrypt the response.
					encData, nonce, err := pp.Encrypt(rsp, pubkey.Hex(), ee.GetServPrivKey().Hex())
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
			rlt = pp.MakeErrRes(errors.New("bad request"))
			break
		}
		c.JSON(rlt)
	}
}

// reply set the code and response in gin, the gin Security middleware
// will encrypt the response, and send the encryped response to client.
func reply(c *sknet.Context, r interface{}) {
	c.Set("response", r)
}
