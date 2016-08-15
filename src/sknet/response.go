package sknet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net"
)

// ResponseWriter interface for writing response.
type ResponseWriter interface {
	Write(p []byte) (n int, err error)
	SendJSON(data interface{}) error
}

// Response concrete response writer.
type Response struct {
	c    net.Conn
	Body *bytes.Buffer
}

// Write write data directly.
func (res *Response) Write(p []byte) (n int, err error) {
	return res.c.Write(p)
}

// SendJSON marshal the data into json, and then send.
func (res *Response) SendJSON(data interface{}) error {
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

// func (res *Response) Read(r io.Reader) error {
// 	var len uint32
// 	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
// 		return err
// 	}
// 	d := make([]byte, len)
// 	if err := binary.Read(r, binary.BigEndian, &d); err != nil {
// 		return err
// 	}
// 	res.Body = bytes.NewBuffer(d)
// 	return nil
// }
