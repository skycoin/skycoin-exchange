package order

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBook(t *testing.T) {
	var BidOrderList = []Order{
		Order{Price: 100, CreatedAt: 132424, Amount: 1},
		Order{Price: 102, CreatedAt: 132425, Amount: 1},
		Order{Price: 103, CreatedAt: 132428, Amount: 1},
		Order{Price: 101, CreatedAt: 132429, Amount: 1},
		Order{Price: 103, CreatedAt: 132430, Amount: 1},
	}

	var AskOrderList = []Order{
		Order{Price: 100, CreatedAt: 132424, Amount: 1},
		Order{Price: 102, CreatedAt: 132425, Amount: 1},
		Order{Price: 101, CreatedAt: 132429, Amount: 1},
		Order{Price: 103, CreatedAt: 132428, Amount: 1},
		Order{Price: 103, CreatedAt: 132438, Amount: 1},
	}
	bk := Book{}

	for _, bid := range BidOrderList {
		bk.AddBid(bid)
	}

	for _, ask := range AskOrderList {
		bk.AddAsk(ask)
	}

	if bk.bidOrders[0].Price < bk.bidOrders[1].Price {
		t.Fatal("bid price not sorted")
	}

	if bk.askOrders[0].Price > bk.askOrders[1].Price {
		t.Fatal("ask price not sorted")
	}

	if bk.askOrders[3].CreatedAt < bk.askOrders[4].CreatedAt {
		t.Fatal("ask create time not sorted")
	}
}

func TestMatch(t *testing.T) {
	var BidOrderList = []Order{
		Order{Type: Bid, Price: 100, CreatedAt: 132424, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 102, CreatedAt: 132425, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 103, CreatedAt: 132428, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 101, CreatedAt: 132429, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 103, CreatedAt: 132430, Amount: 1, RestAmt: 1},
	}

	var AskOrderList = []Order{
		Order{Type: Ask, Price: 100, CreatedAt: 132424, Amount: 1, RestAmt: 1},
		Order{Type: Ask, Price: 102, CreatedAt: 132425, Amount: 1, RestAmt: 1},
		Order{Type: Ask, Price: 101, CreatedAt: 132429, Amount: 1, RestAmt: 1},
		Order{Type: Ask, Price: 103, CreatedAt: 132428, Amount: 1, RestAmt: 1},
		Order{Type: Ask, Price: 103, CreatedAt: 132438, Amount: 1, RestAmt: 1},
	}

	bk := Book{}
	for _, bid := range BidOrderList {
		bk.AddBid(bid)
	}

	for _, ask := range AskOrderList {
		bk.AddAsk(ask)
	}

	ods := bk.Match()
	// for _, od := range ods {
	// 	fmt.Printf("type:%v, price:%d, amount:%d\n", od.Type, od.Price, od.Amount)
	// }
	// fmt.Println("len(ods):", len(ods))
	assert.Equal(t, len(ods), 6)
}

// none match
func TestNoneMatch(t *testing.T) {
	var BidOrderList = []Order{
		Order{Type: Bid, Price: 101, CreatedAt: 132429, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 103, CreatedAt: 132430, Amount: 3, RestAmt: 3},
	}

	var AskOrderList = []Order{
		Order{Type: Ask, Price: 104, CreatedAt: 132438, Amount: 1, RestAmt: 1},
	}

	bk := Book{}
	for _, bid := range BidOrderList {
		bk.AddBid(bid)
	}

	for _, ask := range AskOrderList {
		bk.AddAsk(ask)
	}

	ods := bk.Match()
	assert.Equal(t, len(ods), 0)
}

