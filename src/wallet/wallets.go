package wallet

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/skycoin/skycoin-exchange/src/coin"
)

// wallets record all wallet, key is wallet id, value is wallet interface.
type wallets struct {
	mtx   sync.Mutex
	Value map[string]Walleter
}

// internal global wallets
var gWallets = wallets{Value: make(map[string]Walleter)}

func (wlts *wallets) add(wlt Walleter) error {
	wlts.mtx.Lock()
	defer wlts.mtx.Unlock()
	if _, ok := wlts.Value[wlt.GetID()]; ok {
		return fmt.Errorf("%s already exist", wlt.GetID())
	}
	wlts.Value[wlt.GetID()] = wlt
	return wlts.store(wlt)
}

func (wlts *wallets) remove(id string) error {
	wlts.mtx.Lock()
	defer wlts.mtx.Unlock()

	if wlt, ok := wlts.Value[id]; ok {
		path := storeAddr(wlt)
		if err := os.RemoveAll(path); err != nil {
			return err
		}
		delete(wlts.Value, id)
	}
	return nil
}

func (wlts *wallets) reset() {
	wlts.mtx.Lock()
	wlts.Value = make(map[string]Walleter)
	wlts.mtx.Unlock()
}

// load from local disk
func (wlts *wallets) mustLoad() {
	// clear wallets in memory.
	wlts.reset()

	fileInfos, _ := ioutil.ReadDir(wltDir)
	for _, fileInfo := range fileInfos {
		name := fileInfo.Name()
		if !strings.HasSuffix(name, ".wlt") {
			continue
		}
		// get the wallet type, the name: $bitcoin_$seed1234.wlt
		typeSeed := strings.SplitN(name, "_", 2)
		if len(typeSeed) != 2 {
			panic("error wallet file name")
		}

		// check coin type
		tp, err := coin.TypeFromStr(typeSeed[0])
		if err != nil {
			panic(err)
		}

		newWlt, ok := gWalletCreators[tp]
		if !ok {
			panic(fmt.Sprintf("%s wallet not supported", tp))
		}

		f, err := os.Open(filepath.Join(wltDir, name))
		if err != nil {
			panic(err)
		}
		defer f.Close()

		wlt := newWlt()
		if err := wlt.Load(f); err != nil {
			panic(err)
		}
		if err := wlts.add(wlt); err != nil {
			panic(err)
		}
	}
}

func (wlts *wallets) newAddresses(id string, num int) ([]coin.AddressEntry, error) {
	wlts.mtx.Lock()
	defer wlts.mtx.Unlock()
	if wlt, ok := wlts.Value[id]; ok {
		addrs, err := wlt.NewAddresses(num)
		if err != nil {
			return []coin.AddressEntry{}, err
		}

		if err := wlts.store(wlt); err != nil {
			return []coin.AddressEntry{}, err
		}
		return addrs, nil
	}
	return []coin.AddressEntry{}, fmt.Errorf("%s wallet does not exist", id)
}

func (wlts *wallets) getAddresses(id string) ([]string, error) {
	wlts.mtx.Lock()
	defer wlts.mtx.Unlock()
	if wlt, ok := wlts.Value[id]; ok {
		return wlt.GetAddresses(), nil
	}
	return []string{}, fmt.Errorf("%s wallet does not exist", id)
}

func (wlts *wallets) isContain(id string, addrs []string) (bool, error) {
	wlts.mtx.Lock()
	defer wlts.mtx.Unlock()
	if wlt, ok := wlts.Value[id]; ok {
		as := wlt.GetAddresses()
		for _, addr := range addrs {
			have := func(addr string) bool {
				for _, a := range as {
					if a == addr {
						return true
					}
				}
				return false
			}
			if !have(addr) {
				return false, nil
			}
		}
		return true, nil
	}
	return false, fmt.Errorf("wallet %s does not exist", id)
}

func (wlts *wallets) getKeypair(id string, addr string) (string, string, error) {
	wlts.mtx.Lock()
	defer wlts.mtx.Unlock()
	if wlt, ok := wlts.Value[id]; ok {
		return wlt.GetKeypair(addr)
	}
	return "", "", fmt.Errorf("%s wallet does not exist", id)
}

func (wlts *wallets) store(wlt Walleter) error {
	path := storeAddr(wlt)
	tmpPath := path + "." + "tmp"

	// write wallet to temp file.
	f, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := wlt.Save(f); err != nil {
		return err
	}

	// create bak file if exist.
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if err := os.Rename(path, path+".bak"); err != nil {
			return err
		}
	}

	return os.Rename(tmpPath, path)
}

func (wlts *wallets) isExist(id string) bool {
	wlts.mtx.Lock()
	defer wlts.mtx.Unlock()
	if _, ok := wlts.Value[id]; ok {
		return true
	}
	return false
}

func storeAddr(wlt Walleter) string {
	return filepath.Join(wltDir, wlt.GetID()+"."+Ext)
}
