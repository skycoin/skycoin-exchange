package bitcoin_interface

import (
	//"errors"

	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"log"
	"net/http"

	"github.com/skycoin/skycoin-exchange/src/server/coin_interface"
	"github.com/skycoin/skycoin/src/cipher"
)

var (
	HideSeckey = true
)

// Utxo unspent output
type Utxo interface {
	GetTxid() string
	GetVout() uint32
	GetAmount() uint64
	GetAddress() string
}

// UtxoWithkey unspent output with privkey.
type UtxoWithkey interface {
	Utxo
	GetPrivKey() string
}

type UtxoOut struct {
	Addr  string
	Value uint64
}

// GenerateAddresses, generate bitcoin addresses.
func GenerateAddresses(seed []byte, num int) (string, []coin_interface.AddressEntry) {
	sd, seckeys := cipher.GenerateDeterministicKeyPairsSeed(seed, num)
	entries := make([]coin_interface.AddressEntry, num)
	for i, sec := range seckeys {
		pub := cipher.PubKeyFromSecKey(sec)
		entries[i].Address = cipher.BitcoinAddressFromPubkey(pub)
		entries[i].Public = pub.Hex()
		if !HideSeckey {
			entries[i].Secret = cipher.BitcoinWalletImportFormatFromSeckey(sec)
		}
	}
	return fmt.Sprintf("%2x", sd), entries
}

// GetBalance, query balance of address through the API of blockexplorer.com.
func GetBalance(addr []string) (uint64, error) {
	// if AddressValid(addr) != nil {
	// 	log.Fatal("Address is invalid")
	// }
	// blkEplUrl := fmt.Sprintf("https://blockexplorer.com/api/addr/%s/balance", addr)
	addrs := strings.Join(addr, "|")
	blkChnUrl := fmt.Sprintf("https://blockchain.info/q/addressbalance/%s", addrs)
	data, err := getDataOfUrl(blkChnUrl)
	if err != nil {
		return 0, err
	}
	b, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}
	return uint64(b), nil
}

// GetUnspentOutputs return the unspent outputs
func GetUnspentOutputs(addrs []string) ([]Utxo, error) {
	return getUtxosBlkExplr(addrs)
}

func NewUtxoWithKey(utxo Utxo, key string) UtxoWithkey {
	return BlkExplrUtxoWithkey{
		BlkExplrUtxo: utxo.(BlkExplrUtxo),
		Privkey:      key,
	}
}

type Manager struct {
	WatchAddresses []string
	UxStateMap     map[string]Utxo //keeps track of state
}

type UxMap map[string]Utxo

//does querry/update
func (self *Manager) UpdateOutputs() {
	log.Println("Update outputs...")
	//get all unspent outputs for all watch addresses
	var list []Utxo
	for _, addr := range self.WatchAddresses {
		ux, err := GetUnspentOutputs([]string{addr})
		if err != nil {
			panic(err)
		}
		list = append(list, ux...)
	}
	latestUxMap := make(map[string]Utxo)
	//do diff
	for _, utxo := range list {
		id := fmt.Sprintf("%s:%d", utxo.GetTxid(), utxo.GetVout())
		latestUxMap[id] = utxo
	}

	//get new
	NewUx := make(map[string]Utxo)
	for id, utxo := range latestUxMap {
		if _, ok := self.UxStateMap[id]; !ok {
			NewUx[id] = utxo
			log.Printf("New output Found:%+v\n", utxo)
		}
	}

	// TODO:
	// make sure outputs that exist, never disappear, without being spent
	// means theft or blockchain fork

	// look for ux that disappeared
	// TODO: make sure output exists and has not disappeared, else panic mode
	// TODO: output should still exist, even if not spendable
	DisappearingUx := make(map[string]Utxo)
	for id, utxo := range self.UxStateMap {
		if _, ok := self.UxStateMap[id]; !ok {
			DisappearingUx[id] = utxo
			log.Printf("Output Disappered: %+v\n", utxo)
		}
	}

	self.UxStateMap = latestUxMap
}

// getDataOfUrl, get data from specific URL.
func getDataOfUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	resp.Body.Close()
	return data, nil
}

func (self *Manager) Init() {
	//UxStateMap     map[string]UnspentOutputJSON
	self.WatchAddresses = make([]string, 0)
	self.UxStateMap = make(map[string]Utxo)
}

func (self *Manager) AddWatchAddress(addr string) {
	if AddressValid(addr) != nil {
		log.Fatal("Address being added to watch list, must be valid")
	}
	self.WatchAddresses = append(self.WatchAddresses, addr)
}

func (self *Manager) Tick() {
	self.UpdateOutputs()
}

//returns error if the address is invalid
func AddressValid(address string) error {
	//return errors.New("Address is invalid")
	return nil
}
