package wallet

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin/src/util/file"
)

// Walleter interface, new wallet type can be supported if it fullfills this interface.
type Walleter interface {
	GetID() string                                     // get wallet id.
	SetID(id string)                                   // set wallet id.
	SetSeed(seed string)                               // init the wallet seed.
	GetType() string                                   // get the wallet coin type.
	NewAddresses(num int) ([]coin.AddressEntry, error) // generate new addresses.
	GetAddresses() []string                            // get all addresses in the wallet.
	GetKeypair(addr string) (string, string, error)    // get pub/sec key pair of specific address
	Save(w io.Writer) error                            // save the wallet.
	Load(r io.Reader) error                            // load wallet from reader.
	Copy() Walleter                                    // copy of self, for thread safe.
}

// wltDir default wallet dir, wallet file name sturct: $type_$seed.wlt.
// example: bitcoin_seed.wlt, skycoin_seed.wlt.
var wltDir = filepath.Join(file.UserHome(), ".exchange-client/wallet")

// Ext wallet file extension name
var Ext = "wlt"

// Creator wallet creator.
type Creator func() Walleter

var gWalletCreators = make(map[string]Creator)

// RegisterCreator when new type wallet need to be supported,
// the wallet must provide Creator, and use this function to register the creator into system.
func RegisterCreator(tp string, ctor Creator) error {
	if _, ok := gWalletCreators[tp]; ok {
		return fmt.Errorf("%s wallet already registered", tp)
	}
	gWalletCreators[tp] = ctor
	return nil
}

// InitDir initialize the wallet file storage dir,
// load wallets if exist.
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

// GetWalletDir return the current wallet dir.
func GetWalletDir() string {
	return wltDir
}

// New create wallet base on seed and coin type.
func New(tp, seed string) (Walleter, error) {
	newWlt, ok := gWalletCreators[tp]
	if !ok {
		return nil, fmt.Errorf("%s wallet not regestered", tp)
	}

	// create wallet base on the wallet creator.
	wlt := newWlt()
	wlt.SetID(MakeWltID(tp, seed))
	wlt.SetSeed(seed)

	if err := gWallets.add(wlt); err != nil {
		return nil, err
	}
	return wlt.Copy(), nil
}

// IsExist check if the wallet is already exist.
func IsExist(id string) bool {
	return gWallets.isExist(id)
}

// MakeWltID make wallet id base on coin type and seed.
func MakeWltID(cp, seed string) string {
	return fmt.Sprintf("%s_%s", cp, seed)
}

// NewAddresses create address
func NewAddresses(id string, num int) ([]coin.AddressEntry, error) {
	return gWallets.newAddresses(id, num)
}

// GetAddresses get all addresses in specific wallet.
func GetAddresses(id string) ([]string, error) {
	return gWallets.getAddresses(id)
}

// IsContain check if the addresses are int the wallet.
func IsContain(id string, addrs []string) (bool, error) {
	return gWallets.isContain(id, addrs)
}

// GetKeypair get pub/sec key pair of specific addresse in wallet.
func GetKeypair(id string, addr string) (string, string, error) {
	return gWallets.getKeypair(id, addr)
}

// Remove remove wallet of specific id.
func Remove(id string) error {
	return gWallets.remove(id)
}
