package rpclient

import (
	"path/filepath"

	"github.com/skycoin/skycoin-exchange/src/rpclient/account"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
)

type Client interface {
	Run(addr string)
	GetServApiRoot() string
	GetServPubkey() cipher.PubKey
	GetLocalPubKey() cipher.PubKey
	GetLocalSecKey() cipher.SecKey
	GetAcntName() string
	CreateAccount() (*account.RpcAccount, error)
	HasAccount() bool
}

type RpcClient struct {
	*account.RpcAccount
	Cfg Config
}

type Config struct {
	ApiRoot    string
	DataDir    string
	AcntName   string
	ServPubkey cipher.PubKey
}

func New(cfg Config) Client {
	if cfg.DataDir == "" {
		homeDir := util.UserHome()
		cfg.DataDir = filepath.Join(homeDir, ".skycoin-exchange")
	}
	// init data dir.
	util.InitDataDir(cfg.DataDir)

	// init account dir.
	account.InitDir(filepath.Join(cfg.DataDir, "account/client"))

	cli := &RpcClient{
		Cfg: cfg,
	}

	// load account if exist.
	if account.IsExist(cfg.AcntName) {
		cli.RpcAccount = account.Load(cfg.AcntName)
	}

	return cli
}

func (rc *RpcClient) CreateAccount() (*account.RpcAccount, error) {
	// new account.
	rc.RpcAccount = account.New()

	// store the account
	if err := account.Store(rc.Cfg.AcntName, *rc.RpcAccount); err != nil {
		return nil, err
	}
	return rc.RpcAccount, nil
}

func (rc RpcClient) GetServApiRoot() string {
	return rc.Cfg.ApiRoot
}

func (rc RpcClient) GetServPubkey() cipher.PubKey {
	return rc.Cfg.ServPubkey
}

func (rc RpcClient) GetLocalPubKey() cipher.PubKey {
	return rc.Pubkey
}

func (rc RpcClient) GetLocalSecKey() cipher.SecKey {
	return rc.Seckey
}

func (rc RpcClient) HasAccount() bool {
	return rc.RpcAccount != nil
}

func (rc RpcClient) GetAcntName() string {
	return rc.Cfg.AcntName
}

func (rc *RpcClient) Run(addr string) {
	r := NewRouter(rc)
	r.Run(addr)
}
