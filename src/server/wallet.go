package server

import (
	"fmt"
	"sync"

	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/wallet"
)

// wallets wrap up the wallet package.
type wallets struct {
	ids map[string]string // key wallet type, value wallet id.
}

type walletItem struct {
	Type string // coin type
	Seed string // seed
}

var initWalletOnce sync.Once

func makeWallets(dir string, items []walletItem) (wallets, error) {
	f := func() {
		logger.Debug("wallet dir:%s", dir)
		wallet.InitDir(dir)
	}
	initWalletOnce.Do(f)
	wlts := wallets{ids: make(map[string]string)}
	// create wallets if not exist.
	for _, item := range items {
		id := wallet.MakeWltID(item.Type, item.Seed)
		if !wallet.IsExist(id) {
			_, err := wallet.New(item.Type, item.Seed)
			if err != nil {
				return wallets{}, err
			}
		}
		wlts.ids[item.Type] = id
	}
	return wlts, nil
}

// NewAddresses create specific coin addresses.
func (wlts *wallets) NewAddresses(cp string, num int) ([]coin.AddressEntry, error) {
	if id, ok := wlts.ids[cp]; ok {
		return wallet.NewAddresses(id, num)
	}
	return []coin.AddressEntry{}, fmt.Errorf("%s wallet not supported", cp)
}

// GetKeypair get pub/sec keys of specific address.
func (wlts wallets) GetKeypair(cp string, addr string) (string, string, error) {
	if id, ok := wlts.ids[cp]; ok {
		return wallet.GetKeypair(id, addr)
	}
	return "", "", fmt.Errorf("%s wallet not supported", cp)
}

// GetAddresses get all addresses in one specific wallet.
func (wlts wallets) GetAddresses(cp string) ([]string, error) {
	if id, ok := wlts.ids[cp]; ok {
		return wallet.GetAddresses(id)
	}
	return []string{}, fmt.Errorf("%s wallet not supported", cp)
}
