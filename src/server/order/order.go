package order

type Type uint8

const (
	Bid Type = iota
	Ask
)

type Order struct {
	ID          uint64 // order id.
	Type        Type   // order type.
	Price       uint64 // price of this order.
	Amount      uint64 // the amount of this order.
	CreatedTime uint64 // created time of the order.
}

type byPriceThenTime []Order

func (bp byPriceThenTime) Len() int {
	return len(bp)
}

func (bp byPriceThenTime) Less(i, j int) bool {
	a := bp[i]
	b := bp[j]
	if a.Price > b.Price {
		return true
	} else if a.Price == b.Price {
		return a.CreatedTime > b.CreatedTime
	}
	return false
}

func (bp byPriceThenTime) Swap(i, j int) {
	bp[i], bp[j] = bp[j], bp[i]
}
