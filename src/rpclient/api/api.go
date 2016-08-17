package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"gopkg.in/op/go-logging.v1"

	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/sknet"
	"github.com/skycoin/skycoin/src/cipher"
)

var logger = logging.MustGetLogger("client.api")

type Servicer interface {
	GetServKey() cipher.PubKey
	GetServAddr() string
}

// CreateAccount handle the request of creating account.
func CreateAccount(se Servicer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// generate account pubkey/privkey pair, pubkey is the account id.
		errRlt := &pp.EmptyRes{}
		for {
			if r.Method != "POST" {
				logger.Error("require POST method")
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			p, s := cipher.GenerateKeyPair()
			r := pp.CreateAccountReq{
				Pubkey: pp.PtrString(p.Hex()),
			}

			req, err := makeEncryptReq(&r, se.GetServKey().Hex(), s.Hex())
			if err != nil {
				logger.Error("%s", err)
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			rsp, err := sknet.Get(se.GetServAddr(), "/auth/create/account", req)
			if err != nil {
				logger.Error("%s", err)
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			res, err := decodeRsp(rsp.Body, se.GetServKey().Hex(), s.Hex(), &pp.CreateAccountRes{})
			if err != nil {
				logger.Error("%s", err)
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			acntRes := res.(*pp.CreateAccountRes)
			if !acntRes.GetResult().GetSuccess() {
				sendJSON(w, res)
			} else {
				ret := struct {
					Result    pp.Result `json:"result"`
					AccountID string    `json:"account_id"`
					Key       string    `json:"key"`
					CreatedAt int64     `json:"created_at"`
				}{
					Result:    *acntRes.Result,
					AccountID: p.Hex(),
					Key:       s.Hex(),
					CreatedAt: acntRes.GetCreatedAt(),
				}
				sendJSON(w, &ret)
			}
			return
		}
		sendJSON(w, errRlt)
	}
}

func GetNewAddress(se Servicer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errRlt := &pp.EmptyRes{}
		for {
			if r.Method != "POST" {
				logger.Error("require POST method")
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			id, key, err := getAccountAndKey(r)
			if err != nil {
				logger.Error("%s", err)
				errRlt = pp.MakeErrRes(err)
				break
			}

			cointype := r.URL.Query().Get("coin_type")
			if cointype == "" {
				err := errors.New("coin type empty")
				logger.Error("%s", err)
				errRlt = pp.MakeErrRes(err)
				break
			}

			r := pp.GetDepositAddrReq{
				AccountId: &id,
				CoinType:  pp.PtrString(cointype),
			}

			req, err := makeEncryptReq(&r, se.GetServKey().Hex(), key)
			if err != nil {
				logger.Error("%s", err)
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			resp, err := sknet.Get(se.GetServAddr(), "/auth/create/deposit_address", req)
			if err != nil {
				logger.Error("%s", err)
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			res, err := decodeRsp(resp.Body, se.GetServKey().Hex(), key, &pp.GetDepositAddrRes{})
			if err != nil {
				logger.Error("%s", err)
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			sendJSON(w, res)
			return
		}
		sendJSON(w, errRlt)
	}
}

func GetBalance(se Servicer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errRlt := &pp.EmptyRes{}
		for {
			if r.Method != "GET" {
				logger.Error("require GET method")
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			id, key, err := getAccountAndKey(r)
			if err != nil {
				logger.Error("%s", err)
				errRlt = pp.MakeErrRes(err)
				break
			}
			coinType := r.URL.Query().Get("coin_type")
			if coinType == "" {
				err := errors.New("coin type empty")
				logger.Error("%s", err)
				errRlt = pp.MakeErrRes(err)
				break
			}

			gbr := pp.GetBalanceReq{
				AccountId: &id,
				CoinType:  pp.PtrString(coinType),
			}

			req, err := makeEncryptReq(&gbr, se.GetServKey().Hex(), key)
			if err != nil {
				logger.Error("%s", err)
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			resp, err := sknet.Get(se.GetServAddr(), "/auth/get/balance", req)
			if err != nil {
				logger.Error("%s", err)
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			res, err := decodeRsp(resp.Body, se.GetServKey().Hex(), key, &pp.GetBalanceRes{})
			if err != nil {
				logger.Error("%s", err)
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			sendJSON(w, res)
			return
		}
		sendJSON(w, errRlt)
	}
}

func Withdraw(se Servicer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rlt := &pp.EmptyRes{}
		for {
			if r.Method != "POST" {
				logger.Error("require POST method")
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			id, key, err := getAccountAndKey(r)
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrRes(err)
				break
			}

			cointype := r.URL.Query().Get("coin_type")
			if cointype == "" {
				err := errors.New("coin_type empty")
				logger.Error("%s", err)
				rlt = pp.MakeErrRes(err)
				break
			}

			amount := r.URL.Query().Get("amount")
			if amount == "" {
				rlt = pp.MakeErrRes(errors.New("amount empty"))
				break
			}

			toAddr := r.URL.Query().Get("toaddr")
			if toAddr == "" {
				rlt = pp.MakeErrRes(errors.New("toaddr empty"))
				break
			}

			amt, err := strconv.ParseUint(amount, 10, 64)
			if err != nil {
				rlt = pp.MakeErrRes(err)
				break
			}
			wr := pp.WithdrawalReq{
				AccountId:     &id,
				CoinType:      &cointype,
				Coins:         &amt,
				OutputAddress: &toAddr,
			}

			req, err := makeEncryptReq(&wr, se.GetServKey().Hex(), key)
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			resp, err := sknet.Get(se.GetServAddr(), "/auth/withdrawl", req)
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			res, err := decodeRsp(resp.Body, se.GetServKey().Hex(), key, &pp.WithdrawalRes{})
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			sendJSON(w, res)
			return
		}
		sendJSON(w, rlt)
	}
}

func GetCoins(se Servicer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rlt := &pp.EmptyRes{}
		for {
			if r.Method != "GET" {
				logger.Error("require GET method")
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			id, key, err := getAccountAndKey(r)
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrRes(err)
				break
			}
			rq := pp.GetCoinsReq{
				AccountId: pp.PtrString(id),
			}

			req, err := makeEncryptReq(&rq, se.GetServKey().Hex(), key)
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			rsp, err := sknet.Get(se.GetServAddr(), "/auth/get/coins", req)
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			res, err := decodeRsp(rsp.Body, se.GetServKey().Hex(), key, &pp.CoinsRes{})
			if err != nil {
				logger.Error("%s", err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			sendJSON(w, res)
			return
		}
		sendJSON(w, rlt)
	}
}

func getAccountAndKey(r *http.Request) (id string, key string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("invalid id or key")
		}
	}()
	id = r.URL.Query().Get("id")
	if id == "" {
		return "", "", errors.New("id empty")
	}

	if _, err := cipher.PubKeyFromHex(id); err != nil {
		return "", "", errors.New("invalid id")
	}

	key = r.URL.Query().Get("key")
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
	if res.Result.GetSuccess() {
		d, err := pp.Decrypt(res.Encryptdata, res.GetNonce(), pubkey, seckey)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(d, v); err != nil {
			return nil, err
		}
		return v, nil
	} else {
		return res, nil
	}
}
