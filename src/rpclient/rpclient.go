package rpclient

import (
	"path/filepath"

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
	act, _ := LoadAccount("")
	return &RpcClient{
		ServApiRoot: apiRoot,
		ServPubkey:  pk,
		RA:          act,
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
