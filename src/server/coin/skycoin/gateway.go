package skycoin_interface

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin/src/visor"
)

type Gateway struct{}

func (gw *Gateway) GetTx(txid string) (*pp.Tx, error) {
	url := fmt.Sprintf("%s/transaction?txid=%s", ServeAddr, txid)
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
	url := fmt.Sprintf("%s/rawtx?txid=%s", ServeAddr, txid)
	rsp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	res := struct {
		Rawtx string `json:"rawtx"`
	}{}
	if err := json.NewDecoder(rsp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res.Rawtx, nil
}

func (gw *Gateway) InjectTx(rawtx string) (string, error) {
	return BroadcastTx(rawtx)
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
