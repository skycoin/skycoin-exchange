package skycoin_exchange

import (
	"github.com/skycoin/skycoin/src/cipher"
	"time"
)

/*
	User sending coins to exchange
	Handle incoming, deposit addresses and received coins

	Where are deposit address, account pairs stored?

	This is a wallet
	- credits user when coins are deposited
	- handles withdraws
*/

type DepositAddress struct {
	AccountId AccountIdentifier //Account id
	Coin      string            //BTC,SKY
	Address   string            //address for deposit
	Tx        string            //transaction ID
}

type PendingReceivedCoins struct {
	AccountId AccountIdentifier //Account id
	Coin      string            //BTC,SKY
	Address   string            //address for deposit
	Tx        string            //transaction ID
}

//do not receive/credit coins until required number of confirmations
type ReceivedCoins struct {
	AccountId AccountIdentifier //Account id
	Coin      string            //BTC,SKY
	Address   string            //address for deposit
	Tx        string            //transaction ID
}

//manage incoming outputs and deposits
type IncomingManager struct {
	DepositAddresses []DepositAddress
	PendingIncoming  []PendingReceivedCoins
	ReceivedCoins    []ReceivedCoins
}

/*
	Check for new deposits or state
*/

func (self *IncomingManager) Tick() {

}
