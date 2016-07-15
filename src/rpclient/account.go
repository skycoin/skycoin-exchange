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
	// check whether the account dir is exist.
	// if _, err := os.Stat(path); os.IsNotExist(err) {
	// 	// create the dir.
	// 	if err := os.MkdirAll(path, 0777); err != nil {
	// 		return err
	// 	}
	// }
	// // create the file
	// f, err := os.Create(filepath.Join(path, actFileName))
	// if err != nil {
	// 	return err
	// }
	// defer f.Close()
	// d, err := json.MarshalIndent(a, "", " ")
	// if err != nil {
	// 	return err
	// }
	// if _, err := f.WriteString(string(d)); err != nil {
	// 	return err
	// }
	// return nil
}
