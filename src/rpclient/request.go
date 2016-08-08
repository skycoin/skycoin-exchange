package rpclient

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	"github.com/skycoin/skycoin-exchange/src/pp"
)

type Request struct {
	pp.Request
}

func MakeRequest(path string, data []byte) *Request {
	r := &Request{}
	r.Path = pp.PtrString(path)
	r.Data = data[:]
	return r
}

func (r *Request) Write(w io.Writer) error {
	d, err := json.Marshal(r)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(d))

	buf := make([]byte, 4+len(d))
	binary.BigEndian.PutUint32(buf[:], uint32(len(d)))
	copy(buf[4:], d)
	if err := binary.Write(w, binary.BigEndian, buf); err != nil {
		return err
	}
	return nil
}
