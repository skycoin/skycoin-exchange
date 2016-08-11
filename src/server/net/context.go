package net

import (
	"encoding/json"
	"fmt"
)

type Context struct {
	Request  *Request
	Resp     ResponseWriter
	handlers []HandlerFunc
	index    int
	Data     map[string]interface{}
}

func (c *Context) JSON(data interface{}) error {
	return c.Resp.SendJSON(data)
}

func (c *Context) BindJSON(v interface{}) error {
	return json.Unmarshal(c.Request.GetData(), v)
}

func (c *Context) Next() {
	c.index += 1
	if c.index < len(c.handlers) {
		c.handlers[c.index](c)
	}
}

func (c *Context) Set(key string, v interface{}) {
	c.Data[key] = v
}

func (c *Context) Get(key string) (interface{}, bool) {
	if v, ok := c.Data[key]; ok {
		return v, true
	}

	return nil, false
}

func (c *Context) MustGet(key string) interface{} {
	v, exist := c.Get(key)
	if !exist {
		panic(fmt.Sprintf("%s does not exist", key))
	}
	return v
}
