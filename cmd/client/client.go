package main

import (
	"github.com/skycoin/skycoin-exchange/src/rpclient"
	"github.com/skycoin/skycoin/src/cipher"
)

const (
	ServPubkey = "02942e46684114b35fe15218dfdc6e0d74af0446a397b8fcbf8b46fb389f756eb8"
)

func main() {
	pk := cipher.MustPubKeyFromHex(ServPubkey)

	cfg := rpclient.Config{
		ApiRoot:    "http://localhost:8080/api/v1",
		AcntName:   "account.data",
		ServPubkey: pk,
	}

	cli := rpclient.New(cfg)
	cli.Run(":6060")
}
