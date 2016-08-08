package rpclient

import (
	"encoding/binary"
	"io"
	"log"
)

type Response struct {
	Body []byte
}

func (res *Response) Read(r io.Reader) error {
	var len uint32
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return err
	}
	log.Println("resp len:", len)
	res.Body = make([]byte, len)
	if err := binary.Read(r, binary.BigEndian, &res.Body); err != nil {
		return err
	}
	return nil
}
