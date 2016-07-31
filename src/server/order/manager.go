package order

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Manager struct {
	books map[string]*Book
	chans map[string]chan Order
}

func NewManager() *Manager {
	return &Manager{
		books: make(map[string]*Book),
		chans: make(map[string]chan Order),
	}
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
	return nil
}

// GetBook get specific coin pair's order book.
// the return book is an copy of internal book, for thread safe.
func (m *Manager) GetBook(coinPair string) Book {
	return m.books[coinPair].Copy()
}

func (m *Manager) RegisterOrderChan(coinPair string, c chan Order) {
	m.chans[coinPair] = c
}

// Run start the manager, tm is the match tick time, closing is the used to stop the manager from running.
func (m *Manager) Run(tm time.Duration, closing chan bool) {
	wg := sync.WaitGroup{}
	for p, bk := range m.books {
		wg.Add(1)
		go func(b *Book, orderChan chan Order, c chan bool, w *sync.WaitGroup) {
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
				}
			}
		}(bk, m.chans[p], closing, &wg)
	}
	wg.Wait()
}
