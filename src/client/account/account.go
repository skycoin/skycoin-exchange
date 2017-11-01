package account

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	logging "github.com/op/go-logging"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/file"
)

var logger = logging.MustGetLogger("client.account")

// Account client side account.
type Account struct {
	Pubkey string
	Seckey string
	WltIDs map[string]string // key: wallet type, value: wallet id.
}

type accountJSON struct {
	Pubkey string            `json:"pubkey"`
	Seckey string            `json:"seckey"`
	WltIDs map[string]string `json:"wallet_ids,omitempty"`
}

// internal global accounts.
var gAccounts manager

// account storage dir.
var acntDir = filepath.Join(file.UserHome(), ".exchange-client/account")

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
		WltIDs: make(map[string]string),
	}
}

// Set save and persist the account.
func Set(a Account) {
	logger.Debug("active account:%s", a.Pubkey)
	gAccounts.set(a)
}

// Get return account of specific id.
func Get(pubkey string) (Account, error) {
	for _, a := range gAccounts.Accounts {
		if a.Pubkey == pubkey {
			return a, nil
		}
	}
	return Account{}, errors.New("account does not exist")
}

// GetAll get all accounts
func GetAll() []Account {
	return gAccounts.Accounts
}

// GetActive get the current working account.
func GetActive() (Account, error) {
	if gAccounts.ActiveAcount.Pubkey == "" || gAccounts.ActiveAcount.Seckey == "" {
		return Account{}, errors.New("no active account")
	}
	return gAccounts.ActiveAcount, nil
}

// SetActive set the account as active account.
func SetActive(pubkey string) error {
	a, err := Get(pubkey)
	if err != nil {
		return err
	}

	gAccounts.ActiveAcount = a
	return nil
}

type manager struct {
	ActiveAcount Account   // current working account.
	Accounts     []Account // all accounts
}

func (mgr *manager) set(a Account) {
	mgr.ActiveAcount = a
	var exist bool
	for i, act := range mgr.Accounts {
		if act.Pubkey == a.Pubkey {
			exist = true
			mgr.Accounts[i] = a
			break
		}
	}

	if !exist {
		mgr.Accounts = append(mgr.Accounts, a)
	}

	// persist the accounts
	mgr.store()
}

func (mgr *manager) store() error {
	filename := filepath.Join(acntDir, acntName)
	tmpName := filename + ".tmp"
	f, err := os.Create(tmpName)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := mgr.save(f); err != nil {
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

func (mgr manager) save(w io.Writer) error {
	v := struct {
		Account []accountJSON `json:"accounts"`
		Active  string        `json:"active_account"`
	}{
		accounts(mgr.Accounts).toJSON(),
		mgr.ActiveAcount.Pubkey,
	}

	d, err := json.MarshalIndent(&v, "", "    ")
	if err != nil {
		return err
	}

	io.Copy(w, bytes.NewBuffer(d))
	return nil
}

type accounts []Account

func (acts accounts) toJSON() []accountJSON {
	actsJSON := make([]accountJSON, len(acts))
	for i, a := range acts {
		aj := accountJSON{
			Pubkey: a.Pubkey,
			Seckey: a.Seckey,
			WltIDs: make(map[string]string),
		}

		for cp, id := range a.WltIDs {
			aj.WltIDs[cp] = id
		}
		actsJSON[i] = aj
	}
	return actsJSON
}

func makeAccountsFromJSON(actsJSON []accountJSON) ([]Account, error) {
	acts := make([]Account, len(actsJSON))
	for i, aj := range actsJSON {
		act := Account{
			Pubkey: aj.Pubkey,
			Seckey: aj.Seckey,
			WltIDs: make(map[string]string),
		}
		for cp, id := range aj.WltIDs {
			act.WltIDs[cp] = id
		}
		acts[i] = act
	}
	return acts, nil
}

func mustLoadAccounts(filename string) manager {
	mgr := manager{}
	// check the existence of the file.
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return mgr
	}

	d, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	v := struct {
		ActiveAccount string        `json:"active_account"`
		Acounts       []accountJSON `json:"accounts"`
	}{}

	if err := json.Unmarshal(d, &v); err != nil {
		panic(err)
	}

	acts, err := makeAccountsFromJSON(v.Acounts)
	if err != nil {
		panic(err)
	}

	mgr.Accounts = acts
	mgr.ActiveAcount = func() Account {
		for _, a := range acts {
			if a.Pubkey == v.ActiveAccount {
				return a
			}
		}
		return Account{}
	}()

	return mgr
}
