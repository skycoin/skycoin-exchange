package net

import (
	"encoding/binary"
	"encoding/json"
	"io"

	"github.com/skycoin/skycoin-exchange/src/pp"
)

type Request struct {
	Raw        []byte // raw bytes
	pp.Request        // constructed request.
}

func (r *Request) Reset() {
	r.Raw = r.Raw[0:0]
	r.Request.Reset()
}

// Read serialise request from reader.
func (req *Request) Read(r io.Reader) error {
	var len uint32
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return err
	}
	req.Raw = make([]byte, len)
	if err := binary.Read(r, binary.BigEndian, &req.Raw); err != nil {
		return err
	}

	if err := json.Unmarshal(req.Raw[:], &req.Request); err != nil {
		logger.Error("%s", err)
		return err
	}
	return nil
}
