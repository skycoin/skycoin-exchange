package wallet

import (
	"fmt"

	"github.com/skycoin/skycoin-exchange/src/server/coin"
)

// BtcWallet bitcoin wallet.
type BtcWallet struct {
	walletBase
}

// GetCoinType return the wallet coin type.
func (bt BtcWallet) GetCoinType() coin.Type {
	return coin.Bitcoin
}

// SetID set wallet id
func (bt *BtcWallet) SetID(id string) {
	bt.ID = id
}

// InitSeed initialize the wallet seed.
func (bt *BtcWallet) InitSeed(seed string) {
	bt.InitSeed = seed
	bt.Seed = seed
}

// bitcoin wallet generator
func btcWltCreator(seed string) walletGentor {
	return func() Walleter {
		wlt := &BtcWallet{}
		wlt.ID = fmt.Sprintf("%s_%s", coin.Bitcoin, seed)
		wlt.Seed = seed
		wlt.InitSeed = seed
		return wlt
	}
}
