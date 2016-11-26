package sknet

import (
	"io"
	"net"
)

// Response concrete response writer.
type Response struct {
	c    net.Conn
	Body io.Reader
}

// Write write data directly.
func (res *Response) Write(p []byte) (n int, err error) {
	return res.c.Write(p)
}

// SendJSON marshal the data into json, and then send.
func (res *Response) SendJSON(data interface{}) error {
	return Write(res.c, data)
}
