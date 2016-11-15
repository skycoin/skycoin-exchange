package api

import (
	"encoding/json"
	"fmt"

	logging "github.com/op/go-logging"
	"github.com/skycoin/skycoin-exchange/src/sknet"
	"github.com/skycoin/skycoin/src/cipher"
)

var logger = logging.MustGetLogger("exchange.api")

func getRequest(c *sknet.Context, out interface{}) error {
	d := c.MustGet("rawdata").([]byte)
	return json.Unmarshal(d, out)
}

// ReqParams records the request params
type ReqParams struct {
	Values map[string]interface{}
}

// NewReqParams make and init the ReqParams.
func NewReqParams() *ReqParams {
	return &ReqParams{
		Values: make(map[string]interface{}),
	}
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
