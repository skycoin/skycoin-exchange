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
