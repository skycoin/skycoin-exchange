package skycoin_exchange

import (
	"github.com/skycoin/skycoin/src/cipher"
	"time"
)

/*
Manages withdrawls and deposits
- associates external addresses with internal account IDs
- deposits and withdrawls
- credits and debits internal accounts

Keep a list of deposit addresses for each coin, for each account
*/

/*
- keep a list of addresses for accounts
- keep a list of transactions/outputs for those addresses
-- make sure to only credit user once and record deposits


*/

/*
	Handle outgoing transactions
*/
type PendingOutgoing struct {
	Coin    string //BTC, SKY
	Address string //cipher.Address
	Amount  uint64 //Satoshis or Drops
}

type CompletedOutgoing struct {
	Coin    string //BTC, SKY
	Address string //cipher.Address
	Amount  uint64 //Satoshis or Drops
	Tx      string //transaction ID
}

type OutgoingManager struct {
	PendingOutgoingTransactions   []PendingOutgoing
	CompletedOutgoingTransactions []CompletedOutgoing
}

/*
	Handle incoming, deposit addresses and received coins
*/

type DepositAddress struct {
	AccountId cipher.Address //Account id
	Coin      string         //BTC,SKY
	Address   string         //address for deposit
	Tx        string         //transaction ID
}

type PendingReceivedCoins struct {
	AccountId cipher.Address //Account id
	Coin      string         //BTC,SKY
	Address   string         //address for deposit
	Tx        string         //transaction ID
}

//do not receive/credit coins until required number of confirmations
type ReceivedCoins struct {
	AccountId cipher.Address //Account id
	Coin      string         //BTC,SKY
	Address   string         //address for deposit
	Tx        string         //transaction ID
}

//manage incoming outputs and deposits
type IncomingManager struct {
	DepositAddresses []DepositAddress
	PendingIncoming  []PendingReceivedCoins
	ReceivedCoins    []ReceivedCoins
}

/*
Overall
*/

type AccountEngine struct {
	IncomingManager IncomingManager
	OutgoingManager OutgoingManager
}

func (self *AccountEngine) Save() {

}

func (self *AccountEngine) Load() {

}
