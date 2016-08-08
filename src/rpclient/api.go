package rpclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin/src/cipher"
)

// CreateAccount handle the request of creating account.
func CreateAccount(cli Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// generate account pubkey/privkey pair, pubkey is the account id.
		errRlt := &pp.EmptyRes{}
		for {
			p, s := cipher.GenerateKeyPair()
			r := pp.CreateAccountReq{
				Pubkey: pp.PtrString(p.Hex()),
			}

			req, _ := pp.MakeEncryptReq(&r, cli.GetServPubkey().Hex(), s.Hex())
			rsp, err := sendRequest("/auth/create/account", req)
			if err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			res := pp.EncryptRes{}
			json.NewDecoder(rsp.Body).Decode(&res)

			// handle the response
			if res.Result.GetSuccess() {
				v := pp.CreateAccountRes{}
				pp.DecryptRes(res, cli.GetServPubkey().Hex(), s.Hex(), &v)
				ret := struct {
					Result    pp.Result `json:"result"`
					AccountID string    `json:"account_id"`
					Key       string    `json:"key"`
					CreatedAt int64     `json:"created_at"`
				}{
					Result:    *v.Result,
					AccountID: p.Hex(),
					Key:       s.Hex(),
					CreatedAt: v.GetCreatedAt(),
				}
				SendJSON(w, &ret)
				return
			} else {
				SendJSON(w, &res)
				return
			}
		}
		SendJSON(w, errRlt)
	}
}

func GetNewAddress(cli Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errRlt := &pp.EmptyRes{}
		for {
			id, key, err := getAccountAndKey(r)
			if err != nil {
				errRlt = pp.MakeErrRes(err)
				break
			}

			cointype := r.URL.Query().Get("coin_type")
			if cointype == "" {
				errRlt = pp.MakeErrRes(errors.New("coin type empty"))
				break
			}

			r := pp.GetDepositAddrReq{
				AccountId: &id,
				CoinType:  pp.PtrString(cointype),
			}

			req, _ := pp.MakeEncryptReq(&r, cli.GetServPubkey().Hex(), key)
			resp, err := sendRequest("/auth/create/deposit_address", req)
			if err != nil {
				log.Println(err)
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			res := pp.EncryptRes{}
			json.NewDecoder(resp.Body).Decode(&res)

			// handle the response
			if res.Result.GetSuccess() {
				v := pp.GetDepositAddrRes{}
				pp.DecryptRes(res, cli.GetServPubkey().Hex(), key, &v)
				SendJSON(w, &v)
				return
			} else {
				SendJSON(w, &res)
				return
			}
		}
		SendJSON(w, errRlt)
	}
}

func GetBalance(cli Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errRlt := &pp.EmptyRes{}
		for {
			id, key, err := getAccountAndKey(r)
			if err != nil {
				errRlt = pp.MakeErrRes(err)
				break
			}
			coinType := r.URL.Query().Get("coin_type")
			if coinType == "" {
				errRlt = pp.MakeErrRes(errors.New("coin type empty"))
				break
			}

			gbr := pp.GetBalanceReq{
				AccountId: &id,
				CoinType:  pp.PtrString(coinType),
			}

			req, _ := pp.MakeEncryptReq(&gbr, cli.GetServPubkey().Hex(), key)
			resp, err := sendRequest("/auth/get/balance", req)
			if err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			res := pp.EncryptRes{}
			json.NewDecoder(resp.Body).Decode(&res)

			// handle the response
			if res.Result.GetSuccess() {
				v := pp.GetBalanceRes{}
				pp.DecryptRes(res, cli.GetServPubkey().Hex(), key, &v)
				SendJSON(w, &v)
				return
			} else {
				SendJSON(w, &res)
				return
			}
		}
		SendJSON(w, errRlt)
	}
}

