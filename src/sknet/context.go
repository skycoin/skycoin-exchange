package sknet

import (
	"encoding/json"
	"fmt"

	"github.com/skycoin/skycoin-exchange/src/pp"
)

// ResponseWriter interface for writing response.
type ResponseWriter interface {
	Write(p []byte) (n int, err error)
	SendJSON(data interface{}) error
}

// Context is the most important part of sknet. It allows us to pass variables between middleware,
// manage the flow, validate the request, decrypt and encrypt request and response.
type Context struct {
	Request    *Request               // Request from client
	Raw        []byte                 // the decrypted raw data
	Pubkey     string                 // client pubkey
	ServSeckey string                 // server seckey
	Resp       ResponseWriter         // Response writer
	handlers   []HandlerFunc          // request handlers, for records the middlewares.
	index      int                    // index points to the current request handler.
	Data       map[string]interface{} // data map, for transafer data between handlers.
}

// JSON encrypt the data and write response.
func (c *Context) SendJSON(data interface{}) error {
	encData, nonce, err := pp.Encrypt(data, c.Pubkey, c.ServSeckey)
	if err != nil {
		return err
	}

	res := pp.EncryptRes{
		Result:      pp.MakeResultWithCode(pp.ErrCode_Success),
		Encryptdata: encData,
		Nonce:       nonce,
	}
	return c.Resp.SendJSON(res)
}

// ErrorJSON write json response.
func (c *Context) Error(data interface{}) error {
	return c.Resp.SendJSON(data)
}

func (c *Context) UnmarshalReq(v interface{}) error {
	return json.Unmarshal(c.Request.GetData(), v)
}

// BindJSON unmarshal data from context.Request.
func (c *Context) BindJSON(v interface{}) error {
	return json.Unmarshal(c.Raw, v)
}

// Next execute the next handler.
func (c *Context) Next() error {
	c.index++
	if c.index < len(c.handlers) {
		return c.handlers[c.index](c)
	}
	return nil
}

// Set write data of key into context.Data.
func (c *Context) Set(key string, v interface{}) {
	c.Data[key] = v
}

// Get read data of key from context.Data.
func (c *Context) Get(key string) (interface{}, bool) {
	if v, ok := c.Data[key]; ok {
		return v, true
	}

	return nil, false
}

// MustGet read data of key, panic if not exist.
func (c *Context) MustGet(key string) interface{} {
	v, exist := c.Get(key)
	if !exist {
		panic(fmt.Sprintf("%s does not exist", key))
	}
	return v
}

// Reset set all value to initial state.
func (c *Context) Reset() {
	c.Data = make(map[string]interface{})
	c.index = 0
	c.handlers = c.handlers[:0]
	c.Request = nil
	c.Resp = nil
}
