package rpclient

import (
	"path/filepath"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
)

var (
	defaultAccountFile = filepath.Join(util.UserHome(), ".skycoin-exchange/account/client/act.data")
	actFileName        = "account.data"
)

type RpcAccount struct {
	Pubkey cipher.PubKey `json:"pubkey"`
	Seckey cipher.SecKey `json:"seckey"`
}

// LoadAccount load acccount info from local disk.
func LoadAccount(path string) (RpcAccount, error) {
	if path == "" {
		path = defaultAccountFile
	}
	a := RpcAccount{}
	err := util.LoadJSON(path, &a)
	return a, err
}

// StoreAccount store the rpcaccount to specific path.
func StoreAccount(a RpcAccount, path string) error {
	if path == "" {
		path = defaultAccountFile
	}

	return util.SaveJSON(path, a, 0777)
}
