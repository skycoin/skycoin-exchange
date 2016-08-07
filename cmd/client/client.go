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
	apiUrl := flag.String("api_url", "http://localhost:8080/api/v1", "server api root")
	port := flag.Int("port", 6060, "rpc port")
	flag.Parse()

	pk := cipher.MustPubKeyFromHex(ServPubkey)

	cfg := rpclient.Config{
		ApiRoot:    *apiUrl,
		ServPubkey: pk,
	}

	cli := rpclient.New(cfg)
	cli.Run(fmt.Sprintf(":%d", *port))
}
