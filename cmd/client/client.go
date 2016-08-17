package main

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/op/go-logging.v1"

	"github.com/skycoin/skycoin-exchange/src/rpclient"
	"github.com/skycoin/skycoin/src/cipher"
)

const (
	ServPubkey = "02942e46684114b35fe15218dfdc6e0d74af0446a397b8fcbf8b46fb389f756eb8"
)

var (
	logger     = logging.MustGetLogger("client.main")
	logFormat  = "[%{module}:%{level}] %{message}"
	logModules = []string{
		"client.main",
	}
)

func main() {
	servAddr := flag.String("s", "localhost:8080", "server address")
	port := flag.Int("port", 6060, "rpc port")
	flag.Parse()

	// init logger.
	initLogging(logging.DEBUG, true)

	pk := cipher.MustPubKeyFromHex(ServPubkey)
	cfg := rpclient.Config{
		ApiRoot:    *servAddr,
		ServPubkey: pk,
	}

	svr := rpclient.New(cfg)
	svr.Run(fmt.Sprintf(":%d", *port))
}

func initLogging(level logging.Level, color bool) {
	format := logging.MustStringFormatter(logFormat)
	logging.SetFormatter(format)
	bk := logging.NewLogBackend(os.Stdout, "", 0)
	bk.Color = true
	bkLvd := logging.AddModuleLevel(bk)
	for _, s := range logModules {
		bkLvd.SetLevel(level, s)
	}

	logging.SetBackend(bkLvd)
}
