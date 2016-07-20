package api

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

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
