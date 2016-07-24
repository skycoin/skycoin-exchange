package skycoin_interface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin-exchange/src/server/coin_interface"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
)

var (
	HideSeckey bool   = false
	ServeAddr  string = "http://127.0.0.1:6420"
)

// Hash              string `json:"txid"` //hash uniquely identifies transaction
// SourceTransaction string `json:"src_tx"`
// Address           string `json:"address"`
// Coins             string `json:"coins"`
// Hours             uint64 `json:"hours"`

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

type UtxoOut struct {
	coin.TransactionOutput
}

type Transaction struct {
	coin.Transaction
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
	return i
}

func (su SkyUtxo) GetHours() uint64 {
	return su.Hours
}

// GenerateAddresses, generate bitcoin addresses.
func GenerateAddresses(seed []byte, num int) (string, []coin_interface.AddressEntry) {
	sd, seckeys := cipher.GenerateDeterministicKeyPairsSeed(seed, num)
	entries := make([]coin_interface.AddressEntry, num)
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
	if len(addrs) > 0 {
		addrParam := strings.Join(addrs, ",")
		url = fmt.Sprintf("%s/outputs?addrs=%s", ServeAddr, addrParam)
	} else {
		url = fmt.Sprintf("%s/outputs")
	}

	rsp, err := http.Get(url)
	if err != nil {
		return []Utxo{}, fmt.Errorf("get outputs from %s faild", url)
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

// NewTransaction create skycoin transaction.
func NewTransaction(utxos []Utxo, keyMap map[string]cipher.SecKey, outs []UtxoOut) Transaction {
	tx := Transaction{}
	keys := make([]cipher.SecKey, len(utxos))
	for i, u := range utxos {
		tx.PushInput(cipher.MustSHA256FromHex(u.GetHash()))
		keys[i] = keyMap[u.GetAddress()]
	}

	for _, o := range outs {
		tx.PushOutput(o.Address, o.Coins, o.Hours)
	}

	tx.SignInputs(keys)
	tx.UpdateHeader()
	return tx
}

// BroadcastTx
func BroadcastTx(tx Transaction) (string, error) {
	rawtx := tx.Serialize()
	v := struct {
		Rawtx []byte `json:"rawtx"`
	}{
		rawtx,
	}

	d, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/injectTransaction", ServeAddr)
	rsp, err := http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return "", fmt.Errorf("post rawtx to %s failed", url)
	}
	defer rsp.Body.Close()
	rslt := struct {
		Txid string `json:"txid"`
	}{}
	if err := json.NewDecoder(rsp.Body).Decode(&rslt); err != nil {
		return "", err
	}
	return rslt.Txid, nil
}
