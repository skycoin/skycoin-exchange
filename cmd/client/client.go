package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"gopkg.in/op/go-logging.v1"

	"github.com/skycoin/skycoin-exchange/src/rpclient"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
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
		APIRoot:    *servAddr,
		ServPubkey: pk,
	}

	svr := rpclient.New(cfg)

	quit := make(chan int)
	go catchInterrupt(quit)

	// Watch for SIGUSR1
	go catchDebug()

	staticDir := util.ResolveResourceDirectory("./src/web-app/static")
	svr.Run(fmt.Sprintf("127.0.0.1:%d", *port), staticDir)

	<-quit

	logger.Info("Goodbye")
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

func catchInterrupt(quit chan<- int) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	signal.Stop(sigchan)
	quit <- 1
}

// Catches SIGUSR1 and prints internal program state
func catchDebug() {
	sigchan := make(chan os.Signal, 1)
	//signal.Notify(sigchan, syscall.SIGUSR1)
	signal.Notify(sigchan, syscall.Signal(0xa)) // SIGUSR1 = Signal(0xa)
	for {
		select {
		case <-sigchan:
			printProgramStatus()
		}
	}
}

func printProgramStatus() {
	fn := "goroutine.prof"
	logger.Debug("Writing goroutine profile to %s", fn)
	p := pprof.Lookup("goroutine")
	f, err := os.Create(fn)
	defer f.Close()
	if err != nil {
		logger.Error("%v", err)
		return
	}
	err = p.WriteTo(f, 2)
	if err != nil {
		logger.Error("%v", err)
		return
	}
}
