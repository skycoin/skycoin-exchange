package rpclient

import (
	"encoding/binary"
	"encoding/json"
	"io"

	"github.com/skycoin/skycoin-exchange/src/pp"
)

type Request struct {
	pp.Request
}

func MakeRequest(path string, data interface{}) (*Request, error) {
	r := &Request{}
	r.Path = pp.PtrString(path)
	if data == nil {
		return r, nil
	}
	d, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	r.Data = d[:]
	return r, nil
}

func (r *Request) Write(w io.Writer) error {
	d, err := json.Marshal(r)
	if err != nil {
		return err
	}

	buf := make([]byte, 4+len(d))
	binary.BigEndian.PutUint32(buf[:], uint32(len(d)))
	copy(buf[4:], d)
	if err := binary.Write(w, binary.BigEndian, buf); err != nil {
		return err
	}
	return nil
}

// WriteRsp send request to server, then read response and return.
func (r *Request) WriteRsp(rw io.ReadWriter) (*Response, error) {
	d, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 4+len(d))
	binary.BigEndian.PutUint32(buf[:], uint32(len(d)))
	copy(buf[4:], d)
	if err := binary.Write(rw, binary.BigEndian, buf); err != nil {
		return nil, err
	}

	rsp := &Response{}
	if err := rsp.Read(rw); err != nil {
		return nil, err
	}
	return rsp, nil
}
