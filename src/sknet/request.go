package sknet

import (
	"encoding/json"

	"github.com/skycoin/skycoin-exchange/src/pp"
)

type Request struct {
	pp.Request // constructed request.
}

func (r *Request) Reset() {
	r.Request.Reset()
}

// MakeRequest creates request
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
