package server

import (
	"fmt"

	"github.com/skycoin/skycoin-exchange/src/server/account"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
)

type Server struct {
	cfg Config
	account.AccountManager
}

type Config struct {
	Port          int
	WalletDataDir string
}

func New(cfg Config) *Server {
	s := &Server{
		cfg:            cfg,
		AccountManager: account.NewExchangeAccountManager(),
	}
	return s
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

func (self *Server) Run() {
	// init the wallet package.
	wallet.Init(self.cfg.WalletDataDir)

	fmt.Println("skycoin-exchange server started, port:%d", self.cfg.Port)

	r := NewRouter(self)

	r.Run(fmt.Sprintf(":%d", self.cfg.Port))
}
