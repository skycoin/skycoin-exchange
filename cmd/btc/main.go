package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strconv"

	"os/user"

	"os"

	"bytes"
	"os/exec"

	"github.com/skycoin/skycoin-exchange/src/api/mobile"
)

var currentWltID string

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	wltDirPath := filepath.Join(usr.HomeDir, ".btcd/wallet")
	if _, err := os.Stat(wltDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(wltDirPath, 0700); err != nil {
			panic(err)
		}
	}

	servAddr := flag.String("s", "127.0.0.1:8080", "exchange server addr")
	pubKey := flag.String("p", "02942e46684114b35fe15218dfdc6e0d74af0446a397b8fcbf8b46fb389f756eb8", "server pubkey")
	flag.Parse()

	cfg := mobile.Config{
		WalletDirPath: wltDirPath,
		ServerAddr:    *servAddr,
		ServerPubkey:  *pubKey,
	}

	mobile.Init(&cfg)

	if len(os.Args) == 1 {
		fmt.Println("Name:")
		fmt.Println("\tbtc-cli - the simple bitcoin command line interface")
		fmt.Println("USAGE:")
		fmt.Println("\tbtc-cli [global options] command [arguments...]")
		fmt.Println("COMMANDS:")
		fmt.Println("\tsend			send bitcoin to given address")
		fmt.Println("\tnewWallet		create wallet")
		fmt.Println("\tnewAddress		create address in specific wallet")
		fmt.Println("\tlistWallets		list all wallets")
		fmt.Println("\twallet			see the wallet content")
		fmt.Println("\tbalance			check the balance of address")
		fmt.Println("GLOBAL OPTIONS:")
		fmt.Println("\t-s		exchange server address, default: 127.0.0.1:8080")
		fmt.Println("\t-p		exchange pubkey")
		return
	}

	switch os.Args[1] {
	case "send":
		if len(os.Args) != 5 {
			fmt.Println("send $wallet_id $recv_addr $amount")
			return
		}
		wltID := os.Args[2]
		recvAddr := os.Args[3]
		amount := os.Args[4]
		if err != nil {
			fmt.Println("Invalid amount,", err)
			return
		}
		s, err := mobile.SendBtc(wltID, recvAddr, amount, "2000")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(s)
	case "newWallet":
		if len(os.Args) != 3 {
			fmt.Println("newWallet $seed")
			return
		}
		seed := os.Args[2]
		id, err := mobile.NewWallet("bitcoin", seed)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("wallet id:", id)
	case "newAddress":
		if len(os.Args) != 4 {
			fmt.Println(`newAddress $wallet_id $address_num`)
			return
		}
		wltID := os.Args[2]
		numStr := os.Args[3]
		n, err := strconv.Atoi(numStr)
		if err != nil {
			fmt.Println(err)
			return
		}
		s, err := mobile.NewAddress(wltID, n)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(s)
	case "listWallets":
		cmd := exec.Command("ls", wltDirPath)
		out := bytes.Buffer{}
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(out.String())
	case "wallet":
		wltID := os.Args[2]
		wlt := filepath.Join(wltDirPath, wltID)
		fmt.Println(wlt)
		cmd := exec.Command("less", wlt)
		out := bytes.Buffer{}
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(out.String())
	case "balance":
		addr := os.Args[2]
		s, err := mobile.GetBalance("bitcoin", addr)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(s)
	}
}
