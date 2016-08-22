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
		return fmt.Errorf("%s does exist", wlt.GetID())
	}
	wlts.Value[wlt.GetID()] = wlt
	return wlt.Save()
}

func (wlts *wallets) remove(id string) error {
	wlts.mtx.Lock()
	defer wlts.mtx.Unlock()

	if wlt, ok := wlts.Value[id]; ok {
		if err := wlt.Clear(); err != nil {
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

	fmt.Println(wltDir)
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
			panic(fmt.Sprintf("%s type wallet not registered", tp))
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
		return wlt.NewAddresses(num)
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

func (wlts *wallets) getKeypair(id string, addr string) (pubkey, seckey, error) {
	wlts.mtx.Lock()
	defer wlts.mtx.Unlock()
	if wlt, ok := wlts.Value[id]; ok {
		return wlt.GetKeypair(addr), nil
	}
	return "", "", fmt.Errorf("%s wallet does not exist", id)
}
