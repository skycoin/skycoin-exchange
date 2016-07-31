package order

import (
	"sort"
	"sync"
)

// order book, which records the bid and ask order list.
type Book struct {
	bidOrders []Order
	askOrders []Order
	bidMtx    sync.Mutex
	askMtx    sync.Mutex
}

type BookJson struct {
	BidOrders []Order `json:"bids"`
	AskOrders []Order `json:"asks"`
}

type OrderPair struct {
	Bid Order
	Ask Order
}

func (bk *Book) AddBid(bid Order) {
	bk.bidMtx.Lock()
	bk.bidOrders = append(bk.bidOrders, bid)
	sort.Sort(byPriceThenTimeDesc(bk.bidOrders))
	bk.bidMtx.Unlock()
}

func (bk *Book) AddAsk(ask Order) {
	bk.askMtx.Lock()
	bk.askOrders = append(bk.askOrders, ask)
	sort.Sort(byPriceThenTimeAsc(bk.askOrders))
	bk.askMtx.Unlock()
}

func (bk *Book) Copy() Book {
	newBk := Book{}
	bk.bidMtx.Lock()
	newBk.bidOrders = make([]Order, len(bk.bidOrders))
	copy(newBk.bidOrders, bk.bidOrders)
	bk.bidMtx.Unlock()

	bk.askMtx.Lock()
	newBk.askOrders = make([]Order, len(bk.askOrders))
	copy(newBk.askOrders, bk.askOrders)
	bk.askMtx.Unlock()
	return newBk
}

func (bk Book) getMaxOrderID() uint64 {
	// sort the book with priority of order id.
	orders := append(bk.bidOrders, bk.askOrders...)
	sort.Sort(byOrderID(orders))
	return orders[0].ID
}

// Match check if there're bids and asks are matched,
// if matched, remove from the order book, and return the orders for
// further use.
func (bk *Book) Match() []Order {
	bk.bidMtx.Lock()
	bk.askMtx.Lock()
	if len(bk.bidOrders) == 0 || len(bk.askOrders) == 0 {
		bk.askMtx.Unlock()
		bk.bidMtx.Unlock()
		return []Order{}
	}

	// the highest buy price < the lowest sell price, no order match.
	if bk.bidOrders[0].Price < bk.askOrders[0].Price {
		bk.askMtx.Unlock()
		bk.bidMtx.Unlock()
		return []Order{}
	}

	// var bidIndex, askIndex int = 0, 0
	var bidOrders []Order
	var askOrders []Order

	for i, bid := range bk.bidOrders {
		restAmt, askOrderNum := checkAskOrders(bid, &bk.askOrders)
		if restAmt == bid.Amount {
			break
		}

		bk.bidOrders[i].RestAmt = restAmt

		// append fullfilled ask orders
		askOrders = append(askOrders, bk.askOrders[:askOrderNum]...)
		// remove fullfilled ask orders from book.
		bk.askOrders = bk.askOrders[askOrderNum:]

		if restAmt == 0 {
			bidOrders = append(bidOrders, bk.bidOrders[i])
		} else if restAmt > 0 {
			break
		}
	}

	bk.bidOrders = bk.bidOrders[len(bidOrders):]
	bk.askMtx.Unlock()
	bk.bidMtx.Unlock()

	return append(bidOrders, askOrders...)
}

func (bk Book) ToMarshalable() BookJson {
	bj := BookJson{
		BidOrders: make([]Order, len(bk.bidOrders)),
		AskOrders: make([]Order, len(bk.askOrders)),
	}

	copy(bj.BidOrders, bk.bidOrders)
	copy(bj.AskOrders, bk.askOrders)
	return bj
}

func NewBookFromJson(bj BookJson) *Book {
	bk := &Book{
		bidOrders: make([]Order, len(bj.BidOrders)),
		askOrders: make([]Order, len(bj.AskOrders)),
	}

	copy(bk.bidOrders, bj.BidOrders)
	copy(bk.askOrders, bj.AskOrders)
	return bk
}

func checkAskOrders(bid Order, askOrders *[]Order) (uint64, uint64) {
	if bid.RestAmt == 0 {
		panic("the bid amount already fullfilled")
	}

	var askNum uint64
	for i, ask := range *askOrders {
		// if ask.RestAmt > 0 {
		if bid.Price < ask.Price {
			return bid.RestAmt, askNum
		}

		if bid.RestAmt < ask.RestAmt {
			(*askOrders)[i].RestAmt -= bid.RestAmt
			return 0, 0
		} else if bid.RestAmt == ask.RestAmt {
			(*askOrders)[i].RestAmt = 0
			askNum += 1
			return 0, askNum
		} else if bid.RestAmt > ask.RestAmt {
			bid.RestAmt -= ask.RestAmt
			(*askOrders)[i].RestAmt = 0
			askNum += 1
		}
	}
	return bid.RestAmt, askNum
}
