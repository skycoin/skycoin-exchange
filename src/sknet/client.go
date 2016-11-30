package sknet

import (
	"errors"
	"net"

	"github.com/skycoin/skycoin/src/cipher"
)

var gPubkey = "02942e46684114b35fe15218dfdc6e0d74af0446a397b8fcbf8b46fb389f756eb8"

var gSeckey string

func init() {
	_, s := cipher.GenerateKeyPair()
	gSeckey = s.Hex()
}

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

// SetPubkey updates the server's pubkey
func SetPubkey(key string) {
	gPubkey = key
}
