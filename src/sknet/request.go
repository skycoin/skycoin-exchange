package sknet

import "github.com/skycoin/skycoin-exchange/src/pp"

type Request struct {
	pp.Request // constructed request.
}

func (r *Request) Reset() {
	r.Request.Reset()
}

// Read serialise request from reader.
// func (req *Request) Read(r io.Reader) error {
// 	var len uint32
// 	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
// 		return err
// 	}
// 	d := make([]byte, len)
// 	if err := binary.Read(r, binary.BigEndian, &d); err != nil {
// 		return err
// 	}
//
// 	if err := json.Unmarshal(d, &req.Request); err != nil {
// 		logger.Error(err.Error())
// 		return err
// 	}
// 	return nil
// }
//
