package bitcoin_interface

import (
	"errors"
	"fmt"
	"github.com/skycoin/skycoin/src/cipher"
	"http"
	"sort"
	"strings"
	"time"
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

//returns nil, if address is valid
//returns error if the address is invalid
func AddressValid(address string) error {
	//return errors.New("Address is invalid")
	return nil
}

type UnspentOutput1JSONResponse struct {
	UnspentOutputArray []UnspentOutputJSON `json:"unspent_outputs"`
}

type UnspentOutputJSON struct {
	Tx_age      uint64 `json:"tx_age"`
	Tx_hash     string `json:"tx_hash"`
	Tx_index    string `json:"tx_index"`
	Tx_output_n uint64 `json:"tx_output_n"`
	Script      string `json:"script"`
	Value       uint64 `json:"value"`
}

// ByAge implements sort.Interface for []Person based on
// the Age field.

/*
type ByHash []Person

func (a ByHash) Len() int           { return len(a) }
func (a ByHash) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHash) Less(i, j int) bool { return a[i].Age < a[j].Age }
*/

//https://blockchain.info/unspent?active=1SakrZuzQmGwn7MSiJj5awqJZjSYeBWC3

func GetUnspentOutputs(addr string) {

	if AddressValid(addr) != nil {
		log.Fatal("Address is invalid")
	}

	url := strings.Sprintf
	req := fmt.Sprint("https://blockchain.info/unspent?active=%s", addr)

	//reader := strings.NewReader(`{"active":123}`)
	request, err := http.NewRequest("GET", req)
	// TODO: check err
	client := &http.Client{}
	resp, err := client.Do(request)
	// TODO: check err

	fmt.Printf()
}

type Manager struct {
	WatchAddresses []string
	UxStateMap     map[string]UnspentOutputJSON //keeps track of state
}

type UxMap map[string]UnspentOutputJSON

//does querry/update
func (self *Manager) UpdateOutputs() {

	//get all unspent outputs for all watch addresses
	var list []UnspentOutputJSONResponse
	for addr := range self.WatchAddresses {
		list = append(list, GetUnspentOutputs(addr))
	}
	var uxMap map[string]UnspentOutputJSON

	//do diff
	for _, j := range list {
		id = fmt.Sprint("%s:%i", j.Tx_hash, j.Tx_index)
		fmt.Printf("ID = id\n")
		uxMap[id] = j
	}

	//get new
	var NewUx map[string]UnspentOutputJSON

	//check existing state, compare to new state
	for i, j := range self.UxStateMap {
		val, ok := uxMap[i]
		if !ok {
			//new unspent output found
			NewUx[i] = j
			log.Printf("New Output Found: %x", j)
		}
	}

	// TODO:
	// make sure outputs that exist, never disappear, without being spent
	// means theft or blockchain fork

	//look for ux that disappeared
	// TODO: make sure output exists and has not disappeared, else panic mode
	// TODO: output should still exist, even if not spendable
	var DisappearingUx map[string]UnspentOutputJSON

	for i, j := range UxMap {
		val, ok := self.UxStateMap[i]
		if !ok {
			//output disappeared
			//means it has been spent
			//ensure output never disappears!!!
			DisappearingUx[i] = self.UxStateMap[i]
			log.Printf("Output Disappeared: %x", self.UxStateMap[i])
		}
	}

	self.UxStateMap = UxMap
}

func (self *Manager) Tick() {
	UpdateOutputs()
}
