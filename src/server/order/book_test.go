package order

import (
	"fmt"
	"testing"
)

var BidOrderList = []Order{
	Order{Price: 100, CreatedTime: 132424, Amount: 1},
	Order{Price: 102, CreatedTime: 132425, Amount: 1},
	Order{Price: 103, CreatedTime: 132428, Amount: 1},
	Order{Price: 101, CreatedTime: 132429, Amount: 1},
}

var AskOrderList = []Order{
	Order{Price: 100, CreatedTime: 132424, Amount: 1},
	Order{Price: 102, CreatedTime: 132425, Amount: 1},
	Order{Price: 101, CreatedTime: 132429, Amount: 1},
	Order{Price: 103, CreatedTime: 132428, Amount: 1},
	Order{Price: 103, CreatedTime: 132438, Amount: 1},
}

func TestBook(t *testing.T) {
	bk := Book{
		askOrders: AskOrderList,
		bidOrders: BidOrderList,
	}

	bk.Sort()
	if bk.bidOrders[0].Price < bk.bidOrders[1].Price {
		t.Fail()
	}

	fmt.Println("asks")
	if bk.askOrders[0].Price != bk.askOrders[1].Price {
		t.Fail()
	}

	if bk.askOrders[0].CreatedTime < bk.askOrders[1].CreatedTime {
		t.Fail()
	}
}
