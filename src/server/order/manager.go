package order

import (
	"errors"
	"fmt"
	"time"
)

type Manager struct {
	books     map[string]*Book
	chans     map[string]chan Order
	matchTick time.Duration
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

func (m *Manager) Run(closing chan bool) {
	for p, bk := range m.books {
		go func(b *Book, orderChan chan Order, c chan bool) {
			orders := []Order{}
			for {
				select {
				case <-c:
					return
				case <-time.After(m.matchTick):
					orders = b.Match()
					for _, o := range orders {
						orderChan <- o
					}
				}
			}
		}(bk, m.chans[p], closing)
	}
}
