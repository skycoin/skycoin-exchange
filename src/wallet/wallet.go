package wallet

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin/src/util"
)

// Walleter interface, new wallet type can be supported if it fullfills this interface.
type Walleter interface {
	GetID() string                                     // get wallet id.
	SetID(id string)                                   // set wallet id.
	InitSeed(seed string)                              // init the wallet seed.
	GetCoinType() coin.Type                            // get the wallet coin type.
	NewAddresses(num int) ([]coin.AddressEntry, error) // generate new addresses.
	GetAddresses() []string                            // get all addresses in the wallet.
	GetKeypair(addr string) (string, string)           // get pub/sec key pair of specific address
	Save() error                                       // save the wallet.
	Load(r io.Reader) error                            // load wallet from reader.
	Clear() error                                      // remove the wallet file from local disk.
}

// wltDir default wallet dir, wallet file name sturct: $type_$seed.wlt.
// example: btc_seed.wlt, sky_seed.wlt.
var wltDir = filepath.Join(util.UserHome(), ".exchange-client/wallet")
var wltExt = "wlt"

// InitDir initialize the wallet file storage dir,
// load wallets in the dir if it does exist.
func InitDir(path string) {
	if path == "" {
		path = wltDir
	} else {
		wltDir = path
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		//create the dir.
		if err := os.MkdirAll(path, 0777); err != nil {
			panic(err)
		}
	}

	// load wallets.
	gWallets.mustLoad()
}

// GetWalletDir return the current client wallet dir.
func GetWalletDir() string {
	return wltDir
}

// New create wallet base on seed and coin type.
func New(tp coin.Type, seed string) (Walleter, error) {
	newWlt, ok := gWalletCreators[tp]
	if !ok {
		return nil, fmt.Errorf("%s coin wallet not regestered", tp)
	}

	// create wallet base on the wallet creator.
	wlt := newWlt()
	wlt.SetID(fmt.Sprintf("%s_%s", tp, seed))
	wlt.InitSeed(seed)

	if err := gWallets.add(wlt); err != nil {
		return nil, err
	}

	return wlt, nil
}

// wallet creator.
type walletCreator func() Walleter

var gWalletCreators = make(map[coin.Type]walletCreator)

func init() {
	// register bitcoin wallet creator
	gWalletCreators[coin.Bitcoin] = btcWltCreator()
}

// wallets record all wallet, key is wallet id, value is wallet interface.
type wallets struct {
	Value map[string]Walleter
}

// internal global wallets
var gWallets = wallets{Value: make(map[string]Walleter)}

func (wlts *wallets) add(wlt Walleter) error {
	if _, ok := wlts.Value[wlt.GetID()]; ok {
		return fmt.Errorf("%s does exist", wlt.GetID())
	}
	wlts.Value[wlt.GetID()] = wlt
	return wlt.Save()
}

func (wlts *wallets) remove(id string) error {
	if wlt, ok := wlts.Value[id]; ok {
		return wlt.Clear()
	}
	delete(wlts.Value, id)
	return nil
}

// load from local disk
func (wlts *wallets) mustLoad() {
	fileInfos, _ := ioutil.ReadDir(wltDir)
	for _, fileInfo := range fileInfos {
		name := fileInfo.Name()
		if !strings.HasSuffix(name, ".wlt") {
			continue
		}
		// get the wallet type, the name: $bitcoin_$seed1234.wlt
		typeSeed := strings.Split(name, "_")
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

		f, err := os.OpenFile(filepath.Join(wltDir, name))
		if err != nil {
			panic(err)
		}
		defer f.Close()

		wlt := newWlt()
		if err := wlt.Load(f); err != nil {
			panic(err)
		}
		if err := wallets.add(wlt); err != nil {
			panic(err)
		}
	}
}