// zero bid n asks match.
func TestMatchZero2N(t *testing.T) {
	var BidOrderList = []Order{
		Order{Type: Bid, Price: 100, CreatedAt: 132424, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 102, CreatedAt: 132425, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 103, CreatedAt: 132428, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 101, CreatedAt: 132429, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 103, CreatedAt: 132430, Amount: 7, RestAmt: 7},
	}

	var AskOrderList = []Order{
		Order{Type: Ask, Price: 100, CreatedAt: 132424, Amount: 1, RestAmt: 1},
		Order{Type: Ask, Price: 102, CreatedAt: 132425, Amount: 1, RestAmt: 1},
		Order{Type: Ask, Price: 101, CreatedAt: 132429, Amount: 1, RestAmt: 1},
		Order{Type: Ask, Price: 103, CreatedAt: 132428, Amount: 1, RestAmt: 1},
		Order{Type: Ask, Price: 103, CreatedAt: 132438, Amount: 1, RestAmt: 1},
	}

	bk := Book{}
	for _, bid := range BidOrderList {
		bk.AddBid(bid)
	}

	for _, ask := range AskOrderList {
		bk.AddAsk(ask)
	}

	ods := bk.Match()
	// for _, od := range ods {
	// 	fmt.Printf("type:%v, price:%d, amount:%d\n", od.Type, od.Price, od.Amount)
	// }
	assert.Equal(t, len(ods), 5)
}

// one bid match n asks.
func TestMatchOne2N(t *testing.T) {
	var BidOrderList = []Order{
		Order{Type: Bid, Price: 100, CreatedAt: 132424, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 102, CreatedAt: 132425, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 103, CreatedAt: 132428, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 101, CreatedAt: 132429, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 103, CreatedAt: 132430, Amount: 3, RestAmt: 3},
	}

	var AskOrderList = []Order{
		Order{Type: Ask, Price: 100, CreatedAt: 132424, Amount: 1, RestAmt: 1},
		Order{Type: Ask, Price: 102, CreatedAt: 132425, Amount: 1, RestAmt: 1},
		Order{Type: Ask, Price: 101, CreatedAt: 132429, Amount: 1, RestAmt: 1},
		Order{Type: Ask, Price: 103, CreatedAt: 132428, Amount: 1, RestAmt: 1},
		Order{Type: Ask, Price: 103, CreatedAt: 132438, Amount: 1, RestAmt: 1},
	}

	bk := Book{}
	for _, bid := range BidOrderList {
		bk.AddBid(bid)
	}

	for _, ask := range AskOrderList {
		bk.AddAsk(ask)
	}

	ods := bk.Match()
	// for _, od := range ods {
	// 	fmt.Printf("type:%v, price:%d, amount:%d\n", od.Type, od.Price, od.Amount)
	// }
	// fmt.Println("len(ods):", len(ods))
	assert.Equal(t, len(ods), 6)
}

// n bid match one asks.
func TestMatchN2One(t *testing.T) {
	var BidOrderList = []Order{
		Order{Type: Bid, Price: 100, CreatedAt: 132424, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 102, CreatedAt: 132425, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 103, CreatedAt: 132428, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 101, CreatedAt: 132429, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 103, CreatedAt: 132430, Amount: 1, RestAmt: 1},
	}

	var AskOrderList = []Order{
		Order{Type: Ask, Price: 100, CreatedAt: 132424, Amount: 4, RestAmt: 4},
		Order{Type: Ask, Price: 102, CreatedAt: 132425, Amount: 1, RestAmt: 1},
		Order{Type: Ask, Price: 101, CreatedAt: 132429, Amount: 1, RestAmt: 1},
		Order{Type: Ask, Price: 103, CreatedAt: 132428, Amount: 1, RestAmt: 1},
		Order{Type: Ask, Price: 103, CreatedAt: 132438, Amount: 1, RestAmt: 1},
	}

	bk := Book{}
	for _, bid := range BidOrderList {
		bk.AddBid(bid)
	}

	for _, ask := range AskOrderList {
		bk.AddAsk(ask)
	}

	ods := bk.Match()
	// for _, od := range ods {
	// 	fmt.Printf("type:%v, price:%d, amount:%d\n", od.Type, od.Price, od.Amount)
	// }
	// fmt.Println("len(ods):", len(ods))
	assert.Equal(t, len(ods), 5)
}

