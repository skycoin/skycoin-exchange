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

func (c *Context) Get(key string) (interface{}, error) {
	if v, ok := c.Data[key]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("%s does not exist", key)
}

func (c *Context) MustGet(key string) interface{} {
	v, err := c.Get(key)
	if err != nil {
		panic(err)
	}
	return v
}
