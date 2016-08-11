package net

import (
	"encoding/binary"
	"encoding/json"
	"net"
)

// ResponseWriter interface for writing response.
type ResponseWriter interface {
	Write(p []byte) (n int, err error)
	SendJSON(data interface{}) error
}

// NetResponse concrete response writer.
type NetResponse struct {
	c net.Conn
}

// Write write data directly.
func (res *NetResponse) Write(p []byte) (n int, err error) {
	return res.c.Write(p)
}

// SendJSON marshal the data into json, and then send.
func (res *NetResponse) SendJSON(data interface{}) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	buf := make([]byte, 4+len(d))
	binary.BigEndian.PutUint32(buf[:], uint32(len(d)))
	copy(buf[4:], d)
	if err := binary.Write(res.c, binary.BigEndian, buf); err != nil {
		return err
	}
	return nil
}
