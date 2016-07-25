package engine

import (
	"time"

	"github.com/skycoin/skycoin-exchange/src/server/account"
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

	ChooseUtxos(ct wallet.CoinType, amount uint64, tm time.Duration) (interface{}, error)
	PutUtxos(ct wallet.CoinType, utxos interface{})

	GetPrivKey(ct wallet.CoinType, addr string) (string, error)
	SaveAccount() error
}
