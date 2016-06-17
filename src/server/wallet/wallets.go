package wallet

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
)

const WalletExt = "wlt"
const WalletTimestampFormat = "2006_01_02"

// Wallets, maintains all wallets in the server.
type Wallets struct {
	wallets map[string]Wallet // key of the map is wallet id.
	mtx     sync.Mutex        // protect concurrent access of the wallets.
}

var (
	GWallets Wallets
	dataDir  string = filepath.Join(util.UserHome(), ".skycoin-exchange/wallets")
)

func init() {
	util.InitDataDir(dataDir)
	wallets, err := loadWallets(dataDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	GWallets.wallets = wallets
}

// NewWallet, create deterministic wallet.
func NewWallet(seed string) (Wallet, error) {
	return GWallets.NewWallet(seed, Deterministic)
}

// Create new wallet
func (self *Wallets) NewWallet(seed string, wltType WalletType) (Wallet, error) {
	for {
		wlt := newWallet(seed, wltType)
		id := newWalletID()
		self.mtx.Lock()
		if _, ok := self.wallets[id]; !ok {
			wlt.SetID(id)
			self.wallets[id] = wlt
			self.mtx.Unlock()
			// save wallet.
			if err := wlt.Save(dataDir); err != nil {
				// remove from wallets
				self.RemoveWallet(id)
				return wlt, err
			}
			return wlt, nil
		}
		self.mtx.Unlock()
	}
}

func (self *Wallets) RemoveWallet(id string) {
	self.mtx.Lock()
	delete(self.wallets, id)
	self.mtx.Unlock()
}

// load wallets from local disk.
func loadWallets(dir string) (map[string]Wallet, error) {
	// TODO -- don't load duplicate wallets.
	// TODO -- save a last_modified value in wallets to decide which to load
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	//have := make(map[WalletID]Wallet, len(entries))
	wallets := make(map[string]Wallet)
	for _, e := range entries {
		if e.Mode().IsRegular() {
			name := e.Name()
			if !strings.HasSuffix(name, WalletExt) {
				continue
			}
			fullpath := filepath.Join(dir, name)
			w, err := loadWalletFromFile(fullpath)
			if err != nil {
				return nil, err
			}
			concretWlt, err := w.newConcretWallet()
			if err != nil {
				return nil, err
			}
			wallets[concretWlt.GetID()] = concretWlt
		}
	}
	return wallets, nil
}

func newWallet(seed string, wltType WalletType) Wallet {
	if seed == "" {
		seed = cipher.SumSHA256(cipher.RandByte(1024)).Hex()
	}

	switch wltType {
	case Deterministic:
		return &DeterministicWallet{
			Seed: seed,
			// WalletType:     WalletTypeStr[Deterministic],
			AddressEntries: make(map[string][]AddressEntry)}
	default:
		panic(fmt.Sprintf("unknow wallet type:%d", wltType))
	}
}

// newID, generate new wallet id, check for collisions and retry if failure.
func newWalletID() string {
	timestamp := time.Now().Format(WalletTimestampFormat)
	padding := hex.EncodeToString((cipher.RandByte(2)))
	return fmt.Sprintf("%s_%s.%s", timestamp, padding, WalletExt)
}
