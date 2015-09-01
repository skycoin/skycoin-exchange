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
	"os"
	"net/http"
)


//Context *gnet.MessageContext


/*
	The server gets events from the client and processes them
	- get balance
	- get deposit addresses
	- withdrawl bitcoin
	- etc


*/

//raw event is authed event with a response handler
//is designed for playback for testing
type RawEvent struct {
	EventResponse chan string //response
	MsgAuth MsgAuth
	Type string
	Msg []byte
}

type event_handler func(string) error

var handlers map[string]handle {
	"test" : func(v string) {return nil},
}



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
	UserStates //anon, inherits methods/state

	RawEventChannel chan RawEvent
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
				err := self.RawEventToEvent(e) //push onto event loops

			}
		}
	}

//main:

}




//read in raw event into processing loop
func (self *Server) RawEventToEvent(e) error {

}


func (self *Server) Stop() {

}

//set of all users
type UserStates struct {
	Users map[uint64]UserState
	DepositAddresses map[cipher.Address]UserId
}

//store state of user on server
type UserState struct {
	UserId  uint64
	Address cipher.Address //user id

	BitcoinBalance uint64 //Bitcoin balance in satoshis
	SkycoinBalance uint64 //Skycoin balance in drops
}


type BitcoinDepositAddress struct {
	BitcoinAddress string
	Owner          UserId
	Created time.Time
}

/*
Events
- bitcoin deposit event
- skycoin deposit event

Events
- 

*/