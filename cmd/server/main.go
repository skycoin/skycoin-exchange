package main

//"encoding/json"
//"errors"
// "fmt"
//"github.com/go-goodies/go_oops"
//"github.com/l3x/jsoncfgo"
// "html/template"
//"io/ioutil"
// "log"
// "net/http"
//"regexp"
// "github.com/skycoin/skycoin/src/cipher"
// "os"

import (
	skycoin_exchange "github.com/skycoin/skycoin-exchange/src/server"
)

func main() {

	server := skycoin_exchange.Server

	server.Init()
	server.Run()
}
