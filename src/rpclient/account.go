package rpclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
)

var (
	defaultAccountDir = filepath.Join(util.UserHome(), ".skycoin-exchange/account/client")
	actFileName       = "account.data"
)

type RpcAccount struct {
	Pubkey cipher.PubKey `json:"pubkey"`
	Seckey cipher.SecKey `json:"seckey"`
}

// LoadAccount load acccount info from local disk.
func LoadAccount(path string) (RpcAccount, error) {
	if path == "" {
		path = defaultAccountDir
	}
	// check whether the account file is exist.
	actFile := filepath.Join(path, actFileName)
	if _, err := os.Stat(actFile); os.IsExist(err) {
		d, err := ioutil.ReadFile(actFile)
		if err != nil {
			return RpcAccount{}, err
		}
		a := RpcAccount{}
		err = json.Unmarshal(d, &a)
		if err != nil {
			return RpcAccount{}, err
		}
		return a, nil
	}
	return RpcAccount{}, fmt.Errorf("%s not exist", actFile)
}

func StoreAccount(a RpcAccount, path string) error {
	if path == "" {
		path = defaultAccountDir
	}

	// check whether the account dir is exist.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// create the dir.
		if err := os.MkdirAll(path, 0777); err != nil {
			return err
		}
	}
	// create the file
	f, err := os.Create(filepath.Join(path, actFileName))
	if err != nil {
		return err
	}
	defer f.Close()
	d, err := json.MarshalIndent(a, "", " ")
	if err != nil {
		return err
	}
	if _, err := f.WriteString(string(d)); err != nil {
		return err
	}
	return nil
}
