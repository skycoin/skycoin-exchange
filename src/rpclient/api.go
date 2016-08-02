package rpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin/src/cipher"
)

// CreateAccount handle the request of creating account.
func CreateAccount(cli Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// generate account pubkey/privkey pair, pubkey is the account id.
		errRlt := &pp.EmptyRes{}
		for {
			p, s := cipher.GenerateKeyPair()
			r := pp.CreateAccountReq{
				Pubkey: pp.PtrString(p.Hex()),
			}

			req, _ := pp.MakeEncryptReq(&r, cli.GetServPubkey().Hex(), s.Hex())
			d, _ := json.Marshal(req)

			// send req to server.
			url := fmt.Sprintf("%s/accounts", cli.GetServApiRoot())
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(d))
			if err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			res := pp.EncryptRes{}
			json.NewDecoder(resp.Body).Decode(&res)
			defer resp.Body.Close()

			// handle the response
			if res.Result.GetSuccess() {
				v := pp.CreateAccountRes{}
				pp.DecryptRes(res, cli.GetServPubkey().Hex(), s.Hex(), &v)
				// store the account
				// account.Store(cli.GetAcntName(), *act)
				ret := struct {
					AccountID string `json:"account_id"`
					Key       string `json:"key"`
					CreatedAt int64  `json:"created_at"`
				}{
					AccountID: p.Hex(),
					Key:       s.Hex(),
					CreatedAt: v.GetCreatedAt(),
				}
				c.JSON(200, &ret)
				return
			} else {
				c.JSON(200, res)
				return
			}
		}
		c.JSON(200, *errRlt)
	}
}

func GetNewAddress(cli Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		errRlt := &pp.EmptyRes{}
		for {
			id, key, err := getAccountAndKey(c)
			if err != nil {
				errRlt = pp.MakeErrRes(err)
				break
			}

			cointype, exist := c.GetQuery("coin_type")
			if !exist {
				errRlt = pp.MakeErrRes(errors.New("coin type empty"))
				break
			}

			r := pp.GetDepositAddrReq{
				AccountId: &id,
				CoinType:  pp.PtrString(cointype),
			}

			req, _ := pp.MakeEncryptReq(&r, cli.GetServPubkey().Hex(), key)
			reqjson, _ := json.Marshal(req)

			// send req to server.
			url := fmt.Sprintf("%s/deposit_address", cli.GetServApiRoot())
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqjson))
			if err != nil {
				log.Println(err)
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			res := pp.EncryptRes{}
			json.NewDecoder(resp.Body).Decode(&res)
			defer resp.Body.Close()

			// handle the response
			if res.Result.GetSuccess() {
				v := pp.GetDepositAddrRes{}
				pp.DecryptRes(res, cli.GetServPubkey().Hex(), key, &v)
				c.JSON(200, v)
				return
			} else {
				c.JSON(200, res)
				return
			}
		}
		c.JSON(200, *errRlt)
	}
}

func GetBalance(cli Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		errRlt := &pp.EmptyRes{}
		for {
			id, key, err := getAccountAndKey(c)
			if err != nil {
				errRlt = pp.MakeErrRes(err)
				break
			}

			coinType, exist := c.GetQuery("coin_type")
			if !exist {
				errRlt = pp.MakeErrRes(errors.New("coin type empty"))
				break
			}

			gbr := pp.GetBalanceReq{
				AccountId: &id,
				CoinType:  pp.PtrString(coinType),
			}

			req, _ := pp.MakeEncryptReq(&gbr, cli.GetServPubkey().Hex(), key)
			js, _ := json.Marshal(req)

			url := fmt.Sprintf("%s/account/balance", cli.GetServApiRoot())
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(js))
			if err != nil {
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			res := pp.EncryptRes{}
			json.NewDecoder(resp.Body).Decode(&res)
			defer resp.Body.Close()

			// handle the response
			if res.Result.GetSuccess() {
				v := pp.GetBalanceRes{}
				pp.DecryptRes(res, cli.GetServPubkey().Hex(), key, &v)
				c.JSON(200, v)
				return
			} else {
				c.JSON(200, res)
				return
			}
		}
		c.JSON(200, *errRlt)
	}
}

