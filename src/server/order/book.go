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

func (bk *Book) AddBid(bid Order) {
	bk.bidMtx.Lock()
	bk.bidOrders = append(bk.bidOrders, bid)
	bk.bidMtx.Unlock()
}

func (bk *Book) AddAsk(ask Order) {
	bk.askMtx.Lock()
	bk.askOrders = append(bk.askOrders, ask)
	bk.askMtx.Unlock()
}

// Copy copy the bid and ask order list safely,
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

// Sort the bid and ask order list, it's not thread safe,
func (bk *Book) Sort() {
	sort.Sort(byPriceThenTime(bk.bidOrders))
	sort.Sort(byPriceThenTime(bk.askOrders))
}