func Withdraw(cli Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rlt := &pp.EmptyRes{}
		for {
			id, key, err := getAccountAndKey(r)
			if err != nil {
				rlt = pp.MakeErrRes(err)
				break
			}

			cointype := r.URL.Query().Get("coin_type")
			if cointype == "" {
				rlt = pp.MakeErrRes(errors.New("coin_type empty"))
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

			req, _ := pp.MakeEncryptReq(&wr, cli.GetServPubkey().Hex(), key)
			resp, err := sendRequest("/auth/withdrawl", req)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			res := pp.EncryptRes{}
			json.NewDecoder(resp.Body).Decode(&res)

			// handle the response
			if res.Result.GetSuccess() {
				v := pp.WithdrawalRes{}
				pp.DecryptRes(res, cli.GetServPubkey().Hex(), key, &v)
				SendJSON(w, &v)
				return
			} else {
				SendJSON(w, &res)
				return
			}
		}
		SendJSON(w, rlt)
	}
}

func CreateBidOrder(cli Client) http.HandlerFunc {
	return createOrder(cli, "bid")
}

func CreateAskOrder(cli Client) http.HandlerFunc {
	return createOrder(cli, "ask")
}

func createOrder(cli Client, tp string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rlt := &pp.EmptyRes{}
		for {
			rawReq := pp.OrderReq{}
			if err := BindJSON(r, &rawReq); err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			id, key, err := getAccountAndKey(r)
			if err != nil {
				rlt = pp.MakeErrRes(err)
				break
			}

			rawReq.AccountId = &id
			req, _ := pp.MakeEncryptReq(&rawReq, cli.GetServPubkey().Hex(), key)
			resp, err := sendRequest(fmt.Sprintf("/auth/create/order/%s", rawReq.GetType()), req)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			res := pp.EncryptRes{}
			json.NewDecoder(resp.Body).Decode(&res)

			// handle the response
			if res.Result.GetSuccess() {
				v := pp.OrderRes{}
				pp.DecryptRes(res, cli.GetServPubkey().Hex(), key, &v)
				SendJSON(w, &v)
				return
			} else {
				SendJSON(w, &res)
				return
			}
		}
		SendJSON(w, rlt)
	}
}

func GetBidOrders(cli Client) http.HandlerFunc {
	return getOrders(cli, "bid")
}

func GetAskOrders(cli Client) http.HandlerFunc {
	return getOrders(cli, "ask")
}

func getOrders(cli Client, tp string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rlt := &pp.EmptyRes{}
		for {
			cp := r.URL.Query().Get("coin_pair")
			st := r.URL.Query().Get("start")
			ed := r.URL.Query().Get("end")
			if cp == "" || st == "" || ed == "" || tp == "" {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			start, err := strconv.ParseInt(st, 10, 64)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			end, err := strconv.ParseInt(ed, 10, 64)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}

			req := pp.GetOrderReq{
				CoinPair: &cp,
				Type:     pp.PtrString(tp),
				Start:    &start,
				End:      &end,
			}
			resp, err := sendRequest("/get/orders", req)
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			res := pp.GetOrderRes{}
			if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			SendJSON(w, &res)
			return
		}
		SendJSON(w, rlt)
	}
}

func GetCoins(cli Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rlt := &pp.EmptyRes{}
		for {
			rsp, err := sendRequest("/get/coins", nil)
			if err != nil {
				log.Println(err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			res := pp.CoinsRes{}
			if err := json.NewDecoder(rsp.Body).Decode(&res); err != nil {
				log.Println(err)
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			SendJSON(w, &res)
			return
		}
		SendJSON(w, rlt)
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
func SendJSON(w http.ResponseWriter, msg interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		panic(err)
	}
}

func BindJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func sendRequest(path string, data interface{}) (*Response, error) {
	c, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		return nil, err
	}
	defer c.Close()

	r, err := MakeRequest(path, data)
	if err != nil {
		return nil, err
	}
	return r.WriteRsp(c)
}
