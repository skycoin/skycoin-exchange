package main

import "github.com/skycoin/skycoin-exchange/src/server"

func main() {
	cfg := server.Config{
		Port:          6060,
		WalletDataDir: ""}
	s := server.New(cfg)
	s.Run()
}
