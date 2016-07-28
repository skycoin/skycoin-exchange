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
