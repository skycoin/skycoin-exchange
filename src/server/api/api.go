package api

import (
	"encoding/json"

	"github.com/skycoin/skycoin-exchange/src/sknet"

	"gopkg.in/op/go-logging.v1"
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
