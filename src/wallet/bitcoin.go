package wallet

import "github.com/skycoin/skycoin-exchange/src/coin"

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
func btcWltCreator() walletGentor {
	return func() Walleter {
		return &BtcWallet{}
	}
}
