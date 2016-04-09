package client

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strings"
	// "github.com/btcsuite/btcd/wire"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcrpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/deiwin/interact"
	// "time"
	//ui "github.com/gizak/termui" // <- ui shortcut, optional
)

// var (
// 	checkNotEmpty = func(input string) error {
// 		// note that the inputs provided to these checks are already trimmed
// 		if input == "" {
// 			return errors.New("Input should not be empty!")
// 		}
// 		return nil
// 	}
// 	checkIsAPositiveNumber = func(input string) error {
// 		if n, err := strconv.Atoi(input); err != nil {
// 			return err
// 		} else if n < 0 {
// 			return errors.New("The number can not be negative!")
// 		}
// 		return nil
// 	}
// )

func GetCerts(app string) []byte {
	homeDir := btcutil.AppDataDir(app, false)
	certs, err := ioutil.ReadFile(filepath.Join(homeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	return certs
}

func BtClient(username string, password string, certs []byte, host string) *btcrpcclient.Client {
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
		User:         username,
		Pass:         password,
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

func ConvertToValue(param_type reflect.Type, param string) (reflect.Value, error) {
	switch param_type.Kind() {
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		i, err := strconv.Atoi(param)
		return reflect.ValueOf(i), err
	case reflect.Bool:
		b, err := strconv.ParseBool(param)
		return reflect.ValueOf(b), err
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(param, 64)
		return reflect.ValueOf(f), err
	case reflect.Interface:
		switch param_type.Name() {
		case "Address":
			a, err := btcutil.DecodeAddress(param, &chaincfg.MainNetParams)
			return reflect.ValueOf(a), err
		default:
			return reflect.ValueOf(nil), errors.New(fmt.Sprintf("Unhandled interface %s", param_type.Name()))
		}

	case reflect.String:
		return reflect.ValueOf(param), nil
	default:
		return reflect.ValueOf(nil), errors.New(fmt.Sprintf("Unhandled type %s", param_type))
	}
}

const help = `
Usage example:
	go run ./cmd/client/client.go -u $RPCUSER -p $RPCPASS -f ./examples/example.cdsl -d  WALLET_PASSPHRASE=$WALLET_PASSPHRASE

Assumes
	btcd.conf is setup
	btcwallet.conf is setup
	btcd is running
	btwallet is running
	~/go/bin/btcwallet -u $RPCUSER -P $RPCPASS --create

TODO:

deposit skycoin 1
	Prompts for type of coin
	Prompts for value
	Generate and return address for deposit
withdraw skycoin 1
	Prompts for type of coin
	Prompts for value
	Prompts for destination
bid skycoin 1 bitcoin 10
	Prompts for type of coin from
	Prompts for coin source quantity
	Prompts for type of coin to
	Prompts for coin target quantity
`

func LoadCommands(inputFile string) []string {
	inputs, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(inputs), "\n")
}

func Run() {

	username := flag.String("u", "skycoin", `specify the username for btcd/btwallet`)
	password := flag.String("p", "skycoin", `specify the password for btcd/btwallet`)
	filename := flag.String("f", "", `specify the input command-file to slurp`)
	defines := flag.String("d", "", `specify list of defines NAME1=VALUE1,NAME2=VALUE2`)
	verbose_flag := flag.Bool("v", false, `verbose flag`)

	flag.Parse()

	commands := []string{}
	counter := 0
	filename_s := string(*filename)
	if filename_s != "" {
		commands = append(commands, LoadCommands(filename_s)...)
	}

	username_s := string(*username)
	password_s := string(*password)
	variable_to_value := make(map[string]string)
	if *defines != "" {
		for _, kv := range strings.Split(string(*defines), ",") {
			p := strings.Split(kv, "=")
			variable_to_value[p[0]] = p[1]
		}
	}

	btcd := BtClient(username_s, password_s, GetCerts("btcd"), "localhost:8334")
	btwallet := BtClient(username_s, password_s, GetCerts("btcwallet"), "localhost:8332")
	default_current := reflect.ValueOf(btwallet)
	actor := interact.NewActor(os.Stdin, os.Stdout)
	verbose := *verbose_flag

	rvars := regexp.MustCompile(`\$\w+`)
	for {
		err := error(nil)
		command := ""

		if len(commands) > counter {
			command = commands[counter]
			counter++
		} else {
			command, err = actor.Prompt("command")
			if err != nil {
				log.Fatal(err)
			}
		}

		trimmed := strings.TrimSpace(command)
		if trimmed == "" {
			continue
		}

		subbed := rvars.ReplaceAllStringFunc(trimmed, func(s string) string {
			return variable_to_value[s[1:]]
		})
		fmt.Println(subbed)

		tokens := strings.Split(subbed, " ")
		cmd := tokens[0]

		switch cmd {
		case "help":
			fmt.Println(help)

		case "verbose":
			verbose = !verbose

		case "load":
			if len(tokens) > 1 {
				filename_s = tokens[1]
			} else {
				filename_s, err = actor.Prompt("filename")
				if err != nil {
					log.Fatal(err)
				}
			}
			commands = append(commands, LoadCommands(filename_s)...)

		case "use":
			target, err := actor.Prompt("target")
			if err != nil {
				log.Fatal(err)
			}
			switch target {
			case "btcd":
				default_current = reflect.ValueOf(btcd)
			case "btwallet":
				default_current = reflect.ValueOf(btwallet)
			}

		default:
			pkg_cmd := strings.Split(tokens[0], ".")
			current := default_current
			if len(pkg_cmd) > 1 {
				switch pkg_cmd[0] {
				case "btcd":
					current = reflect.ValueOf(btcd)
				case "btwallet":
					current = reflect.ValueOf(btwallet)
				}
				cmd = pkg_cmd[1]
			}

			method := current.MethodByName(cmd)
			if !method.IsValid() {
				fmt.Printf("Method '%s' not found\n", tokens[0])
				continue
			}
			args := tokens[1:]
			method_value := method.Interface()
			method_type := reflect.TypeOf(method_value)

			if verbose {
				fmt.Println(method_type)
			}

			in := []reflect.Value{}

			for i := 0; i < method_type.NumIn(); i++ {
				param_type := method_type.In(i)
				var param string
				if i < len(args) {
					param = args[i]
				} else {
					param, err = actor.Prompt(fmt.Sprintf("%s (%s)", param_type, param_type.Kind()))
					if err != nil {
						fmt.Printf("Invalid input %s\n", err)
						continue
					}
				}

				value, err := ConvertToValue(param_type, param)
				if err != nil {
					fmt.Printf("Could not parse: %s because %s\n", param, err)
					break
				}
				converted := value.Convert(param_type)
				in = append(in, converted)
			}

			return_values := method.Call(in)
			for i, r := range return_values {
				if verbose {
					fmt.Printf("%s: $%d=%s\n", r.Kind(), i, fmt.Sprint(r))
				} else {
					fmt.Println(r)
				}
				variable_to_value[fmt.Sprint(i)] = fmt.Sprint(r)
			}
		}
	}
}
