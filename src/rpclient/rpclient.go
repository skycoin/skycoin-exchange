package rpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/codahale/chacha20"
	"github.com/gin-gonic/gin"
	"github.com/skycoin/skycoin-exchange/src/server"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
)

func init() {
	homeDir := util.UserHome()
	util.InitDataDir(filepath.Join(homeDir, ".skycoin-exchange"))
}

type Client interface {
	Run(addr string)
	GetServApiRoot() string
	GetServPubkey() cipher.PubKey
	GetLocalPubKey() cipher.PubKey
	GetLocalSecKey() cipher.SecKey
	CreateAccount() *RpcAccount
}

type RpcClient struct {
	RA          RpcAccount
	ServApiRoot string
	ServPubkey  cipher.PubKey
}

func New(apiRoot string, servPubkey string) Client {
	pk := cipher.MustPubKeyFromHex(servPubkey)
	return &RpcClient{
		ServApiRoot: apiRoot,
		ServPubkey:  pk,
	}
}

func (rc *RpcClient) CreateAccount() *RpcAccount {
	p, s := cipher.GenerateKeyPair()
	rc.RA = RpcAccount{
		Pubkey: p,
		Seckey: s,
	}
	return &rc.RA
}

func (rc RpcClient) GetServApiRoot() string {
	return rc.ServApiRoot
}

func (rc RpcClient) GetServPubkey() cipher.PubKey {
	return rc.ServPubkey
}

func (rc RpcClient) GetLocalPubKey() cipher.PubKey {
	return rc.RA.Pubkey
}

func (rc RpcClient) GetLocalSecKey() cipher.SecKey {
	return rc.RA.Seckey
}

func (rc *RpcClient) Run(addr string) {
	r := NewRouter(rc)
	r.Run(addr)
}

func CreateAccount(cli Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// generate account pubkey/privkey pair, pubkey is the account id.
		act := cli.CreateAccount()
		r := server.CreateAccountRequest{
			Pubkey: act.Pubkey.Hex(),
		}

		key := cipher.ECDH(cli.GetServPubkey(), act.Seckey)
		req := server.MustToContentRequest(r, act.Pubkey.Hex(), key).MustJson()

		// send req to server.
		url := fmt.Sprintf("%s/accounts", cli.GetServApiRoot())
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(req))
		if err != nil {
			panic(err)
		}

		// handle the response
		if resp.StatusCode == 200 || resp.StatusCode == 201 {
			rawdata, err := DecryptResponseBody(resp, cli.GetServPubkey(), act.Seckey)
			if err != nil {
				panic(err)
			}
			car := server.CreateAccountResponse{}
			if err := json.Unmarshal(rawdata, &car); err != nil {
				panic(err)
			}
			c.JSON(resp.StatusCode, car)
			return
		}

		em, err := GetErrorMsg(resp)
		if err != nil {
			panic(err)
		}
		c.JSON(resp.StatusCode, em)
	}
}

func GetNewAddress(cli Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		cointype, exist := c.GetQuery("cointype")
		if !exist {
			c.JSON(400, server.ErrorMsg{Code: 400, Error: "cointype empty"})
			return
		}

		r := server.DepositAddressRequest{
			AccountID: cli.GetLocalPubKey().Hex(),
			CoinType:  cointype,
		}

		key := cipher.ECDH(cli.GetServPubkey(), cli.GetLocalSecKey())
		req := EncryptContentRequest(r, cli.GetLocalPubKey().Hex(), key).MustJson()

		// send req to server.
		url := fmt.Sprintf("%s/deposit_address", cli.GetServApiRoot())
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(req))
		if err != nil {
			panic(err)
		}

		// handle the response
		if resp.StatusCode == 200 || resp.StatusCode == 201 {
			rawdata, err := DecryptResponseBody(resp, cli.GetServPubkey(), cli.GetLocalSecKey())
			if err != nil {
				panic(err)
			}
			car := server.DepositAddressResponse{}
			if err := json.Unmarshal(rawdata, &car); err != nil {
				panic(err)
			}
			c.JSON(resp.StatusCode, car)
			return
		}

		em, err := GetErrorMsg(resp)
		if err != nil {
			panic(err)
		}
		c.JSON(resp.StatusCode, em)
	}
}

func Withdraw(cli Client) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func DecryptResponseBody(resp *http.Response, servPubkey cipher.PubKey, cliSeckey cipher.SecKey) ([]byte, error) {
	// decrypt the data.
	cnt := server.ContentRequest{}
	rd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	fmt.Println(string(rd))
	defer resp.Body.Close()
	if err := json.Unmarshal(rd, &cnt); err != nil {
		return []byte{}, err
	}

	key := cipher.ECDH(servPubkey, cliSeckey)
	return ChaCha20(cnt.Data, key, cnt.Nonce)
}

func EncryptContentRequest(r interface{}, id string, key []byte) server.ContentRequest {
	d, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	nonce := cipher.RandByte(8)
	data, err := ChaCha20(d, key, nonce)
	if err != nil {
		panic(err)
	}
	return server.ContentRequest{
		AccountID: id,
		Data:      data,
		Nonce:     nonce,
	}
}

func GetErrorMsg(resp *http.Response) (server.ErrorMsg, error) {
	em := server.ErrorMsg{}
	rd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return em, err
	}

	if err := json.Unmarshal(rd, &em); err != nil {
		return em, err
	}

	return em, nil
}

// decrypt the data

func ChaCha20(data []byte, key []byte, nonce []byte) ([]byte, error) {
	e := make([]byte, len(data))
	c, err := chacha20.New(key, nonce)
	if err != nil {
		return []byte{}, err
	}
	c.XORKeyStream(e, data)
	return e, nil
}
