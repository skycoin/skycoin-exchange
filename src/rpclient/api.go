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
		// pk := cli.GetLocalPubKey()
		// r := struct {
		// 	Rlt       *pp.EmptyRes
		// 	AccountId []byte
		// }{
		// 	Rlt:       rlt,
		// 	AccountId: pk[:],
		// }
		c.JSON(200, rlt)
	}
}

//
// func DecryptResponseBody(resp *http.Response, servPubkey cipher.PubKey, cliSeckey cipher.SecKey) ([]byte, error) {
// 	// decrypt the data.
// 	cnt := server.ContentRequest{}
// 	rd, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return []byte{}, err
// 	}
// 	fmt.Println(string(rd))
// 	defer resp.Body.Close()
// 	if err := json.Unmarshal(rd, &cnt); err != nil {
// 		return []byte{}, err
// 	}
//
// 	key := cipher.ECDH(servPubkey, cliSeckey)
// 	return ChaCha20(cnt.Data, key, cnt.Nonce)
// }
//
// func EncryptContentRequest(r interface{}, id string, key []byte) server.ContentRequest {
// 	d, err := json.Marshal(r)
// 	if err != nil {
// 		panic(err)
// 	}
// 	nonce := cipher.RandByte(xchacha20.NonceSize)
// 	data, err := ChaCha20(d, key, nonce)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return server.ContentRequest{
// 		AccountID: id,
// 		Data:      data,
// 		Nonce:     nonce,
// 	}
// }
//
// func GetErrorMsg(resp *http.Response) (server.ErrorMsg, error) {
// 	em := server.ErrorMsg{}
// 	rd, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return em, err
// 	}
//
// 	if err := json.Unmarshal(rd, &em); err != nil {
// 		return em, err
// 	}
//
// 	return em, nil
// }
//
// // decrypt the data
//
// func ChaCha20(data []byte, key []byte, nonce []byte) ([]byte, error) {
// 	e := make([]byte, len(data))
// 	c, err := chacha20.New(key, nonce)
// 	if err != nil {
// 		return []byte{}, err
// 	}
// 	c.XORKeyStream(e, data)
// 	return e, nil
// }
//
// func makeErrorMsg(code int, reason string) server.ErrorMsg {
// 	return server.ErrorMsg{Code: code, Error: reason}
// }
