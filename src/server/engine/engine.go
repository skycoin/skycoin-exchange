package engine

import (
	"time"

	"github.com/skycoin/skycoin-exchange/src/server/account"
	"github.com/skycoin/skycoin-exchange/src/server/order"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

type Exchange interface {
	Server
	Accounter
	Addresser
	Order
	Utxor
}

type Accounter interface {
	CreateAccountWithPubkey(pubkey string) (account.Accounter, error)
	GetAccount(id string) (account.Accounter, error)
	SaveAccount() error
}

type Addresser interface {
	WatchAddress(ct wallet.CoinType, addr string)
	GetNewAddress(coinType wallet.CoinType) string
	GetAddrPrivKey(ct wallet.CoinType, addr string) (string, error)
}

type Order interface {
	AddOrder(cp string, odr order.Order) (uint64, error)
	GetOrders(cp string, tp order.Type, start, end int64) ([]order.Order, error)
}

type Utxor interface {
	ChooseUtxos(ct wallet.CoinType, amount uint64, tm time.Duration) (interface{}, error)
	PutUtxos(ct wallet.CoinType, utxos interface{})
}

type Server interface {
	Run()
	GetBtcFee() uint64
	GetServPrivKey() cipher.SecKey
	GetSupportCoins() []string
}
