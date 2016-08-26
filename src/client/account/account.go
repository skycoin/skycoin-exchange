package account

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
)

// Account client side account.
type Account struct {
	Pubkey string
	Seckey string
	WltIDs map[coin.Type]string // key: wallet type, value: wallet id.
}

type accountJSON struct {
	Pubkey string            `json:"pubkey"`
	Seckey string            `json:"seckey"`
	WltIDs map[string]string `json:"wallet_ids,omitempty"`
}

// internal global accounts.
var gAccounts accounts

// account storage dir.
var acntDir = filepath.Join(util.UserHome(), ".exchange-client/account")

// account file name
var acntName = "data.act"

// FileName return account storage file name.
func FileName() string {
	return acntName
}

// InitDir initialize account storage dir.
func InitDir(path string) {
	if path == "" {
		path = acntDir
	} else {
		acntDir = path
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		//create the dir.
		if err := os.MkdirAll(path, 0777); err != nil {
			panic(err)
		}
	}

	// load wallets.
	gAccounts = mustLoadAccounts(filepath.Join(acntDir, acntName))
}

// New create an account.
func New() Account {
	p, s := cipher.GenerateKeyPair()
	return Account{
		Pubkey: p.Hex(),
		Seckey: s.Hex(),
		WltIDs: make(map[coin.Type]string),
	}
}

// Set save and persist the account.
func Set(a Account) {
	gAccounts.set(a)
}

// Get return account of specific id.
func Get(id string) (Account, error) {
	for _, a := range gAccounts {
		if a.Pubkey == id {
			return a, nil
		}
	}
	return Account{}, errors.New("account does not exist")
}

type accounts []Account

func (acts *accounts) set(a Account) {
	var exist bool
	for i, act := range *acts {
		if act.Pubkey == a.Pubkey {
			exist = true
			(*acts)[i] = a
			break
		}
	}

	if !exist {
		*acts = append(*acts, a)
	}

	// persist the accounts
	acts.store()
}

func (acts *accounts) store() error {
	filename := filepath.Join(acntDir, acntName)
	tmpName := filename + ".tmp"
	f, err := os.Create(tmpName)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := acts.save(f); err != nil {
		return err
	}

	// check the existence of filename
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		if err := os.Rename(filename, filename+".bak"); err != nil {
			return err
		}
	}

	return os.Rename(tmpName, filename)
}

func (acts accounts) save(w io.Writer) error {
	v := struct {
		Account []accountJSON `json:"accounts"`
	}{
		acts.toJSON(),
	}

	d, err := json.MarshalIndent(&v, "", "    ")
	if err != nil {
		return err
	}

	io.Copy(w, bytes.NewBuffer(d))
	return nil
}

func (acts accounts) toJSON() []accountJSON {
	actsJSON := make([]accountJSON, len(acts))
	for i, a := range acts {
		aj := accountJSON{
			Pubkey: a.Pubkey,
			Seckey: a.Seckey,
			WltIDs: make(map[string]string),
		}

		for cp, id := range a.WltIDs {
			aj.WltIDs[cp.String()] = id
		}
		actsJSON[i] = aj
	}
	return actsJSON
}

func makeAccountsFromJSON(actsJSON []accountJSON) (accounts, error) {
	acts := make([]Account, len(actsJSON))
	for i, aj := range actsJSON {
		act := Account{
			Pubkey: aj.Pubkey,
			Seckey: aj.Seckey,
			WltIDs: make(map[coin.Type]string),
		}
		for cp, id := range aj.WltIDs {
			cpt, err := coin.TypeFromStr(cp)
			if err != nil {
				return accounts{}, err
			}
			act.WltIDs[cpt] = id
		}
		acts[i] = act
	}
	return acts, nil
}

func mustLoadAccounts(filename string) accounts {
	acts := accounts{}
	// check the existence of the file.
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return acts
	}

	d, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	v := struct {
		Acounts []accountJSON `json:"accounts"`
	}{}
	if err := json.Unmarshal(d, &v); err != nil {
		panic(err)
	}

	acts, err = makeAccountsFromJSON(v.Acounts)
	if err != nil {
		panic(err)
	}
	return acts
}
