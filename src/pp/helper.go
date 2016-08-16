package pp

import (
	"encoding/json"
	"errors"

	"github.com/codahale/chacha20"
	"github.com/skycoin/skycoin/src/cipher"
)

func PtrBool(b bool) *bool {
	return &b
}

func PtrInt32(i int32) *int32 {
	return &i
}

func PtrUint32(i uint32) *uint32 {
	return &i
}

func PtrInt64(i int64) *int64 {
	return &i
}

func PtrUint64(i uint64) *uint64 {
	return &i
}

func PtrString(s string) *string {
	return &s
}

func MakeErrRes(err error) *EmptyRes {
	res := &EmptyRes{}
	res.Result = MakeResult(ErrCode_WrongFormat, err.Error())
	return res
}

func MakeErrResWithCode(code ErrCode) *EmptyRes {
	res := &EmptyRes{}
	res.Result = MakeResultWithCode(code)
	return res
}

func MakeResult(code ErrCode, reason string) *Result {
	result := &Result{}
	result.SetErrCode(code)
	result.SetReason(reason)
	return result
}

func MakeResultWithCode(code ErrCode) *Result {
	return MakeResult(code, code.String())
}

func (r *Result) SetErrCode(code ErrCode) {
	r.Errcode = PtrInt32(int32(code))
	r.Success = PtrBool(code == ErrCode_Success)
	r.Reason = PtrString(code.String())
}

func (r *Result) SetReason(reason string) {
	r.Reason = PtrString(reason)
}

func (r *Result) SetCodeAndReason(code ErrCode, reason string) {
	r.SetErrCode(code)
	r.SetReason(reason)
}

func MakeEncryptReq(r interface{}, pubkey string, seckey string) (req EncryptReq, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("encrypt faild")
		}
	}()
	sp := cipher.MustPubKeyFromHex(pubkey)
	cs := cipher.MustSecKeyFromHex(seckey)
	cp := cipher.PubKeyFromSecKey(cs).Hex()
	nonce := cipher.RandByte(chacha20.NonceSize)
	d, err := json.Marshal(r)
	if err != nil {
		return EncryptReq{}, err
	}
	key := cipher.ECDH(sp, cs)
	ed, err := cipher.Chacha20Encrypt([]byte(d), key, nonce)
	if err != nil {
		return EncryptReq{}, err
	}
	return EncryptReq{
		Pubkey:      &cp,
		Nonce:       nonce,
		Encryptdata: ed,
	}, nil
}

func DecryptRes(res EncryptRes, pubkey string, seckey string, v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("encrypt faild")
		}
	}()

	p := cipher.MustPubKeyFromHex(pubkey)
	s := cipher.MustSecKeyFromHex(seckey)
	key := cipher.ECDH(p, s)
	d, err := cipher.Chacha20Decrypt(res.Encryptdata, key, res.GetNonce())
	if err != nil {
		return err
	}

	// unmarshal the data
	return json.Unmarshal(d, v)
}

func BytesToPubKey(b []byte) cipher.PubKey {
	pk := cipher.PubKey{}
	copy(pk[:], b[:])
	return pk
}
