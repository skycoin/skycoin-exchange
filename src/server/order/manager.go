package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/util/file"
)

type Manager struct {
	books map[string]*Book
	chans map[string]chan Order
	idg   map[string]*IDGenerator
}

func NewManager() *Manager {
	return &Manager{
		books: make(map[string]*Book),
		chans: make(map[string]chan Order),
		idg:   make(map[string]*IDGenerator),
	}
}

func LoadManager() (*Manager, error) {
	// check if the order dir exists
	if _, err := os.Stat(orderDir); os.IsNotExist(err) {
		return nil, err
	}

	files, err := ioutil.ReadDir(orderDir)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, os.ErrNotExist
	}

	m := NewManager()
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), orderExt) {
			continue
		}
		d, err := ioutil.ReadFile(filepath.Join(orderDir, f.Name()))
		if err != nil {
			return nil, err
		}
		bj := BookJson{}
		if err := json.Unmarshal(d, &bj); err != nil {
			return nil, err
		}
		p := strings.Split(f.Name(), ".")
		pair := strings.Split(p[0], "_")
		if len(pair) != 2 {
			panic("error order book file name")
		}
		cp := strings.Join(pair, "/")
		m.books[cp] = NewBookFromJson(bj)

		// init order id generator.
		m.idg[cp] = newIDGenerator(cp)
	}

	return m, nil
}

// AddBook add the order book of specific coin pair to manager,
// the stored book is an copy book, for thread safe.
func (m *Manager) AddBook(coinPair string, book *Book) error {
	if coinPair == "" {
		return errors.New("coin pair is empty")
	}

	if _, ok := m.books[coinPair]; ok {
		return fmt.Errorf("book of coin pair: %s already exists", coinPair)
	}
	bk := book.Copy()
	m.books[coinPair] = &bk

	m.idg[coinPair] = newIDGenerator(coinPair)
	return nil
}

func (m *Manager) IsExist(coinPair string) bool {
	if _, ok := m.books[coinPair]; ok {
		return true
	}
	return false
}

// AddOrder add bid or ask order to order book.
func (m *Manager) AddOrder(coinPair string, order Order) (uint64, error) {
	bk, ok := m.books[coinPair]
	if !ok {
		return 0, fmt.Errorf("coin pair:%s not supported", coinPair)
	}

	idg, ok := m.idg[coinPair]
	if !ok {
		return 0, fmt.Errorf("coin pair:%s's id generator not supported", coinPair)
	}

	switch order.Type {
	case Bid:
		order.ID = idg.GetID()
		bk.AddBid(order)
		return order.ID, nil
	case Ask:
		order.ID = idg.GetID()
		bk.AddAsk(order)
		return order.ID, nil
	default:
		return 0, errors.New("unknow order type")
	}
}

// GetBook get specific coin pair's order book.
// the return book is an copy of internal book, for thread safe.
func (m *Manager) GetBook(coinPair string) Book {
	return m.books[coinPair].Copy()
}

func (m *Manager) GetOrders(cp string, tp Type, start, end int64) ([]Order, error) {
	if _, ok := m.books[cp]; !ok {
		return []Order{}, errors.New("get orders faile, err: unknow coin pair")
	}
	return m.books[cp].GetOrders(tp, start, end), nil
}

func (m *Manager) RegisterOrderChan(coinPair string, c chan Order) {
	m.chans[coinPair] = c
}

// Run start the manager, tm is the match tick time, closing is used for stopping the manager from running.
func (m *Manager) Start(tm time.Duration, closing chan bool) {
	// start the id generators
	for _, g := range m.idg {
		go g.Run(closing)
	}

	// start the match timer.
	wg := sync.WaitGroup{}
	for p, bk := range m.books {
		wg.Add(1)
		go func(cp string, b *Book, orderChan chan Order, c chan bool, w *sync.WaitGroup) {
			orders := []Order{}
			for {
				select {
				case <-c:
					w.Done()
					return
				case <-time.After(tm):
					orders = b.Match()
					for _, o := range orders {
						orderChan <- o
					}
					// update order book in local disk.
					pairs := strings.Split(cp, "/")
					if len(pairs) != 2 {
						panic("error coin pair name")
					}
					filename := strings.Join(pairs, "_")
					if err := file.SaveJSON(filepath.Join(orderDir, filename+"."+orderExt), b.Copy().ToMarshalable(), 0600); err != nil {
						panic(err)
					}
				}
			}
		}(p, bk, m.chans[p], closing, &wg)
	}
	wg.Wait()
}