func Withdraw(cli Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		rlt := &pp.EmptyRes{}
		for {
			id, key, err := getAccountAndKey(c)
			if err != nil {
				rlt = pp.MakeErrRes(err)
				break
			}

			cointype, exist := c.GetQuery("coin_type")
			if !exist {
				rlt = pp.MakeErrRes(errors.New("coin_type empty"))
				break
			}

			amount, exist := c.GetQuery("amount")
			if !exist {
				rlt = pp.MakeErrRes(errors.New("amount empty"))
				break
			}

			toAddr, exist := c.GetQuery("toaddr")
			if !exist {
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
			js, _ := json.Marshal(req)
			url := fmt.Sprintf("%s/account/withdrawal", cli.GetServApiRoot())
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(js))
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			res := pp.EncryptRes{}
			json.NewDecoder(resp.Body).Decode(&res)
			defer resp.Body.Close()

			// handle the response
			if res.Result.GetSuccess() {
				v := pp.WithdrawalRes{}
				pp.DecryptRes(res, cli.GetServPubkey().Hex(), key, &v)
				c.JSON(200, v)
				return
			} else {
				c.JSON(200, res)
				return
			}
		}
		c.JSON(200, rlt)
	}
}

func BidOrder(cli Client) gin.HandlerFunc {
	return orderHandler("bid", cli)
}

func AskOrder(cli Client) gin.HandlerFunc {
	return orderHandler("ask", cli)
}

func orderHandler(tp string, cli Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		rlt := &pp.EmptyRes{}
		for {
			rawReq := pp.OrderReq{}
			if err := c.BindJSON(&rawReq); err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_WrongRequest)
				break
			}
			id, key, err := getAccountAndKey(c)
			if err != nil {
				rlt = pp.MakeErrRes(err)
				break
			}

			rawReq.AccountId = &id
			req, _ := pp.MakeEncryptReq(&rawReq, cli.GetServPubkey().Hex(), key)
			js, _ := json.Marshal(req)
			url := fmt.Sprintf("%s/account/%s", cli.GetServApiRoot(), tp)
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(js))
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			res := pp.EncryptRes{}
			json.NewDecoder(resp.Body).Decode(&res)
			defer resp.Body.Close()

			// handle the response
			if res.Result.GetSuccess() {
				v := pp.OrderRes{}
				pp.DecryptRes(res, cli.GetServPubkey().Hex(), key, &v)
				c.JSON(200, v)
				return
			} else {
				c.JSON(200, res)
				return
			}
		}
		c.JSON(200, rlt)
	}
}

func GetOrderBook(cli Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		rlt := &pp.EmptyRes{}
		for {
			tp := c.Param("type")
			cp := c.Query("coin_pair")
			st := c.Query("start")
			ed := c.Query("end")
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
				Type:     &tp,
				Start:    &start,
				End:      &end,
			}
			jsn, _ := json.Marshal(req)
			url := fmt.Sprintf("%s/orders", cli.GetServApiRoot())
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsn))
			if err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			res := pp.GetOrderRes{}
			defer resp.Body.Close()
			if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			c.JSON(200, res)
			return
		}
		c.JSON(200, rlt)
	}
}

func getAccountAndKey(c *gin.Context) (id string, key string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("invalid id or key")
		}
	}()
	var exist bool
	id, exist = c.GetQuery("id")
	if !exist {
		return "", "", errors.New("id empty")
	}

	if _, err := cipher.PubKeyFromHex(id); err != nil {
		return "", "", errors.New("invalid id")
	}

	key, exist = c.GetQuery("key")
	if !exist {
		return "", "", errors.New("key empty")
	}

	if _, err := cipher.SecKeyFromHex(key); err != nil {
		return "", "", errors.New("invalid key")
	}

	return id, key, nil
}
