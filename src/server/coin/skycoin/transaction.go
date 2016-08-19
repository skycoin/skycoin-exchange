package skycoin_interface

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/skycoin/encoder"
	"github.com/skycoin/skycoin/src/cipher"
	skycoin "github.com/skycoin/skycoin/src/coin"
)

type Transaction struct {
	skycoin.Transaction
}

// NewTransaction create skycoin transaction.
func NewTransaction(utxos []Utxo, keys []cipher.SecKey, outs []UtxoOut) *Transaction {
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

// BroadcastTx
func BroadcastTx(tx Transaction) (string, error) {
	rawtx, err := tx.Serialize()
	if err != nil {
		return "", err
	}

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
		Success bool   `json:"success"`
		Reason  string `json:"reason"`
		Txid    string `json:"txid"`
	}{}

	if err := json.NewDecoder(rsp.Body).Decode(&rslt); err != nil {
		return "", err
	}
	if rslt.Success {
		return rslt.Txid, nil
	}
	return "", errors.New(rslt.Reason)
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