// n bid match n asks.
func TestMatchN2N(t *testing.T) {
	var BidOrderList = []Order{
		Order{Type: Bid, Price: 100, CreatedAt: 132424, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 102, CreatedAt: 132425, Amount: 2, RestAmt: 2},
		Order{Type: Bid, Price: 103, CreatedAt: 132428, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 101, CreatedAt: 132429, Amount: 1, RestAmt: 1},
		Order{Type: Bid, Price: 103, CreatedAt: 132430, Amount: 1, RestAmt: 1},
	}

	var AskOrderList = []Order{
		Order{Type: Ask, Price: 100, CreatedAt: 132424, Amount: 2, RestAmt: 2},
		Order{Type: Ask, Price: 102, CreatedAt: 132425, Amount: 1, RestAmt: 1},
		Order{Type: Ask, Price: 102, CreatedAt: 132440, Amount: 1, RestAmt: 1},
		Order{Type: Ask, Price: 101, CreatedAt: 132429, Amount: 1, RestAmt: 1},
		Order{Type: Ask, Price: 103, CreatedAt: 132428, Amount: 1, RestAmt: 1},
		Order{Type: Ask, Price: 103, CreatedAt: 132438, Amount: 1, RestAmt: 1},
	}

	bk := Book{}
	for _, bid := range BidOrderList {
		bk.AddBid(bid)
	}

	for _, ask := range AskOrderList {
		bk.AddAsk(ask)
	}

	ods := bk.Match()
	// for _, od := range ods {
	// 	fmt.Printf("type:%v, price:%d, amount:%d\n", od.Type, od.Price, od.Amount)
	// }
	// fmt.Println("len(ods):", len(ods))
	assert.Equal(t, len(ods), 6)
}

// zero bid and ask
func TestMatchZero(t *testing.T) {
	var BidOrderList = []Order{}

	var AskOrderList = []Order{}

	bk := Book{}
	for _, bid := range BidOrderList {
		bk.AddBid(bid)
	}

	for _, ask := range AskOrderList {
		bk.AddAsk(ask)
	}

	ods := bk.Match()
	// for _, od := range ods {
	// 	fmt.Printf("type:%v, price:%d, amount:%d\n", od.Type, od.Price, od.Amount)
	// }
	// fmt.Println("len(ods):", len(ods))
	assert.Equal(t, len(ods), 0)
}

func TestCopy(t *testing.T) {
	var BidOrderList = []Order{
		Order{Price: 100, CreatedAt: 132424, Amount: 1},
		Order{Price: 102, CreatedAt: 132425, Amount: 1},
		Order{Price: 103, CreatedAt: 132428, Amount: 1},
		Order{Price: 101, CreatedAt: 132429, Amount: 1},
		Order{Price: 103, CreatedAt: 132430, Amount: 1},
	}

	var AskOrderList = []Order{
		Order{Price: 100, CreatedAt: 132424, Amount: 1},
		Order{Price: 102, CreatedAt: 132425, Amount: 1},
		Order{Price: 101, CreatedAt: 132429, Amount: 1},
		Order{Price: 103, CreatedAt: 132428, Amount: 1},
		Order{Price: 103, CreatedAt: 132438, Amount: 1},
	}

	bk := Book{}
	for _, bid := range BidOrderList {
		bk.AddBid(bid)
	}

	for _, ask := range AskOrderList {
		bk.AddAsk(ask)
	}

	copyBk := bk.Copy()
	assert.Equal(t, len(bk.bidOrders), len(copyBk.bidOrders))
	for i := 0; i < len(bk.bidOrders); i++ {
		assert.NotEqual(t, fmt.Sprintf("%p", &bk.bidOrders[i]), fmt.Sprintf("%p", &copyBk.bidOrders[i]))
	}

	assert.Equal(t, len(bk.askOrders), len(copyBk.askOrders))
	for i := 0; i < len(bk.askOrders); i++ {
		assert.NotEqual(t, fmt.Sprintf("%p", &bk.askOrders[i]), fmt.Sprintf("%p", &copyBk.askOrders[i]))
	}
}
