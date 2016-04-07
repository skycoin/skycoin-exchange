package client

import (
	"errors"
	"fmt"
	"reflect"
	// "github.com/btcsuite/btcd/wire"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/btcsuite/btcrpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/deiwin/interact"
	// "time"
	//ui "github.com/gizak/termui" // <- ui shortcut, optional
)

var (
	checkNotEmpty = func(input string) error {
		// note that the inputs provided to these checks are already trimmed
		if input == "" {
			return errors.New("Input should not be empty!")
		}
		return nil
	}
	checkIsAPositiveNumber = func(input string) error {
		if n, err := strconv.Atoi(input); err != nil {
			return err
		} else if n < 0 {
			return errors.New("The number can not be negative!")
		}
		return nil
	}
)

func GetCerts(app string) []byte {
	homeDir := btcutil.AppDataDir(app, false)
	certs, err := ioutil.ReadFile(filepath.Join(homeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	return certs
}

func BtClient(certs []byte, host string) *btcrpcclient.Client {
	// Only override the handlers for notifications you care about.
	// Also note most of these handlers will only be called if you register
	// for notifications.  See the documentation of the btcrpcclient
	// NotificationHandlers type for more details about each handler.
	ntfnHandlers := btcrpcclient.NotificationHandlers{
	// OnBlockConnected: func(hash *wire.ShaHash, height int32, time time.Time) {
	// 	log.Printf("Block connected: %v (%d) %v", hash, height, time)
	// },
	// OnBlockDisconnected: func(hash *wire.ShaHash, height int32, time time.Time) {
	// 	log.Printf("Block disconnected: %v (%d) %v", hash, height, time)
	// },
	}
	connCfg := &btcrpcclient.ConnConfig{
		Host:         host,
		Endpoint:     "ws",
		User:         "skycoin",
		Pass:         "skycoin2016",
		Certificates: certs,
	}
	client, err := btcrpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		log.Fatal(err)
	}

	// Register for block connect and disconnect notifications.
	if err := client.NotifyBlocks(); err != nil {
		log.Fatal(err)
	}
	log.Println("NotifyBlocks: Registration Complete")

	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block count: %d", blockCount)

	return client
}

func ConvertToString(r reflect.Value) string {
	switch r.Kind() {
	case reflect.Uint64:
		return fmt.Sprintf("%d", r)
	default:
		log.Fatal("Unsupported kind", r.Kind())
	}
	log.Fatal("Unsupported kind", r.Kind())
	return "not reachable"
}

func Run() {
	help := `
	This is the help.

	Deposit
		Prompts for type of coin
		Prompts for value
		Generate and return address for deposit
	Withdraw
		Prompts for type of coin
		Prompts for value
		Prompts for destination

	`
	btcd := BtClient(GetCerts("btcd"), "localhost:8334")
	btwallet := BtClient(GetCerts("btcwallet"), "localhost:8332")
	current := reflect.ValueOf(btwallet)
	actor := interact.NewActor(os.Stdin, os.Stdout)

	for {
		command, err := actor.Prompt("command")
		if err != nil {
			log.Fatal(err)
		}
		switch command {
		case "help":
			fmt.Println(help)
		case "use":
			target, err := actor.Prompt("target")
			if err != nil {
				log.Fatal(err)
			}
			switch target {
			case "btcd":
				current = reflect.ValueOf(btcd)
			case "btwallet":
				current = reflect.ValueOf(btwallet)
			}
		case "ListAccounts":
			accounts, e := btwallet.ListAccounts()
			if e != nil {
				fmt.Println("Error", e)
				continue
			}
			for k, v := range accounts {
				fmt.Printf("%s = %s\n", k, v.String())
			}
		case "CheckBalance":
		case "Spend":
		default:
			//in := []reflect.Value{reflect.ValueOf(nil)}
			method := current.MethodByName(command)
			if !method.IsValid() {
				fmt.Printf("Method %s not found\n", command)
				continue
			}
			method_value := method.Interface()
			method_type := reflect.TypeOf(method_value)
			fmt.Println(method_type)

			in := []reflect.Value{}

			for i := 0; i < method_type.NumIn(); i++ {
				param_type := method_type.In(i)
				param, err := actor.Prompt(fmt.Sprintf("%s", param_type))
				if err != nil {
					fmt.Printf("Invalid input %s\n", err)
					continue
				}
				switch param_type.Kind() {
				case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
					converted, err := strconv.Atoi(param)
					if err != nil {
						fmt.Printf("Invalid input %s\n", err)
						continue
					}
					in = append(in, reflect.ValueOf(converted).Convert(param_type))
				case reflect.Bool:
					converted, err := strconv.ParseBool(param)
					if err != nil {
						fmt.Printf("Invalid input %s\n", err)
						continue
					}
					in = append(in, reflect.ValueOf(converted).Convert(param_type))
				case reflect.Float32, reflect.Float64:
					converted, err := strconv.ParseFloat(param, 64)
					if err != nil {
						fmt.Printf("Invalid input %s\n", err)
						continue
					}
					in = append(in, reflect.ValueOf(converted).Convert(param_type))
				default:
					in = append(in, reflect.ValueOf(param).Convert(param_type))
				}

			}

			return_values := method.Call(in)
			for _, r := range return_values {
				fmt.Printf("%s: ", r.Kind())
				fmt.Println(r)
			}
		}
	}
}
