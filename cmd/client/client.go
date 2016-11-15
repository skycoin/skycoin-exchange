package main

import (
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"syscall"

	logging "github.com/op/go-logging"
	"github.com/skycoin/skycoin-exchange/src/client"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
)

const (
	servPubkey = "02942e46684114b35fe15218dfdc6e0d74af0446a397b8fcbf8b46fb389f756eb8"
)

var (
	logger     = logging.MustGetLogger("client.main")
	logFormat  = "[%{module}:%{level}] %{message}"
	logModules = []string{
		"client.main",
	}
)

func main() {
	var cfg client.Config
	home := util.UserHome()

	flag.StringVar(&cfg.ServAddr, "s", "localhost:8080", "server address")
	flag.IntVar(&cfg.Port, "p", 6060, "rpc port")
	flag.StringVar(&cfg.GuiDir, "gui-dir", "./src/web-app/static", "webapp static dir")
	flag.StringVar(&cfg.WalletDir, "wlt-dir", filepath.Join(home, ".exchange-client/wallet"), "wallet dir")
	flag.StringVar(&cfg.AccountDir, "account-dir", filepath.Join(home, ".exchange-client/account"), "account dir")

	flag.Parse()

	cfg.ServPubkey = cipher.MustPubKeyFromHex(servPubkey)
	cfg.GuiDir = util.ResolveResourceDirectory(cfg.GuiDir)

	// init logger.
	initLogging(logging.DEBUG, true)

	quit := make(chan int)
	go catchInterrupt(quit)

	// Watch for SIGUSR1
	go catchDebug()

	client.New(cfg).Run()

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
		logger.Error(err.Error())
		return
	}
	err = p.WriteTo(f, 2)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
