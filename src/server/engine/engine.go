package engine

import (
	"time"

	"github.com/skycoin/skycoin-exchange/src/server/account"
	"github.com/skycoin/skycoin-exchange/src/server/coin_interface/bitcoin"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

type Exchange interface {
	Run()
	CreateAccountWithPubkey(pubkey cipher.PubKey) (account.Accounter, error)
	GetAccount(id account.AccountID) (account.Accounter, error)
	GetFee() uint64
	GetServPrivKey() cipher.SecKey
	AddWatchAddress(ct wallet.CoinType, addr string)
	GetNewAddress(coinType wallet.CoinType) string
	ChooseUtxos(coinType wallet.CoinType, amount uint64, tm time.Duration) ([]bitcoin_interface.UtxoWithkey, error)
	PutUtxos(ct wallet.CoinType, utxos []bitcoin_interface.UtxoWithkey)
}
