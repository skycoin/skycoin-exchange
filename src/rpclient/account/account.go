package account

import (
	"log"
	"os"
	"path/filepath"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
)

var (
	acntDir = filepath.Join(util.UserHome(), ".skycoin-exchange/account/client")
)

type RpcAccount struct {
	Pubkey cipher.PubKey
	Seckey cipher.SecKey
}

type AccountJson struct {
	Pubkey []byte `json:"pubkey"`
	Seckey []byte `json:"seckey"`
}

func InitDir(path string) {
	if path == "" {
		path = acntDir
	} else {
		acntDir = path
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		// create the dir.
		if err := os.MkdirAll(path, 0777); err != nil {
			panic(err)
		}
	}
}

func IsExist(acntName string) bool {
	if _, err := os.Stat(filepath.Join(acntDir, acntName)); os.IsNotExist(err) {
		return false
	}
	return true
}

func New() *RpcAccount {
	p, s := cipher.GenerateKeyPair()
	return &RpcAccount{
		Pubkey: p,
		Seckey: s,
	}
}

// Load load acccount info from local disk.
func Load(acntName string) *RpcAccount {
	p := filepath.Join(acntDir, acntName)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		panic(err)
	}

	aj := AccountJson{}
	if err := util.LoadJSON(p, &aj); err != nil {
		panic(err)
	}
	a := RpcAccount{}
	copy(a.Pubkey[:], aj.Pubkey[0:33])
	copy(a.Seckey[:], aj.Pubkey[0:32])
	return &a
}

// StoreAccount store the rpcaccount to specific path.
func Store(acntName string, a RpcAccount) error {
	path := filepath.Join(acntDir, acntName)
	log.Println(path)
	aj := AccountJson{
		Pubkey: a.Pubkey[:],
		Seckey: a.Seckey[:],
	}
	return util.SaveJSON(path, aj, 0777)
}
