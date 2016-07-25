package skycoin_interface

import (
	"errors"
	"math/rand"
	"time"

	"github.com/golang/glog"
)

var CheckTick = 5 * time.Second

type UtxoManager interface {
	Start(closing chan bool)
	ChooseUtxos(amt uint64, tm time.Duration) ([]Utxo, error)
	// GetUtxo() chan Utxo // get utxo from utxo pool
	PutUtxo(utxo Utxo) // put utxo into utxo pool
	WatchAddresses(addrs []string)
}

type ExUtxoManager struct {
	WatchAddress []string
	UtxosCh      chan Utxo
	UtxoStateMap map[string]Utxo
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
	glog.Info("start the skycoin utxo manager")
	t := time.Tick(CheckTick)
	for {
		select {
		case <-closing:
			return
		case <-t:
			// check skycoin new utxos.
			newUtxos, err := eum.checkNewUtxo()
			if err != nil {
				glog.Error(err)
				break
			}

			for _, utxo := range newUtxos {
				glog.Info("new skycoin utxo:", utxo.GetHash(), " ", utxo.GetCoins(), " ", utxo.GetHours())
				eum.UtxosCh <- utxo
			}
		}
	}
}

func (eum *ExUtxoManager) GetUtxo() chan Utxo {
	return eum.UtxosCh
}

func (eum *ExUtxoManager) PutUtxo(utxo Utxo) {
	glog.Info("skycoin utxo put back:", utxo.GetHash())
	eum.UtxosCh <- utxo
}

func (eum *ExUtxoManager) WatchAddresses(addrs []string) {
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
	newUtxos := []Utxo{}
	for id, utxo := range latestUxMap {
		if _, ok := eum.UtxoStateMap[id]; !ok {
			newUtxos = append(newUtxos, utxo)
		}
	}

	eum.UtxoStateMap = latestUxMap
	return newUtxos, nil
}

// chooseUtxos choose appropriate utxos, if time out, and not found enough utxos,
// the utxos got before will put back to the utxos pool, and return error.
// the tm is millisecond
func (eum *ExUtxoManager) chooseUtxos(amount uint64, tm time.Duration) ([]Utxo, error) {
	glog.Info("skycoin choose utxos, amount:", amount)
	var totalAmount uint64
	// utxos := []bitcoin.UtxoWithkey{}
	utxos := []Utxo{}
	for {
		select {
		case utxo := <-eum.UtxosCh:
			glog.Info("get utxo:", utxo.GetHash(), " ", utxo.GetCoins())
			utxos = append(utxos, utxo)
			totalAmount += utxo.GetCoins()
			if totalAmount >= amount {
				return utxos, nil
			}

		case <-time.After(tm):
			// put utxos back
			glog.Info("choose time out, put back utxos")
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
					// return []Utxo{}, err
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
