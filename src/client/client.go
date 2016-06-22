package client

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcrpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/deiwin/interact"
	"github.com/skycoin/skycoin-exchange/src/skyclient"
)

func GetCerts(app string) []byte {
	homeDir := btcutil.AppDataDir(app, false)
	certs, err := ioutil.ReadFile(filepath.Join(homeDir, "rpc.cert"))
	fmt.Println(filepath.Join(homeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	return certs
}

func BtClient(config *Config, certs []byte, host string) *btcrpcclient.Client {
	ntfnHandlers := btcrpcclient.NotificationHandlers{}
	connCfg := &btcrpcclient.ConnConfig{
		Host:         host,
		Endpoint:     "ws",
		User:         config.username,
		Pass:         config.password,
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

func SkyClient(config *Config, host string) (*skyclient.Client, error) {
	return skyclient.New()
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
	btcwallet is running
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

type Config struct {
	username       string
	password       string
	filename       string
	defines        string
	verbose        bool
	testnet        bool
	btcd_host      string
	btcwallet_host string
}

func Run() {
	var config Config = Config{
		username:       "skycoin",
		password:       "skycoin",
		filename:       "",
		defines:        "",
		verbose:        false,
		testnet:        false,
		btcd_host:      "localhost:8334",
		btcwallet_host: "localhost:8332",
	}

	flag.StringVar(&config.username, "u", config.username, `specify the username for btcd/btcwallet`)
	flag.StringVar(&config.password, "p", config.password, `specify the password for btcd/btcwallet`)
	flag.StringVar(&config.filename, "f", config.filename, `specify the input command-file to slurp`)
	flag.StringVar(&config.defines, "d", config.defines, `specify list of defines NAME1=VALUE1,NAME2=VALUE2`)
	flag.StringVar(&config.btcd_host, "b", config.btcd_host, `btcd host`)
	flag.StringVar(&config.btcwallet_host, "w", config.btcwallet_host, `btcwallet host`)
	flag.BoolVar(&config.verbose, "v", config.verbose, `verbose flag`)
	flag.BoolVar(&config.testnet, "t", config.verbose, `testnet flag`)

	flag.Parse()

	if config.testnet {
		config.btcd_host = "localhost:18334"
		config.btcwallet_host = "localhost:18332"
	}

	commands := []string{}
	counter := 0
	if config.filename != "" {
		commands = append(commands, LoadCommands(config.filename)...)
	}

	variable_to_value := make(map[string]string)
	if config.defines != "" {
		for _, kv := range strings.Split(string(config.defines), ",") {
			p := strings.Split(kv, "=")
			variable_to_value[p[0]] = p[1]
		}
	}

	btcd := BtClient(&config, GetCerts("btcd"), config.btcd_host)
	btcwallet := BtClient(&config, GetCerts("btcwallet"), config.btcwallet_host)
	skycoin := SkyClient(&config, config.btcwallet_host)

	actor := interact.NewActor(os.Stdin, os.Stdout)

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
			config.verbose = !config.verbose

		case "load":
			var filename string
			if len(tokens) > 1 {
				filename = tokens[1]
			} else {
				filename, err = actor.Prompt("filename")
				if err != nil {
					log.Fatal(err)
				}
			}
			commands = append(commands, LoadCommands(filename)...)

		default:
			pkg_cmd := strings.Split(tokens[0], ".")
			current := reflect.ValueOf(nil)
			if len(pkg_cmd) > 1 {
				switch pkg_cmd[0] {
				case "btcd":
					current = reflect.ValueOf(btcd)
				case "btcwallet":
					current = reflect.ValueOf(btcwallet)
				case "skycoin":
					current = reflect.ValueOf(skycoin)
				default:
					log.Fatal("Unsupported target ", pkg_cmd[0])
				}
				cmd = pkg_cmd[1]
			} else {
				log.Fatal("Most specify target for reflection ", tokens)
			}

			method := current.MethodByName(cmd)
			if !method.IsValid() {
				fmt.Printf("Method '%s' not found\n", tokens[0])
				continue
			}
			args := tokens[1:]
			method_value := method.Interface()
			method_type := reflect.TypeOf(method_value)

			if config.verbose {
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
				if config.verbose {
					fmt.Printf("%s: $%d=%s\n", r.Kind(), i, fmt.Sprint(r))
				} else {
					fmt.Println(r)
				}
				variable_to_value[fmt.Sprint(i)] = fmt.Sprint(r)
			}
		}
	}
}
