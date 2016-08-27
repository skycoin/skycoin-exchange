package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"gopkg.in/op/go-logging.v1"

	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin/src/cipher"
)

var logger = logging.MustGetLogger("client.api")

// Servicer api service interface
type Servicer interface {
	GetServKey() cipher.PubKey
	GetServAddr() string
}

func getPubkey(r *http.Request) (pubkey string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("invalid pubkey")
		}
	}()

	pubkey = r.FormValue("pubkey")
	if pubkey == "" {
		return "", errors.New("pubkey empty")
	}

	if _, err = cipher.PubKeyFromHex(pubkey); err != nil {
		return "", errors.New("invalid pubkey")
	}
	return
}

func getAccountAndKey(r *http.Request) (id string, key string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("invalid id or key")
		}
	}()
	id = r.FormValue("id")
	if id == "" {
		return "", "", errors.New("id empty")
	}

	if _, err := cipher.PubKeyFromHex(id); err != nil {
		return "", "", errors.New("invalid id")
	}

	key = r.FormValue("key")
	if key == "" {
		return "", "", errors.New("key empty")
	}

	if _, err := cipher.SecKeyFromHex(key); err != nil {
		return "", "", errors.New("invalid key")
	}

	return id, key, nil
}

// JSON to an http response
func sendJSON(w http.ResponseWriter, msg interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		panic(err)
	}
}

func bindJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func makeEncryptReq(r interface{}, pubkey string, seckey string) (*pp.EncryptReq, error) {
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

func decodeRsp(r io.Reader, pubkey string, seckey string, v interface{}) (interface{}, error) {
	res := pp.EncryptRes{}
	if err := json.NewDecoder(r).Decode(&res); err != nil {
		return nil, err
	}

	// handle the response
	if !res.Result.GetSuccess() {
		return res, nil
	}
	d, err := pp.Decrypt(res.Encryptdata, res.GetNonce(), pubkey, seckey)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(d, v); err != nil {
		return nil, err
	}
	return v, nil
}
