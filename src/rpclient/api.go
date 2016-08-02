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
	"github.com/skycoin/skycoin-exchange/src/rpclient/account"
)

// CreateAccount handle the request of creating account.
func CreateAccount(cli Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// generate account pubkey/privkey pair, pubkey is the account id.
		errRlt := &pp.EmptyRes{}
		for {
			act, err := cli.CreateAccount()
			if err != nil {
				log.Println(err)
				errRlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}
			r := pp.CreateAccountReq{
				Pubkey: act.Pubkey[:],
			}

			req, _ := pp.MakeEncryptReq(&r, cli.GetServPubkey().Hex(), act.Seckey.Hex())
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
				pp.DecryptRes(res, cli.GetServPubkey().Hex(), act.Seckey.Hex(), &v)
				// store the account
				account.Store(cli.GetAcntName(), *act)
				c.JSON(200, &v)
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
			if !cli.HasAccount() {
				errRlt = pp.MakeErrRes(errors.New("no account found"))
				break
			}

			cointype, exist := c.GetQuery("cointype")
			if !exist {
				errRlt = pp.MakeErrRes(errors.New("coin type empty"))
				break
			}

			pk := cli.GetLocalPubKey()
			r := pp.GetDepositAddrReq{
				AccountId: pk[:],
				CoinType:  pp.PtrString(cointype),
			}

			req, _ := pp.MakeEncryptReq(&r, cli.GetServPubkey().Hex(), cli.GetLocalSecKey().Hex())
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
				pp.DecryptRes(res, cli.GetServPubkey().Hex(), cli.GetLocalSecKey().Hex(), &v)
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
			if !cli.HasAccount() {
				errRlt = pp.MakeErrRes(errors.New("no account found"))
				break
			}

			coinType := c.Query("coin_type")
			pk := cli.GetLocalPubKey()
			gbr := pp.GetBalanceReq{
				AccountId: pk[:],
				CoinType:  pp.PtrString(coinType),
			}

			req, _ := pp.MakeEncryptReq(&gbr, cli.GetServPubkey().Hex(), cli.GetLocalSecKey().Hex())
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
				pp.DecryptRes(res, cli.GetServPubkey().Hex(), cli.GetLocalSecKey().Hex(), &v)
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
			cointype := c.Query("cointype")
			amount := c.Query("amount")
			toAddr := c.Query("toaddr")

			if cointype == "" || amount == "" || toAddr == "" {
				rlt = pp.MakeErrRes(errors.New(""))
				break
			}

			pk := cli.GetLocalPubKey()
			amtmp, err := strconv.Atoi(amount)
			if err != nil {
				rlt = pp.MakeErrRes(err)
				break
			}
			amt := uint64(amtmp)
			wr := pp.WithdrawalReq{
				AccountId:     pk[:],
				CoinType:      &cointype,
				Coins:         &amt,
				OutputAddress: &toAddr,
			}

			req, _ := pp.MakeEncryptReq(&wr, cli.GetServPubkey().Hex(), cli.GetLocalSecKey().Hex())
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
				pp.DecryptRes(res, cli.GetServPubkey().Hex(), cli.GetLocalSecKey().Hex(), &v)
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

			pk := cli.GetLocalPubKey()
			rawReq.AccountId = pk[:]

			req, _ := pp.MakeEncryptReq(&rawReq, cli.GetServPubkey().Hex(), cli.GetLocalSecKey().Hex())
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
				pp.DecryptRes(res, cli.GetServPubkey().Hex(), cli.GetLocalSecKey().Hex(), &v)
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
