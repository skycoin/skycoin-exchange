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

// NewAddresses, create addresses in specific wallet.
func NewAddresses(wltID string, cointype CoinType, num int) ([]string, error) {
	wlt, err := GWallets.GetWallet(wltID)
	if err != nil {
		return []string{}, err
	}

	addrEntries, err := wlt.NewAddresses(cointype, num)
	if err != nil {
		return []string{}, err
	}
	addrs := make([]string, len(addrEntries))
	for i, e := range addrEntries {
		addrs[i] = e.Address
	}
	return addrs, nil
}

// GetBalance, query balance of specific address.
// func GetBalance(addr string, cointype CoinType) (string, error) {
// 	switch cointype {
// 	case Bitcoin:
// 		return bitcoin.GetBalance(addr)
// 	default:
// 		return "", fmt.Errorf("unknow cointype:%d", cointype)
// 	}
// }
//
// func GetUnspentOutputs(addr string, cointype CoinType) ([]bitcoin.UnspentOutputJSON, error) {
// 	switch cointype {
// 	case Bitcoin:
// 		return bitcoin.GetUnspentOutputs(addr), nil
// 	default:
// 		return []bitcoin.UnspentOutputJSON{}, fmt.Errorf("unknow cointype:%d", cointype)
// 	}
// }

// Create new wallet
func (self *Wallets) NewWallet(seed string, wltType WalletType) (Wallet, error) {
	for {
		// wlt := newWallet(seed, wltType)
		id := newWalletID()
		self.mtx.Lock()
		if _, ok := self.wallets[id]; !ok {
			wlt, err := newWallet(id, seed, wltType)
			if err != nil {
				self.mtx.Unlock()
				return nil, err
			}
			self.wallets[id] = wlt
			self.mtx.Unlock()
			return wlt, nil
		}
		self.mtx.Unlock()
	}
}

func (self *Wallets) GetWallet(id string) (Wallet, error) {
	self.mtx.Lock()
	defer self.mtx.Unlock()
	if _, ok := self.wallets[id]; !ok {
		return nil, fmt.Errorf("invalide wallet id:%s", id)
	}
	return self.wallets[id], nil
}

// removeWallet, this function should not be called manual.
func (self *Wallets) removeWallet(id string) {
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

func newWallet(id string, seed string, wltType WalletType) (Wallet, error) {
	if seed == "" {
		seed = cipher.SumSHA256(cipher.RandByte(1024)).Hex()
	}

	switch wltType {
	case Deterministic:
		return &DeterministicWallet{
			ID:             id,
			Seed:           seed,
			InitSeed:       seed,
			AddressEntries: make(map[string][]AddressEntry)}, nil
	default:
		return nil, fmt.Errorf("newWallet fail, unknow wallet type:%d", wltType)
	}
}

// newID, generate new wallet id, check for collisions and retry if failure.
func newWalletID() string {
	timestamp := time.Now().Format(WalletTimestampFormat)
	padding := hex.EncodeToString((cipher.RandByte(2)))
	return fmt.Sprintf("%s_%s.%s", timestamp, padding, WalletExt)
}
