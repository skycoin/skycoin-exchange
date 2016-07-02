package server

import "github.com/codahale/chacha20"

// Encrypt chacha20
func Encrypt(data []byte, key []byte, nonce []byte) ([]byte, error) {
	return encOrDec(data, key, nonce)
}

func Decrypt(data []byte, key []byte, nonce []byte) ([]byte, error) {
	return encOrDec(data, key, nonce)
}

func encOrDec(data []byte, key []byte, nonce []byte) ([]byte, error) {
	e := make([]byte, len(data))
	c, err := chacha20.New(key, nonce)
	if err != nil {
		return []byte{}, err
	}
	c.XORKeyStream(e, data)
	return e, nil
}
