package skycoin_exchange

import (
	//"encoding/json"
	//"errors"

	"fmt"
	//"github.com/go-goodies/go_oops"
	//"github.com/l3x/jsoncfgo"

	//"io/ioutil"

	"net/http"
	//"regexp"
	//"github.com/skycoin/skycoin/src/daemon/gnet"
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

func (self *Server) HandleStatus(in PingStatus, out *PongStatus) {
	addr := in.Auth.Address //user ID
	// get account id from the address.

	account, err := self.AccountManager.GetAccount(AccountID(addr))
	if err != nil {
		out.Err = err.Error()
		return
	}
	balanceMap := account.GetBalanceMap()
	for ctp, bal := range balanceMap {
		out.Balance[ctp.String()] = uint64(bal)
	}
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

func (self *Server) HandleWithdrawl(in PingStatus, out *PongStatus) {
	// addr := in.Auth.Address //user ID
	//send to account manager
}

//Add bid/ask

type PingOrder struct {
	Auth      MsgAuth
	OrderType string //"bid" or "ask"
	CoinPair  string //"SKY/BTC"
	price     uint64
	quantity  uint64
}

type PongOrder struct {
	Err     string //error?
	OrderId uint64
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

	// auth := r.URL.Query()["MsgAuth"]
	// msg_type := r.URL.Query()["msg_type"]
	// msg := r.URL.Query()["msg"]

	// value, ok := handlers[msg_type] // return value if found or ok=false if not found
	//
	// if !ok {
	// 	HttpError(w, http.StatusBadRequest, "request type does not exist", nil)
	// }
	//
	// _ = value

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
	AccountManager  //anon, inherits methods/state
	RawEventChannel chan RawEvent
}

//read in raw event into processing loop
func (self *Server) RawEventToEvent(e RawEvent) error {
	return nil
}

//begin event processing loop
func (self *Server) Run(quit chan int) {
	//takes in raw events and push the onto channel
	go func() {
		for {
			select {
			case e := <-self.RawEventChannel:
				{
					err := CheckMsgAuth(e.MsgAuth)
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

/*
Events
- bitcoin deposit event
- skycoin deposit event

Events
-

*/
