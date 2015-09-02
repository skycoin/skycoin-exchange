package skycoin_exchange

import (
	"github.com/skycoin/skycoin/src/cipher"
)

/*
Manages withdrawls and deposits
- associates external addresses with internal account IDs
- deposits and withdrawls
- credits and debits internal accounts

Keep a list of deposit addresses for each coin, for each account
*/

/*



*/

type PendingOutgoing struct {
	Coin    string //BTC, SKY
	Address cipher.Address
	Amount  uint64 //Satoshis or Drops
}

type AccountEngine struct {
}
