package main

import (
	"flag"
	//"fmt"
	//"skycoin/"
	bitcoin_interface "github.com/skycoin/skycoin-exchange/src/server/coin_interface/bitcoin"

	"time"
)

/*
	Listens for bitcoin unspent outputs on input address
*/

//bitcoin_interface

var Addr string

//https://blockchain.info/unspent?active=1SakrZuzQmGwn7MSiJj5awqJZjSYeBWC3

func main() {
	//flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
	//func StringVar(p *string, name string, value string, usage string)
	flag.StringVar(&Addr, "Address", "1SakrZuzQmGwn7MSiJj5awqJZjSYeBWC3", "address to watch for outputs")

	btcd := bitcoin_interface.Manager{}
	btcd.Init()
	btcd.AddWatchAddress(Addr)

	for {
		btcd.Tick()
		time.Sleep(1 * time.Second)
	}

}
