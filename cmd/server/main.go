package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/skycoin/skycoin-exchange/src/server"
	"github.com/skycoin/skycoin/src/cipher"
)

var sk = "38d010a84c7b9374352468b41b076fa585d7dfac67ac34adabe2bbba4f4f6257"

func registerFlags(cfg *server.Config) {
	flag.IntVar(&cfg.Port, "port", 8080, "server listen port")
	flag.IntVar(&cfg.Fee, "fee", 10000, "transaction fee in satoish")
	flag.StringVar(&cfg.DataDir, "dataDir", ".skycoin-exchange", "data directory")
	flag.StringVar(&cfg.Seed, "seed", "", "wallet's seed")
	flag.StringVar(&cfg.AcntName, "acntName", "account.data", "accounts file name")
	flag.IntVar(&cfg.UtxoPoolSize, "poolsize", 1000, "utxo pool size")

	// flag.Set("log_dir", logDir)
	flag.Set("logtostderr", "true")
	flag.Parse()
}

func main() {
	cfg := server.Config{}
	registerFlags(&cfg)
	cfg.WalletName = cfg.Seed + ".wlt"
	key, err := cipher.SecKeyFromHex(sk)
	if err != nil {
		glog.Fatal(err)
	}
	cfg.Seckey = key
	s := server.New(cfg)
	s.Run()
}
