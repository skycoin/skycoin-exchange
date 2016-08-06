package api

import (
	"encoding/json"

	"gopkg.in/op/go-logging.v1"

	"github.com/gin-gonic/gin"
)

var logger = logging.MustGetLogger("exchange.api")

func getRequest(c *gin.Context, out interface{}) error {
	d := c.MustGet("rawdata").([]byte)
	return json.Unmarshal(d, out)
}

type ReqParams struct {
	Values map[string]interface{}
}

func NewReqParams() *ReqParams {
	return &ReqParams{
		Values: make(map[string]interface{}),
	}
}
