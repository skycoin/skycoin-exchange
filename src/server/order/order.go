package order

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/skycoin/skycoin/src/util/file"
)

type Type uint8

const (
	Bid Type = iota
	Ask
)

var (
	orderDir string = filepath.Join(file.UserHome(), ".skycoin-exchange/orderbook")
	orderExt string = "ods"
	idExt    string = "id"
)

type Order struct {
	ID        uint64 `json:"id"` // order id.
	AccountID string `json:"account_id"`
	Type      Type   `json:"type"`       // order type.
	Price     uint64 `json:"price"`      // price of this order.
	Amount    uint64 `json:"amount"`     // total amount of this order.
	RestAmt   uint64 `json:"reset_amt"`  // rest amount.
	CreatedAt int64  `json:"created_at"` // created time of the order.
}

type byPriceThenTimeDesc []Order
type byPriceThenTimeAsc []Order
type byOrderID []Order

func (bp byPriceThenTimeDesc) Len() int {
	return len(bp)
}

func (bp byPriceThenTimeDesc) Less(i, j int) bool {
	a := bp[i]
	b := bp[j]
	if a.Price > b.Price {
		return true
	} else if a.Price == b.Price {
		return a.CreatedAt > b.CreatedAt
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
		return a.CreatedAt > b.CreatedAt
	}
	return false
}

func (bp byPriceThenTimeAsc) Swap(i, j int) {
	bp[i], bp[j] = bp[j], bp[i]
}

func (bo byOrderID) Len() int {
	return len(bo)
}

func (bo byOrderID) Less(i, j int) bool {
	return bo[i].ID > bo[j].ID
}

func (bo byOrderID) Swap(i, j int) {
	bo[i], bo[j] = bo[j], bo[i]
}

func InitDir(path string) {
	if path == "" {
		path = orderDir
	} else {
		orderDir = path
	}
	// create the account dir if not exist.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			panic(err)
		}
	}
}

func New(aid string, tp Type, price uint64, amount uint64) *Order {
	return &Order{
		AccountID: aid,
		Type:      tp,
		Price:     price,
		Amount:    amount,
		RestAmt:   amount,
		CreatedAt: time.Now().Unix(),
	}
}

func (tp Type) String() string {
	switch tp {
	case Bid:
		return "bid"
	case Ask:
		return "ask"
	default:
		return ""
	}
}

func TypeFromStr(tp string) (Type, error) {
	switch tp {
	case "bid":
		return Bid, nil
	case "ask":
		return Ask, nil
	default:
		return 0, fmt.Errorf("unknow order type:%s", tp)
	}
}
