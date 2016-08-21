package wallet

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/skycoin/skycoin-exchange/src/server/coin"
	"github.com/skycoin/skycoin/src/util"
)

// Walleter interface, new wallet type can be supported if it fullfills this interface.
type Walleter interface {
	GetID() string                                            // get wallet id.
	GetCoinType() coin.Type                                   // get current wallet coin type.
	NewAddresses(num int) ([]coin.AddressEntry, error)        // generate new address
	GetAddresses() []string                                   // get all addresses in the wallet.
	GetAddressEntries() []coin.AddressEntry                   // get all address enties in the wallet.
	GetAddrEntyByAddr(addr string) (coin.AddressEntry, error) // get address enty by address
	Save() error                                              // save the wallet into local disk.
	Clear() error                                             // remove the wallet file from local disk.
}

// default wallet dir, wallet file name sturct: $seedname_$type.wlt.
// example: btc_seed.wlt, sky_seed.wlt.
var wltDir = filepath.Join(util.UserHome(), ".exchange-client/wallet")

// InitDir initialize the wallet file storage dir.
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

	gWallets.mustLoad()
}

// GetWalletDir return the current client wallet dir.
func GetWalletDir() string {
	return wltDir
}

// New create wallet base on seed and coin type.
func New(tp coin.Type, seed string) (Walleter, error) {
	var wlt Walleter
	switch tp {
	case coin.Bitcoin:
		wlt = &BtcWallet{
			walletBase: walletBase{
				ID:       fmt.Sprintf("btc_%s.wlt", seed),
				Type:     tp.String(),
				InitSeed: seed,
				Seed:     seed},
		}
	case coin.Skycoin:
		return nil, nil
	default:
		return nil, errors.New("unknow wallet coin type")
	}

	if err := gWallets.add(wlt); err != nil {
		return nil, err
	}

	return wlt, nil
}

type walletBase struct {
	ID             string              `json:"id"` // wallet id
	Type           string              `json:"type"`
	InitSeed       string              `json:"init_seed"`         // Init seed, used to recover the wallet.
	Seed           string              `json:"seed"`              // used to track the latset seed
	AddressEntries []coin.AddressEntry `json:"entries,omitempty"` // address entries.
}

var gWallets = wallets{Value: make(map[string]Walleter)}

type wallets struct {
	Value map[string]Walleter
}

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
func mustLoad() *wallets {
	fileInfos, _ := ioutil.ReadDir(wltDir)
	for _, fileInfo := range fileInfos {
		if !strings.HasSuffix(fileInfo.Name(), ".wlt") {
			continue
		}

	}

	// err := util.LoadJSON(filename, &w)
	// if err != nil {
	// 	return WalletBase{}, err
	// }
	// return w, nil
}
