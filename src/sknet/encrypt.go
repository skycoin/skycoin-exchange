package sknet

import (
	"encoding/json"
	"io"

	"fmt"

	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin/src/cipher"
)

func encrypt(r interface{}, pubkey string, seckey string) (*pp.EncryptReq, error) {
	encData, nonce, err := pp.Encrypt(r, pubkey, seckey)
	if err != nil {
		return nil, err
	}

	s, err := cipher.SecKeyFromHex(seckey)
	if err != nil {
		return nil, err
	}

	p := cipher.PubKeyFromSecKey(s)
	return &pp.EncryptReq{
		Pubkey:      pp.PtrString(p.Hex()),
		Nonce:       nonce,
		Encryptdata: encData,
	}, nil
}

func decrypt(r io.Reader, pubkey string, seckey string, v interface{}) error {
	res := pp.EncryptRes{}
	if err := json.NewDecoder(r).Decode(&res); err != nil {
		return err
	}

	// handle the response
	if !res.Result.GetSuccess() {
		return fmt.Errorf("%v", res.Result.GetReason())
	}
	d, err := pp.Decrypt(res.Encryptdata, res.GetNonce(), pubkey, seckey)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(d, v); err != nil {
		return err
	}
	return nil
}
