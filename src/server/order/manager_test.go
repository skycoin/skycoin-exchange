package order

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/util/file"
	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	m := NewManager()
	coinPair := "btc/sky"
	m.AddBook(coinPair, &Book{})
	btcSkyChan := make(chan Order, 100)
	m.RegisterOrderChan(coinPair, btcSkyChan)
	closing := make(chan bool)
	go m.Start(time.Duration(1)*time.Second, closing)

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

	for _, od := range BidOrderList {
		m.AddOrder(coinPair, od)
	}

	for _, od := range AskOrderList {
		m.AddOrder(coinPair, od)
	}

	totalMath := 0
	go func(orders chan Order, c chan bool) {
		for {
			select {
			case od := <-orders:
				// assert.Equal(t, od.RestAmt, 0)
				if od.RestAmt != 0 {
					t.Fatal("match order's reset amt is not zero")
				}
				totalMath += 1
				// fmt.Printf("match order: type:%v, price:%d, amount:%d, restamt:%d\n", od.Type, od.Price, od.Amount, od.RestAmt)
			case <-c:
				return
			}
		}
	}(btcSkyChan, closing)
	time.Sleep(2 * time.Second)
	// fmt.Println("add new bid: price 104")
	m.AddOrder(coinPair, Order{Type: Bid, Price: 104, Amount: 1, RestAmt: 1})
	time.Sleep(2 * time.Second)
	close(closing)
	assert.Equal(t, totalMath, 8)
}

func TestLoadManager(t *testing.T) {
	// prepare data
	coinPair := []string{"test", "sky"}
	bk := Book{}
	var BidOrderList = []Order{
		Order{ID: 1, Type: Bid, Price: 100, CreatedAt: 132424, Amount: 1, RestAmt: 1},
		Order{ID: 2, Type: Bid, Price: 102, CreatedAt: 132425, Amount: 1, RestAmt: 1},
		Order{ID: 3, Type: Bid, Price: 103, CreatedAt: 132428, Amount: 1, RestAmt: 1},
		Order{ID: 4, Type: Bid, Price: 101, CreatedAt: 132429, Amount: 1, RestAmt: 1},
		Order{ID: 5, Type: Bid, Price: 103, CreatedAt: 132430, Amount: 1, RestAmt: 1},
	}

	var AskOrderList = []Order{
		Order{ID: 6, Type: Ask, Price: 100, CreatedAt: 132424, Amount: 1, RestAmt: 1},
		Order{ID: 7, Type: Ask, Price: 102, CreatedAt: 132425, Amount: 1, RestAmt: 1},
		Order{ID: 8, Type: Ask, Price: 101, CreatedAt: 132429, Amount: 1, RestAmt: 1},
		Order{ID: 9, Type: Ask, Price: 103, CreatedAt: 132428, Amount: 1, RestAmt: 1},
		Order{ID: 10, Type: Ask, Price: 103, CreatedAt: 132438, Amount: 1, RestAmt: 1},
	}

	for _, od := range BidOrderList {
		bk.AddBid(od)
	}

	for _, od := range AskOrderList {
		bk.AddAsk(od)
	}

	// write book to files.
	err := file.SaveJSON(filepath.Join(orderDir, strings.Join(coinPair, "_")+"."+orderExt), bk.ToMarshalable(), 0600)
	assert.Nil(t, err)
	m, err := LoadManager()
	assert.Nil(t, err)
	bk1 := m.GetBook(strings.Join(coinPair, "/"))
	assert.Equal(t, bk, bk1)
}
