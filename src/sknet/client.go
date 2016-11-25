package sknet

import (
	"errors"
	"net"
)

const gPubkey = "02942e46684114b35fe15218dfdc6e0d74af0446a397b8fcbf8b46fb389f756eb8"

var gSeckey string

// Get send request to server, then read response and return.
func Get(addr string, path string, v interface{}) (*Response, error) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	r, err := MakeRequest(path, v)
	if err != nil {
		return nil, err
	}

	if err := Write(c, r); err != nil {
		return nil, err
	}
	// d, err := json.Marshal(r)
	// if err != nil {
	// 	return nil, err
	// }

	// buf := make([]byte, 4+len(d))
	// binary.BigEndian.PutUint32(buf[:], uint32(len(d)))
	// copy(buf[4:], d)
	// if err := binary.Write(c, binary.BigEndian, buf); err != nil {
	// 	return nil, err
	// }

	rsp := Response{}
	if err := Read(c, &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

// EncryGet will encrypt the request and decrypt the response.
func EncryGet(addr string, path string, req interface{}, res interface{}) error {
	if gSeckey == "" {
		return errors.New("private key is empty")
	}

	encReq, err := encrypt(req, gPubkey, gSeckey)
	if err != nil {
		return err
	}

	resp, err := Get(addr, path, encReq)
	if err != nil {
		return err
	}

	// decode the response.
	return decrypt(resp.Body, gPubkey, gSeckey, res)
}

// SetKey set local private key
func SetKey(key string) {
	gSeckey = key
}
