package server

import "github.com/skycoin/skycoin-exchange/src/server/wallet"

type UtxoManager interface {
	GetUtxo(coinType wallet.CoinType) chan bitcoin.Utxo  // get utxo from utxo pool
	PutUtxo(coinType wallet.CoinType, utxo bitcoin.Utxo) // put utxo into utxo pool
	AddWatchAddress(ct wallet.CoinType, addr string)
	GetAddressOfUtxo(ct wallet.CoinType, utxo bitcoin.Utxo) (string, error)
}
