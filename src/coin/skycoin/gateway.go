package skycoin_interface

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/skycoin/skycoin-exchange/src/coin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

// Gateway skycoin gateway.
type Gateway struct{}

// GetTx get skycoin verbose transaction.
func (gw *Gateway) GetTx(txid string) (*pp.Tx, error) {
	url := fmt.Sprintf("http://%s/transaction?txid=%s", ServeAddr, txid)
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	tx := visor.TransactionResult{}
	if err := json.NewDecoder(rsp.Body).Decode(&tx); err != nil {
		return nil, err
	}
	return newPPTx(&tx), nil
}

// GetRawTx get raw tx by txid.
func (gw *Gateway) GetRawTx(txid string) (string, error) {
	url := fmt.Sprintf("http://%s/rawtx?txid=%s", ServeAddr, txid)
	rsp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	s, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	return strings.Trim(string(s), "\""), nil
}

// InjectTx inject skycoin transaction.
func (gw *Gateway) InjectTx(rawtx string) (string, error) {
	return BroadcastTx(rawtx)
}

// GetBalance get skycoin balance of specific addresses.
func (gw *Gateway) GetBalance(addrs []string) (pp.Balance, error) {
	url := fmt.Sprintf("http://%s/balance?addrs=%s", ServeAddr, strings.Join(addrs, ","))
	rsp, err := http.Get(url)
	if err != nil {
		return pp.Balance{}, err
	}
	defer rsp.Body.Close()
	bal := struct {
		Confirmed wallet.Balance `json:"confirmed"`
		Predicted wallet.Balance `json:"predicted"`
	}{}
	if err := json.NewDecoder(rsp.Body).Decode(&bal); err != nil {
		return pp.Balance{}, err
	}
	return pp.Balance{
		Amount: pp.PtrUint64(bal.Confirmed.Coins),
		Hours:  pp.PtrUint64(bal.Confirmed.Hours)}, nil
}

func (gw *Gateway) ValidateTxid(txid string) bool {
	_, err := cipher.SHA256FromHex(txid)
	return err == nil
}

func newPPTx(tx *visor.TransactionResult) *pp.Tx {
	return &pp.Tx{
		Sky: &pp.SkyTx{
			Length:    pp.PtrUint32(tx.Transaction.Length),
			Type:      pp.PtrInt32(int32(tx.Transaction.Type)),
			Hash:      pp.PtrString(tx.Transaction.Hash),
			InnerHash: pp.PtrString(tx.Transaction.InnerHash),
			Sigs:      tx.Transaction.Sigs,
			Inputs:    tx.Transaction.In,
			Outputs:   newSkyTxOutputArray(tx.Transaction.Out),
			Unknow:    pp.PtrBool(tx.Status.Unknown),
			Confirmed: pp.PtrBool(tx.Status.Confirmed),
			Height:    pp.PtrUint64(tx.Status.Height),
		},
	}
}

func newSkyTxOutputArray(ops []visor.ReadableTransactionOutput) []*pp.SkyTxOutput {
	outs := make([]*pp.SkyTxOutput, len(ops))
	for i, op := range ops {
		outs[i] = &pp.SkyTxOutput{
			Hash:    pp.PtrString(op.Hash),
			Address: pp.PtrString(op.Address),
			Coins:   pp.PtrString(op.Coins),
			Hours:   pp.PtrUint64(op.Hours),
		}
	}
	return outs
}

// CreateRawTx create skycoin raw transaction.
func (gw Gateway) CreateRawTx(txIns []coin.TxIn, txOuts interface{}) (string, error) {
	tx := Transaction{}
	// keys := make([]cipher.SecKey, len(utxos))
	for _, in := range txIns {
		tx.PushInput(cipher.MustSHA256FromHex(in.Txid))
	}

	s := reflect.ValueOf(txOuts)
	if s.Kind() != reflect.Slice {
		return "", errors.New("error tx out type")
	}
	outs := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		outs[i] = s.Index(i).Interface()
	}

	if len(outs) > 2 {
		return "", errors.New("out address more than 2")
	}

	for _, o := range outs {
		out := o.(TxOut)
		if (out.Coins % 1e6) != 0 {
			return "", errors.New("skycoin coins must be multiple of 1e6")
		}
		tx.PushOutput(out.Address, out.Coins, out.Hours)
	}

	tx.UpdateHeader()
	d, err := tx.Serialize()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(d), nil
}

// SignRawTx sign skycoin transaction.
func (gw Gateway) SignRawTx(rawtx string, getKey coin.GetPrivKey) (string, error) {
	// decode the rawtx
	tx := Transaction{}
	b, err := hex.DecodeString(rawtx)
	if err != nil {
		return "", err
	}
	if err := tx.Deserialize(bytes.NewBuffer(b)); err != nil {
		return "", err
	}

	// TODO: need to get the address of the in hash, then get key of those address, and sign.
	hashes := make([]string, len(tx.In))
	for i, in := range tx.In {
		hashes[i] = in.Hex()
	}

	// get utxos of thoes hashes.
	utxos, err := getUnspentOutputsByHashes(hashes)
	if err != nil {
		return "", err
	}

	if len(utxos) != len(hashes) {
		return "", errors.New("failed to search tx in's address")
	}

	hashAddrMap := map[string]string{}
	for _, u := range utxos {
		hashAddrMap[u.GetHash()] = u.GetAddress()
	}

	keys := make([]cipher.SecKey, len(hashes))
	for i, h := range hashes {
		key, err := getKey(hashAddrMap[h])
		if err != nil {
			return "", err
		}

		keys[i], err = cipher.SecKeyFromHex(key)
		if err != nil {
			return "", err
		}
	}

	tx.SignInputs(keys)
	tx.UpdateHeader()
	d, err := tx.Serialize()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(d), nil
}
