package bitcoin_interface

import (
	//"errors"

	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"gopkg.in/op/go-logging.v1"

	"net/http"

	"github.com/skycoin/skycoin-exchange/src/server/coin_interface"
	"github.com/skycoin/skycoin/src/cipher"
)

var (
	HideSeckey = false
	logger     = logging.MustGetLogger("exchange.bitcoin")
)

// Utxo unspent output
type Utxo interface {
	GetTxid() string
	GetVout() uint32
	GetAmount() uint64
	GetAddress() string
}

// UtxoWithkey unspent output with privkey.
type UtxoWithkey interface {
	Utxo
	GetPrivKey() string
}

type UtxoOut struct {
	Addr  string
	Value uint64
}

// GenerateAddresses, generate bitcoin addresses.
func GenerateAddresses(seed []byte, num int) (string, []coin_interface.AddressEntry) {
	sd, seckeys := cipher.GenerateDeterministicKeyPairsSeed(seed, num)
	entries := make([]coin_interface.AddressEntry, num)
	for i, sec := range seckeys {
		pub := cipher.PubKeyFromSecKey(sec)
		entries[i].Address = cipher.BitcoinAddressFromPubkey(pub)
		entries[i].Public = pub.Hex()
		if !HideSeckey {
			entries[i].Secret = cipher.BitcoinWalletImportFormatFromSeckey(sec)
		}
	}
	return fmt.Sprintf("%2x", sd), entries
}

// GetBalance, query balance of address through the API of blockexplorer.com.
func GetBalance(addr []string) (uint64, error) {
	// if AddressValid(addr) != nil {
	// 	log.Fatal("Address is invalid")
	// }
	// blkEplUrl := fmt.Sprintf("https://blockexplorer.com/api/addr/%s/balance", addr)
	addrs := strings.Join(addr, "|")
	blkChnUrl := fmt.Sprintf("https://blockchain.info/q/addressbalance/%s", addrs)
	data, err := getDataOfUrl(blkChnUrl)
	if err != nil {
		return 0, err
	}
	b, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}
	return uint64(b), nil
}

// GetUnspentOutputs return the unspent outputs
func GetUnspentOutputs(addrs []string) ([]Utxo, error) {
	return getUtxosBlkExplr(addrs)
}

func NewUtxoWithKey(utxo Utxo, key string) UtxoWithkey {
	return BlkExplrUtxoWithkey{
		BlkExplrUtxo: utxo.(BlkExplrUtxo),
		Privkey:      key,
	}
}

// getDataOfUrl, get data from specific URL.
func getDataOfUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	resp.Body.Close()
	return data, nil
}
