package server

import (
	"fmt"
	"log"
	"time"

	"github.com/skycoin/skycoin-exchange/src/server/account"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

type Server interface {
	Run()
	CreateAccountWithPubkey(pubkey cipher.PubKey) (account.Accounter, error)
	GetAccount(id account.AccountID) (account.Accounter, error)
	GetNonceKeyLifetime() time.Duration
}

// Config store server's configuration.
type Config struct {
	Port             int
	WalletDataDir    string
	NonceKeyLifetime time.Duration
}

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

type ExchangeServer struct {
	account.AccountManager
	cfg Config
}

func New(cfg Config) Server {
	s := &ExchangeServer{
		cfg:            cfg,
		AccountManager: account.NewExchangeAccountManager(),
	}
	return s
}

func (self *ExchangeServer) Run() {
	// init the wallet package.
	wallet.Init(self.cfg.WalletDataDir)

	log.Println("skycoin-exchange server started, port:", self.cfg.Port)

	r := NewRouter(self)

	r.Run(fmt.Sprintf(":%d", self.cfg.Port))
}

func (self ExchangeServer) GetNonceKeyLifetime() time.Duration {
	return self.cfg.NonceKeyLifetime
}
