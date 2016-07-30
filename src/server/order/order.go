package order

import "github.com/skycoin/skycoin-exchange/src/server/account"

type Type uint8

const (
	Bid Type = iota
	Ask
)

type Order struct {
	ID          uint64 // order id.
	AccountId   account.AccountID
	Type        Type   // order type.
	Price       uint64 // price of this order.
	Amount      uint64 // total amount of this order.
	RestAmt     uint64 // rest amount.
	CreatedTime uint64 // created time of the order.
}

type byPriceThenTimeDesc []Order
type byPriceThenTimeAsc []Order

func (bp byPriceThenTimeDesc) Len() int {
	return len(bp)
}

func (bp byPriceThenTimeDesc) Less(i, j int) bool {
	a := bp[i]
	b := bp[j]
	if a.Price > b.Price {
		return true
	} else if a.Price == b.Price {
		return a.CreatedTime > b.CreatedTime
	}
	return false
}

func (bp byPriceThenTimeDesc) Swap(i, j int) {
	bp[i], bp[j] = bp[j], bp[i]
}

func (bp byPriceThenTimeAsc) Len() int {
	return len(bp)
}

func (bp byPriceThenTimeAsc) Less(i, j int) bool {
	a := bp[i]
	b := bp[j]
	if a.Price < b.Price {
		return true
	} else if a.Price == b.Price {
		return a.CreatedTime > b.CreatedTime
	}
	return false
}

func (bp byPriceThenTimeAsc) Swap(i, j int) {
	bp[i], bp[j] = bp[j], bp[i]
}
