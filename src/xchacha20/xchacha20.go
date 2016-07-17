package xchacha20

import (
	"errors"

	"github.com/codahale/chacha20"
	"github.com/skycoin/skycoin/src/cipher"
)

const (
	NonceSize = 8
)

func Decrypt(data []byte, pubkey cipher.PubKey, seckey cipher.SecKey, nonce []byte) (d []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("server decrypt faild")
		}
	}()

	key := cipher.ECDH(pubkey, seckey)
	d, err = encOrDec(data, key, nonce)
	return
}

func Encrypt(data []byte, pubkey cipher.PubKey, seckey cipher.SecKey, nonce []byte) (d []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("server encrypt faild")
		}
	}()

	key := cipher.ECDH(pubkey, seckey)
	d, err = encOrDec(data, key, nonce)
	return
}

// chacha20
func encOrDec(data []byte, key []byte, nonce []byte) ([]byte, error) {
	e := make([]byte, len(data))
	c, err := chacha20.New(key, nonce)
	if err != nil {
		return []byte{}, err
	}
	c.XORKeyStream(e, data)
	return e, nil
}
