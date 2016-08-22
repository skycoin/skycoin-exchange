package skycoin_interface

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var CheckTick = 5 * time.Second

type UtxoManager interface {
	Start(closing chan bool)
	ChooseUtxos(amt uint64, tm time.Duration) ([]Utxo, error)
	PutUtxo(utxo Utxo) // put utxo into utxo pool
	WatchAddresses(addrs []string)
}

type ExUtxoManager struct {
	WatchAddress []string
	UtxosCh      chan Utxo
	UtxoStateMap map[string]Utxo
	mutx         sync.Mutex
}

func NewUtxoManager(utxoPoolsize int, watchAddrs []string) UtxoManager {
	eum := &ExUtxoManager{
		UtxosCh:      make(chan Utxo, utxoPoolsize),
		UtxoStateMap: make(map[string]Utxo),
		WatchAddress: watchAddrs,
	}

	return eum
}

func (eum *ExUtxoManager) Start(closing chan bool) {
	logger.Info("start the skycoin utxo manager")
	t := time.Tick(CheckTick)
	for {
		select {
		case <-closing:
			return
		case <-t:
			// check skycoin new utxos.
			newUtxos, err := eum.checkNewUtxo()
			if err != nil {
				logger.Error(err.Error())
				break
			}

			for _, utxo := range newUtxos {
				logger.Debug("new skycoin utxo: hash:%s coins:%d hours:%d",
					utxo.GetHash(), utxo.GetCoins(), utxo.GetHours())
				eum.UtxosCh <- utxo
			}
		}
	}
}

func (eum *ExUtxoManager) PutUtxo(utxo Utxo) {
	logger.Debug("skycoin utxo put back: %s", utxo.GetHash())
	eum.UtxosCh <- utxo
}

func (eum *ExUtxoManager) WatchAddresses(addrs []string) {
	for _, addr := range addrs {
		logger.Debug("skycoin watch address:%s", addr)
	}
	eum.WatchAddress = append(eum.WatchAddress, addrs...)
}

func (eum *ExUtxoManager) checkNewUtxo() ([]Utxo, error) {
	latestUtxos, err := GetUnspentOutputs(eum.WatchAddress)
	if err != nil {
		return []Utxo{}, err
	}

	latestUxMap := make(map[string]Utxo)
	// do diff
	for _, utxo := range latestUtxos {
		id := utxo.GetHash()
		latestUxMap[id] = utxo
	}

	//get new
	eum.mutx.Lock()
	newUtxos := []Utxo{}
	for id, utxo := range latestUxMap {
		if _, ok := eum.UtxoStateMap[id]; !ok {
			newUtxos = append(newUtxos, utxo)
		}
	}

	eum.UtxoStateMap = latestUxMap
	eum.mutx.Unlock()
	return newUtxos, nil
}

func (eum *ExUtxoManager) mustGetUtxos(hash string) Utxo {
	eum.mutx.Lock()
	defer eum.mutx.Unlock()
	if u, ok := eum.UtxoStateMap[hash]; ok {
		return u
	}
	panic(fmt.Sprintf("utxo:%s not found", hash))
}

// chooseUtxos choose appropriate utxos, if time out, and not found enough utxos,
// the utxos got before will put back to the utxos pool, and return error.
// the tm is millisecond
func (eum *ExUtxoManager) chooseUtxos(amount uint64, tm time.Duration) ([]Utxo, error) {
	logger.Debug("skycoin choose utxos, amount:%d", amount)
	var totalAmount uint64
	utxos := []Utxo{}
	for {
		select {
		case utxo := <-eum.UtxosCh:
			u := eum.mustGetUtxos(utxo.GetHash())
			if u.GetCoins() != utxo.GetCoins() {
				panic("utxo coins not equal")
			}
			logger.Debug("get utxo: hash:%s coins:%d hours:%d",
				utxo.GetHash(), utxo.GetCoins(), utxo.GetHours())
			utxos = append(utxos, u)
			totalAmount += utxo.GetCoins() * 1e6
			if totalAmount >= amount {
				return utxos, nil
			}

		case <-time.After(tm):
			// put utxos back
			logger.Debug("choose time out, put back utxos")
			for _, u := range utxos {
				eum.UtxosCh <- u
			}
			return []Utxo{}, nil
		}
	}
}

// ChooseUtxos choose sufficient utxos in specific time.
func (eum *ExUtxoManager) ChooseUtxos(amt uint64, tm time.Duration) ([]Utxo, error) {
	var (
		utxos []Utxo
		err   error
		ch    = make(chan bool)
		ok    = make(chan bool)
	)

	go func(closing chan bool, ok chan bool) {
		for {
			select {
			case <-closing:
				for _, u := range utxos {
					eum.UtxosCh <- u
				}
				return
			default:
				utxos, err = eum.chooseUtxos(amt, randExpireTm())
				if err != nil {
					return
				}

				if len(utxos) > 0 {
					ok <- true
					return
				}
			}
		}
	}(ch, ok)

	for {
		select {
		case <-time.After(tm):
			ch <- true
			return []Utxo{}, errors.New("time out")
		case <-ok:
			return utxos, nil
		}
	}
}

func randExpireTm() time.Duration {
	v := rand.Intn(5)
	return time.Duration(3+v) * time.Second
}
