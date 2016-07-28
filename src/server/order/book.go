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
