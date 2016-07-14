package main

import "github.com/skycoin/skycoin-exchange/src/rpclient"

const (
	ServPubkey = "02942e46684114b35fe15218dfdc6e0d74af0446a397b8fcbf8b46fb389f756eb8"
)

func main() {
	cli := rpclient.New("http://localhost:8080/api/v1", ServPubkey)
	cli.Run(":6060")
}
