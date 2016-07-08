package server

import (
	"errors"

	"github.com/codahale/chacha20"
	"github.com/skycoin/skycoin/src/cipher"
)

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

func Decrypt(data []byte, serverPrivKey cipher.SecKey, clientPubkey cipher.PubKey, nonce []byte) (d []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("server decrypt faild")
		}
	}()

	key := cipher.ECDH(clientPubkey, serverPrivKey)
	d, err = encOrDec(data, key, nonce)
	return
}

func Encrypt(data []byte, serverPrivKey cipher.SecKey, clientPubkey cipher.PubKey, nonce []byte) (d []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("server encrypt faild")
		}
	}()

	key := cipher.ECDH(clientPubkey, serverPrivKey)
	d, err = encOrDec(data, key, nonce)
	return
}
