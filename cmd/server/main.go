package main

import (
	"flag"
	"os"

	"gopkg.in/op/go-logging.v1"

	"github.com/skycoin/skycoin-exchange/src/server"
	"github.com/skycoin/skycoin/src/cipher"
)

var (
	sk         = "38d010a84c7b9374352468b41b076fa585d7dfac67ac34adabe2bbba4f4f6257"
	logger     = logging.MustGetLogger("exchange.main")
	logFormat  = "[%{module}:%{level}] %{message}"
	logModules = []string{
		"exchange.main",
		"exchange.server",
		"exchange.account",
		"exchange.api",
		"exchange.bitcoin",
		"exchange.skycoin",
		"exchange.gin",
	}
)

func registerFlags(cfg *server.Config) {
	flag.IntVar(&cfg.Port, "port", 8080, "server listen port")
	flag.IntVar(&cfg.BtcFee, "btcFee", 10000, "transaction fee in satoish")
	flag.StringVar(&cfg.DataDir, "dataDir", ".skycoin-exchange", "data directory")
	flag.StringVar(&cfg.Seed, "seed", "", "wallet's seed")
	flag.IntVar(&cfg.UtxoPoolSize, "poolsize", 1000, "utxo pool size")

	flag.Set("logtostderr", "true")
	flag.Parse()
}

func main() {
	initLogging(logging.DEBUG, true)
	cfg := initConfig()
	s := server.New(cfg)
	s.Run()
}

func initConfig() server.Config {
	cfg := server.Config{}
	registerFlags(&cfg)
	if cfg.Seed == "" {
		logger.Error("seed must be set")
		flag.Usage()
	}
	cfg.WalletName = cfg.Seed + ".wlt"

	key, err := cipher.SecKeyFromHex(sk)
	if err != nil {
		logger.Fatal(err)
	}
	cfg.Seckey = key
	return cfg
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
