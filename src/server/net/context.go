package net

import (
	"encoding/json"
	"fmt"
)

type Context struct {
	Request  *Request               // Request from client
	Resp     ResponseWriter         // Response writer
	handlers []HandlerFunc          // request handlers, for records the middlewares.
	index    int                    // index points to the current request handler.
	Data     map[string]interface{} // data map, for transafer data between handlers.
}

// JSON write json response.
func (c *Context) JSON(data interface{}) error {
	return c.Resp.SendJSON(data)
}

// BindJSON marshal data from context.Request.
func (c *Context) BindJSON(v interface{}) error {
	return json.Unmarshal(c.Request.GetData(), v)
}

// Next execute the next handler.
func (c *Context) Next() {
	c.index += 1
	if c.index < len(c.handlers) {
		c.handlers[c.index](c)
	}
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
