package engine

import (
	"time"

	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/server/account"
	"github.com/skycoin/skycoin-exchange/src/server/order"
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
	IsAdmin(pubkey string) bool
}

type Addresser interface {
	WatchAddress(ct, addr string)
	GetNewAddress(coinType string) string
	GetAddrPrivKey(ct, addr string) (string, error)
}

type Order interface {
	AddOrder(cp string, odr order.Order) (uint64, error)
	GetOrders(cp string, tp order.Type, start, end int64) ([]order.Order, error)
}

type Utxor interface {
	ChooseUtxos(ct string, amount uint64, tm time.Duration) (interface{}, error)
	PutUtxos(ct string, utxos interface{})
}

type Server interface {
	Run()
	GetSecKey() string
	GetBtcFee() uint64
	GetSupportCoins() []string
	GetCoin(ct string) (coin.Gateway, error)
	BindCoins(cs ...coin.Gateway) error
}
