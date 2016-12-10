package bitcoin_interface

import (
	//"errors"

	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"net/http"

	logging "github.com/op/go-logging"
	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin/src/cipher"
)

var (
	HideSeckey = false
	logger     = logging.MustGetLogger("exchange.bitcoin")
	// GatewayIns = Gateway{}
	Type = "bitcoin"
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

// TxOut bitcion transaction out struct
type TxOut struct {
	Addr  string
	Value uint64
}

// GenerateAddresses generates bitcoin addresses.
func GenerateAddresses(seed []byte, num int) (string, []coin.AddressEntry) {
	sd, seckeys := cipher.GenerateDeterministicKeyPairsSeed(seed, num)
	entries := make([]coin.AddressEntry, num)
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

// GetBalance query balance of address through the API of blockexplorer.com.
func GetBalance(addr []string) (uint64, error) {
	for _, a := range addr {
		if !validateAddress(a) {
			return 0, fmt.Errorf("invalid bitcoin address %v", a)
		}
	}

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

// NewUtxoWithKey create UtxoWithkey struct
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
		return []byte{}, fmt.Errorf("access %v failed", url)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	resp.Body.Close()
	return data, nil
}

func validateAddress(addr string) bool {
	_, err := cipher.BitcoinDecodeBase58Address(addr)
	return err == nil
}
