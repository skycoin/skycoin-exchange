package bitcoin_interface

import (
	//"errors"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/skycoin/skycoin/src/cipher"
	//"github.com/skycoin/skycoin/src/cipher"
	"net/http"
	//"sort"
	//"strings"
	//"time"
	//"bytes"
	"log"
)

/*



 */

/*
{
    "unspent_outputs":[
        {
            "tx_age":"1322659106",
            "tx_hash":"e6452a2cb71aa864aaa959e647e7a4726a22e640560f199f79b56b5502114c37",
            "tx_index":"12790219",
            "tx_output_n":"0",
            "script":"76a914641ad5051edd97029a003fe9efb29359fcee409d88ac", (Hex encoded)
            "value":"5000661330"
        }
    ]
}
*/

var (
	HideSeckey = true
)

//returns nil, if address is valid
//returns error if the address is invalid
func AddressValid(address string) error {
	//return errors.New("Address is invalid")
	return nil
}

type UnspentOutputJSONResponse struct {
	UnspentOutputArray []UnspentOutputJSON `json:"unspent_outputs"`
}

type UnspentOutputJSON struct {
	// Tx_age      uint64 `json:"tx_age"`
	Tx_hash            string `json:"tx_hash"` // the previous transaction id
	Tx_hash_big_endian string `json:"tx_hash_big_endian"`
	Tx_index           uint64 `json:"tx_index"`
	Tx_output_n        uint64 `json:"tx_output_n"` // the output index of previous transaction
	Script             string `json:"script"`      // pubkey script
	Value              uint64 `json:"value"`       // the bitcoin amount in satoshis
	Value_hex          string `json:"value_hex"`   // alisa the Value, in hex format.
	Confirmations      uint64 `json:"confirmations"`
}

type AddressEntry struct {
	Address string
	Public  string
	Secret  string
}

// ByAge implements sort.Interface for []Person based on
// the Age field.

/*
type ByHash []Person

func (a ByHash) Len() int           { return len(a) }
func (a ByHash) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHash) Less(i, j int) bool { return a[i].Age < a[j].Age }
*/

// GenerateAddresses, generate bitcoin addresses.
func GenerateAddresses(seed []byte, num int) (string, []AddressEntry) {
	sd, seckeys := cipher.GenerateDeterministicKeyPairsSeed(seed, num)
	entries := make([]AddressEntry, num)
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
func GetBalance(addr string) (string, error) {
	if AddressValid(addr) != nil {
		log.Fatal("Address is invalid")
	}

	data, err := getDataOfUrl(fmt.Sprintf("https://blockexplorer.com/api/addr/%s/balance", addr))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// https://blockchain.info/unspent?active=1SakrZuzQmGwn7MSiJj5awqJZjSYeBWC3
// GetUnspentOutputs, using the API from blockchain.info to query the unspent outputs.
func GetUnspentOutputs(addr string) []UnspentOutputJSON {
	if AddressValid(addr) != nil {
		log.Fatal("Address is invalid")
	}

	//url := strings.Sprint
	// fmt.Printf("Address= %s\n", addr)
	resp, err := http.Get(fmt.Sprintf("https://blockchain.info/unspent?active=%s", addr))
	if err != nil {
		log.Fatalf("Get url:%s fail, error:%s", addr, err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Read data from resp body fail, error:%s", err)
	}
	resp.Body.Close()
	// fmt.Println("data:", string(data))

	// parse the JSON.
	utxoResp := UnspentOutputJSONResponse{}
	err = json.Unmarshal(data, &utxoResp)
	if err != nil {
		log.Fatalf("unmasharl fail, error:%s", err)
	}

	return utxoResp.UnspentOutputArray
}

type Manager struct {
	WatchAddresses []string
	UxStateMap     map[string]UnspentOutputJSON //keeps track of state
}

type UxMap map[string]UnspentOutputJSON

//does querry/update
func (self *Manager) UpdateOutputs() {
	log.Println("Update outputs...")
	//get all unspent outputs for all watch addresses
	var list []UnspentOutputJSON
	for _, addr := range self.WatchAddresses {
		ux := GetUnspentOutputs(addr)
		list = append(list, ux...)
	}
	latestUxMap := make(map[string]UnspentOutputJSON)
	//do diff
	for _, utxo := range list {
		id := fmt.Sprintf("%s:%d", utxo.Tx_hash, utxo.Tx_index)
		latestUxMap[id] = utxo
	}

	//get new
	NewUx := make(map[string]UnspentOutputJSON)
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
	DisappearingUx := make(map[string]UnspentOutputJSON)
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
	self.UxStateMap = make(map[string]UnspentOutputJSON)
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
