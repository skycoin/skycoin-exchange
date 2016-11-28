package sknet

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin/src/cipher"
)

// Authorize will decrypt the request, and it's a buildin middleware for skynet.
func Authorize(servSeckey string) HandlerFunc {
	return func(c *Context) error {
		var (
			req pp.EncryptReq
			rlt *pp.EmptyRes
		)

		c.ServSeckey = servSeckey

		for {
			if c.UnmarshalReq(&req) == nil {
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
				c.Pubkey = pubkey.Hex()

				seckey, err := cipher.SecKeyFromHex(servSeckey)
				if err != nil {
					logger.Error(err.Error())
					rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
					break
				}

				key := cipher.ECDH(pubkey, seckey)
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

				c.Raw = data

				return c.Next()
			}
			rlt = pp.MakeErrRes(errors.New("bad request"))
			break
		}
		return c.Error(rlt)
	}
}

// reply set the code and response in gin, the gin Security middleware
// will encrypt the response, and send the encryped response to client.
func reply(c *Context, r interface{}) {
	c.Set("response", r)
}

func validatePubkey(key string) (err error) {
	defer func() {
		// the PubKeyFromHex may panic if the key is invalidate.
		// use recover to catch the panic, and return false.
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	_, err = cipher.PubKeyFromHex(key)
	return
}
