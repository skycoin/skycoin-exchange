package server

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/skycoin/skycoin/src/cipher"
)

// Authorize will decrypt the request, and encrypt the response.
func Authorize(svr Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		// read the ContentReqeust
		r := ContentRequest{}
		if c.BindJSON(&r) == nil {
			pubkey, err := cipher.PubKeyFromHex(r.AccountID)
			if err != nil {
				c.JSON(400, ErrorMsg{Code: 400, Error: "invalide account id"})
				c.Abort()
				return
			}

			data, err := svr.Decrypt(r.Data, pubkey, r.Nonce)
			if err != nil {
				c.JSON(401, ErrorMsg{Code: 401, Error: "unauthorized"})
				c.Abort()
				return
			}

			ok, err := regexp.MatchString(`^\{.*\}$`, string(data))
			if err != nil || !ok {
				c.JSON(401, ErrorMsg{Code: 401, Error: "unauthorized"})
				c.Abort()
				return
			}

			c.Set("rawdata", data)

			c.Next()

			// get response code
			code := c.MustGet("code").(int)
			rsp := c.MustGet("response")
			if code >= 200 && code < 300 {
				// encapsulate the response in ContentResponse.
				cr := MustToContentResponse(svr, pubkey, rsp)
				c.JSON(code, cr)
			}
			c.JSON(code, rsp)
			return
		}
		c.JSON(400, ErrorMsg{Code: 400, Error: "Bad request"})
	}
}

// MustToContentResponse marshal and encrypt the response object,
// return the ContentResponse object.
func MustToContentResponse(svr Server, pubkey cipher.PubKey, rsp interface{}) ContentResponse {
	d, err := json.Marshal(rsp)
	if err != nil {
		panic(err)
	}

	nonce := cipher.RandByte(8)
	resp, err := svr.Encrypt(d, pubkey, nonce)
	if err != nil {
		panic(err)
	}

	return ContentResponse{
		Success: true,
		Nonce:   fmt.Sprintf("%x", nonce),
		Data:    resp,
	}
}
