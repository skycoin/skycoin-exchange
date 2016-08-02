package engine

import (
	"time"

	"github.com/skycoin/skycoin-exchange/src/server/account"
	"github.com/skycoin/skycoin-exchange/src/server/order"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

type Exchange interface {
	Run()
	CreateAccountWithPubkey(pubkey cipher.PubKey) (account.Accounter, error)
	GetAccount(id account.AccountID) (account.Accounter, error)
	GetBtcFee() uint64
	GetServPrivKey() cipher.SecKey
	WatchAddress(ct wallet.CoinType, addr string)
	GetNewAddress(coinType wallet.CoinType) string

	ChooseUtxos(ct wallet.CoinType, amount uint64, tm time.Duration) (interface{}, error)
	PutUtxos(ct wallet.CoinType, utxos interface{})

	GetAddrPrivKey(ct wallet.CoinType, addr string) (string, error)
	SaveAccount() error

	AddOrder(cp string, odr order.Order) (uint64, error)
	GetOrders(cp string, tp order.Type, start, end int64) ([]order.Order, error)
}
