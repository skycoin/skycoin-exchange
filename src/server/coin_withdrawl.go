package skycoin_exchange

/*
How can it be represented?
- should use hashes?
- event types and transactions?

*/

/*
- exchange sending coins to user


*/

type PendingOutgoing struct {
	AccountId AccountID //Account id
	Coin      string    //BTC, SKY
	Address   string    //cipher.Address
	Amount    uint64    //Satoshis or Drops
}

type CompletedOutgoing struct {
	AccountId AccountID //Account id
	Coin      string    //BTC, SKY
	Address   string    //cipher.Address
	Amount    uint64    //Satoshis or Drops
	Tx        string    //transaction ID
}

type OutgoingManager struct {
	//FreezeWithdraws               bool
	PendingOutgoingTransactions   []PendingOutgoing
	CompletedOutgoingTransactions []CompletedOutgoing
}

func (self *OutgoingManager) Tick() {

}
