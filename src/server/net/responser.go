package net

import (
	"encoding/binary"
	"encoding/json"
	"net"
)

type ResponseWriter interface {
	Write(p []byte) (n int, err error)
	SendJSON(data interface{})
}

type NetResponse struct {
	c net.Conn
}

func (res *NetResponse) Write(p []byte) (n int, err error) {
	return res.c.Write(p)
}

func (res *NetResponse) SendJSON(data interface{}) {
	d, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 4+len(d))
	binary.BigEndian.PutUint32(buf[:], uint32(len(d)))

	buf = append(buf, d...)
	if err := binary.Write(res.c, binary.BigEndian, buf); err != nil {
		panic(err)
	}
}
