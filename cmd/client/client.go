package main

import (
	"flag"
	"fmt"

	"github.com/skycoin/skycoin-exchange/src/rpclient"
	"github.com/skycoin/skycoin/src/cipher"
)

const (
	ServPubkey = "02942e46684114b35fe15218dfdc6e0d74af0446a397b8fcbf8b46fb389f756eb8"
)

func main() {
	servAddr := flag.String("s", "localhost:8080", "server address")
	port := flag.Int("port", 6060, "rpc port")
	flag.Parse()

	pk := cipher.MustPubKeyFromHex(ServPubkey)

	cfg := rpclient.Config{
		ApiRoot:    *servAddr,
		ServPubkey: pk,
	}

	svr := rpclient.New(cfg)
	svr.Run(fmt.Sprintf(":%d", *port))
}
