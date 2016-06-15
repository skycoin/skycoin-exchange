package wallet

import (
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
)

const WalletExt = "wlt"
const WalletTimestampFormat = "2006_01_02"

// Wallets, maintains all wallets in the server.
type Wallets struct {
	wallets map[string]Wallet // key of the map is wallet id.
	mtx     sync.Mutex        // protect concurrent access of the wallets.
}

var GWallets Wallets

func init() {
	GWallets.wallets = make(map[string]Wallet)
}

// NewWallet, create deterministic wallet.
func NewWallet(seed string) Wallet {
	return GWallets.NewWallet(seed, Deterministic)
}

// Create new wallet, it's thread safe.
func (self *Wallets) NewWallet(seed string, wltType WalletType) Wallet {
	for {
		wlt := newWallet(seed, wltType)
		id := newWalletID()
		self.mtx.Lock()
		if _, ok := self.wallets[id]; !ok {
			wlt.SetID(id)
			self.wallets[id] = wlt
			self.mtx.Unlock()
			return wlt
		}
		self.mtx.Unlock()
	}
}

func newWallet(seed string, wltType WalletType) Wallet {
	if seed == "" {
		seed = cipher.SumSHA256(cipher.RandByte(1024)).Hex()
	}

	switch wltType {
	case Deterministic:
		return &DeterministicWallet{
			Seed:      seed,
			Addresses: make(map[CoinType][]AddressEntry)}
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
