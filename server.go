package main

import (
	//"encoding/json"
	//"errors"
	"fmt"
	//"github.com/go-goodies/go_oops"
	//"github.com/l3x/jsoncfgo"
	"html/template"
	//"io/ioutil"
	"log"
	"net/http"
	//"regexp"
	"github.com/skycoin/skycoin/src/cipher"
	//"github.com/skycoin/skycoin/src/daemon/gnet"
	"net/http"
	"os"
)

//Context *gnet.MessageContext

/*
	The server gets events from the client and processes them
	- get balance/status
	- get deposit addresses
	- withdrawl bitcoin
	- withdrawl skycoin
	- add bid
	- add ask
	- get order book
*/

/*
Ping is response
Pong is response
*/

//get status
type PingStatus struct {
	Auth MsgAuth //must be included with every message

}

type PongStatus struct {
	Err string //error?

	Balance map[string]uint64

	DepositAddresses []string

	//list of orders unexecuted?
	//bid orders?
	//ask orders?
}

func (self *server) HandleStatus(in PingStatus, out *PongStatus) {
	addr := in.Auth.Address //user ID

	account := self.AccountStates.GetAccount(addr)
	if account != nil {
		out.Err = "account does not exist"
		return
	}

	out.Balance = account.Balance
	//out.SKY_balance = account.SKY_balance
}

//withdrawl coin

type PingWithdrawl struct {
	Auth     MsgAuth //must be included with every message
	CoinType string  //string for coin type
	Address  string  //address
}

type PongWithdrawl struct {
	Err string //error?

}

func (self *server) HandleWithdrawl(in PingStatus, out *PongStatus) {
	addr := in.Auth.Address //user ID
	//send to account manager
}

//raw event is authed event with a response handler
//is designed for playback for testing
type RawEvent struct {
	EventResponse chan string //response
	MsgAuth       MsgAuth
	Type          string
	Msg           []byte
}

//type event_handler func(string) error

//var handlers map[string]event_handler {
//	"test" : func(v string) {return nil},
//}

//handle events
func (self *Server) eventHandler(w http.ResponseWriter, r *http.Request) {

	auth := r.URL.Query()["MsgAuth"]
	msg_type := r.URL.Query()["msg_type"]
	msg := r.URL.Query()["msg"]

	value, ok := handlers[msg_type] // return value if found or ok=false if not found

	if !ok {
		HttpError(w, http.StatusBadRequest, "request type does not exist", nil)
	}

	_ = value

}

/*



*/

func (self *Server) Init() {

	host := "localhost"
	fmt.Printf("host: %v\n", host)

	port := int(6666)
	fmt.Printf("port: %v\n", port)
	addr := fmt.Sprintf("%s:%d", host, port)

	mux := http.NewServeMux()

	mux.Handle("/event", http.HandlerFunc(eventHandler))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "")
	})

	err := http.ListenAndServe(addr, mux)
	fmt.Println(err.Error())
}

/*
	Note:
	Everything must occur in a loop and be deterministic
	- must be syncrounous
	- should not allow asycronous requests
*/

//add a raw event
func (self *Server) InjectEvent(e RawEvent) {
	self.RawEventChannel <- e
}

type Server struct {
	AccountStates //anon, inherits methods/state

	RawEventChannel chan RawEvent
}

//read in raw event into processing loop
func (self *Server) RawEventToEvent(e) error {

}

//begin event processing loop
func (self *Server) Run(quit chan int) {

	//takes in raw events and push the onto channel
	go func() {
		for {
			select {

			case e <- RawEventIn:
				{
					err := CheckMsgAuth(e)
					if err == nil {

						e.EventResponse <- "error invalid message"
						break //ignore messages that are invalid
						//return error to event channel
					}
					err = self.RawEventToEvent(e) //push onto event loops

				}
			}
		}

		//main:

	}()
}

func (self *Server) Stop() {

}

//set of all users
type AccountManager struct {
	Accounts map[cipher.Address]*AccountState
	//AccountMap map[cipher.Address]uint64
}

//store state of user on server
type AccountState struct {

	//AccountId  uint64
	Address cipher.Address //Account id

	Balance map[string]uint64
	//Bitcoin balance in satoshis
	//Skycoin balance in drops

	//Inc1 uint64 //inc every write? Associated with local change
	//Inc2 uint64 //set to last change. Associatd with global event id
}

func (self *AccountState) GetAccount(addr cipher.Address) (*AccountState, error) {
	account, ok := self.Accounts[addr]
	if !ok {
		return nil, errors.New("Account does not exist")
	}
	return account
}

func (self *AccountState) CreateAccount(addr cipher.Address) error {
	if _, ok := self.Accounts[addr]; ok == true {
		return errors.New("Account already exists")
	}

	self.Accounts[addr] = AccountState{Address: addr}
	self.Accounts.Balance = map[string]uint64{
		"BTC": 0,
		"SKY": 0,
	}
}

//persistance to disc. Save as JSON
func (self *AccountManager) Save() {

}

func (self *AccountManager) Load() {

	//load accounts
}

/*
Events
- bitcoin deposit event
- skycoin deposit event

Events
-

*/
