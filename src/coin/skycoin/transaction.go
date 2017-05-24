package skycoin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	sky "github.com/skycoin/skycoin/src/coin"
)

type Transaction struct {
	sky.Transaction
}

// NewTransaction create skycoin transaction.
func NewTransaction(utxos []Utxo, keys []cipher.SecKey, outs []TxOut) *Transaction {
	tx := Transaction{}
	// keys := make([]cipher.SecKey, len(utxos))
	for _, u := range utxos {
		tx.PushInput(cipher.MustSHA256FromHex(u.GetHash()))
	}

	for _, o := range outs {
		tx.PushOutput(o.Address, o.Coins, o.Hours)
	}
	// tx.Verify()

	tx.SignInputs(keys)
	tx.UpdateHeader()
	return &tx
}

// BroadcastTx skycoin broadcast tx.
func BroadcastTx(nodeAddr, rawtx string) (string, error) {
	v := struct {
		Rawtx string `json:"rawtx"`
	}{
		rawtx,
	}

	d, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("http://%s/injectTransaction", nodeAddr)
	rsp, err := http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return "", fmt.Errorf("post rawtx to %s failed", url)
	}
	defer rsp.Body.Close()
	s, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	return strings.Trim(string(s), "\""), nil
}

func (tx *Transaction) Serialize() ([]byte, error) {
	return tx.Transaction.Serialize(), nil
}

func (tx *Transaction) Deserialize(r io.Reader) error {
	d, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if err := encoder.DeserializeRaw(d, tx); err != nil {
		return err
	}
	return nil
}
