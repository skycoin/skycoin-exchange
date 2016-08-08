package rpclient

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Response struct {
	Body *bytes.Buffer
}

func (res *Response) Read(r io.Reader) error {
	var len uint32
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return err
	}
	d := make([]byte, len)
	if err := binary.Read(r, binary.BigEndian, &d); err != nil {
		return err
	}
	res.Body = bytes.NewBuffer(d)
	return nil
}
