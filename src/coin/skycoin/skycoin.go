package skycoin_interface

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/op/go-logging.v1"

	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin/src/cipher"
	skycoin "github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
)

var (
	HideSeckey bool   = false
	ServeAddr  string = "127.0.0.1:6420"
	logger            = logging.MustGetLogger("exchange.skycoin")
	GatewayIns        = Gateway{}
)

type Utxo interface {
	GetHash() string
	GetSrcTx() string
	GetAddress() string
	GetCoins() uint64
	GetHours() uint64
}

type SkyUtxo struct {
	visor.ReadableOutput
}

type TxOut struct {
	skycoin.TransactionOutput
}

func (su SkyUtxo) GetHash() string {
	return su.Hash
}

func (su SkyUtxo) GetSrcTx() string {
	return su.SourceTransaction
}

func (su SkyUtxo) GetAddress() string {
	return su.Address
}

func (su SkyUtxo) GetCoins() uint64 {
	i, err := strconv.ParseUint(su.Coins, 10, 64)
	if err != nil {
		panic(err)
	}
	return i * 1e6
}

func (su SkyUtxo) GetHours() uint64 {
	return su.Hours
}

func MakeUtxoOutput(addr string, amount uint64, hours uint64) TxOut {
	uo := TxOut{}
	uo.Address = cipher.MustDecodeBase58Address(addr)
	uo.Coins = amount
	uo.Hours = hours
	return uo
}

func VerifyAmount(amt uint64) error {
	if (amt % 1e6) != 0 {
		return errors.New("Transaction amount must be multiple of 1e6 ")
	}
	return nil
}

// GenerateAddresses, generate bitcoin addresses.
func GenerateAddresses(seed []byte, num int) (string, []coin.AddressEntry) {
	sd, seckeys := cipher.GenerateDeterministicKeyPairsSeed(seed, num)
	entries := make([]coin.AddressEntry, num)
	for i, sec := range seckeys {
		pub := cipher.PubKeyFromSecKey(sec)
		entries[i].Address = cipher.AddressFromPubKey(pub).String()
		entries[i].Public = pub.Hex()
		if !HideSeckey {
			entries[i].Secret = sec.Hex()
		}
	}
	return fmt.Sprintf("%2x", sd), entries
}

// GetUnspentOutputs return the unspent outputs
func GetUnspentOutputs(addrs []string) ([]Utxo, error) {
	var url string
	if len(addrs) == 0 {
		return []Utxo{}, nil
	}

	addrParam := strings.Join(addrs, ",")
	url = fmt.Sprintf("http://%s/outputs?addrs=%s", ServeAddr, addrParam)

	rsp, err := http.Get(url)
	if err != nil {
		return []Utxo{}, errors.New("get skycoin outputs failed")
	}
	defer rsp.Body.Close()
	outputs := []SkyUtxo{}
	if err := json.NewDecoder(rsp.Body).Decode(&outputs); err != nil {
		return []Utxo{}, err
	}
	ux := make([]Utxo, len(outputs))
	for i, u := range outputs {
		ux[i] = u
	}
	return ux, nil
}

func getUnspentOutputsByHashes(hashes []string) ([]Utxo, error) {
	if len(hashes) == 0 {
		return []Utxo{}, nil
	}

	url := fmt.Sprintf("http://%s/outputs?hashes=%s", ServeAddr, strings.Join(hashes, ","))
	rsp, err := http.Get(url)
	if err != nil {
		return []Utxo{}, err
	}
	defer rsp.Body.Close()
	outputs := []SkyUtxo{}
	if err := json.NewDecoder(rsp.Body).Decode(&outputs); err != nil {
		return []Utxo{}, err
	}
	ux := make([]Utxo, len(outputs))
	for i, u := range outputs {
		ux[i] = u
	}
	return ux, nil
}
