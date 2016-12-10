package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/skycoin/skycoin-exchange/src/client/account"
	"github.com/skycoin/skycoin-exchange/src/coin"
	bitcoin "github.com/skycoin/skycoin-exchange/src/coin/bitcoin"
	skycoin "github.com/skycoin/skycoin-exchange/src/coin/skycoin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/sknet"
	"github.com/skycoin/skycoin-exchange/src/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

// InjectTx broadcast transaction.
// mode: POST
// url: /api/v1/inject_rawtx?rawtx=[:rawtx]&coin_type=[:coin_type]
// params:
// 		rawtx: raw tx that's going to be injected.
// 		coin_type: skycoin or bitcoin.
func InjectTx(se Servicer) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var rlt *pp.EmptyRes
		for {
			// get tx
			rawtx := r.FormValue("rawtx")
			if rawtx == "" {
				err := errors.New("rawtx is empty")
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			// get coin type
			cp := r.FormValue("coin_type")
			if cp == "" {
				logger.Error("empty coin type")
				rlt = pp.MakeErrRes(errors.New("empty coin type"))
				break
			}

			req := pp.InjectTxnReq{
				CoinType: pp.PtrString(cp),
				Tx:       pp.PtrString(rawtx),
			}

			var res pp.InjectTxnRes
			if err := sknet.EncryGet(se.GetServAddr(), "/inject/tx", req, &res); err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			sendJSON(w, res)
			return
		}
		sendJSON(w, rlt)
	}
}

// GetTx get verbose transaction info by transacton id.
// mode: GET
// url: /api/v1/tx?coin_type=[:coin_type]&txid=[:txid]
func GetTx(se Servicer) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var rlt *pp.EmptyRes
		for {
			// get coin type
			cp := r.FormValue("coin_type")
			if cp == "" {
				rlt = pp.MakeErrRes(errors.New("no coin type"))
				break
			}

			// get txid
			txid := r.FormValue("txid")
			if txid == "" {
				rlt = pp.MakeErrRes(errors.New("no txid"))
				break
			}
			req := pp.GetTxReq{
				CoinType: pp.PtrString(cp),
				Txid:     pp.PtrString(txid),
			}

			var res pp.GetTxRes
			if err := sknet.EncryGet(se.GetServAddr(), "/get/tx", req, &res); err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			sendJSON(w, res)
			return
		}
		logger.Error(rlt.GetResult().GetReason())
		sendJSON(w, rlt)
	}
}

// GetRawTx get raw tx by txid.
// mode: GET
// url: /api/v1/rawtx?coin_type=[:coin_type]&txid=[:txid]
func GetRawTx(se Servicer) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var rlt *pp.EmptyRes
		for {
			// get coin type
			cp := r.FormValue("coin_type")
			if cp == "" {
				logger.Error("no coin type")
				rlt = pp.MakeErrRes(errors.New("no coin type"))
				break
			}
			// get txid
			txid := r.FormValue("txid")
			if txid == "" {
				logger.Error("no txid")
				rlt = pp.MakeErrRes(errors.New("no txid"))
				break
			}
			req := pp.GetRawTxReq{
				CoinType: pp.PtrString(cp),
				Txid:     pp.PtrString(txid),
			}
			var res pp.GetRawTxRes
			if err := sknet.EncryGet(se.GetServAddr(), "/get/rawtx", req, &res); err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			sendJSON(w, res)
			return
		}
		sendJSON(w, rlt)
	}
}

type rawTxParams struct {
	TxIns  []coin.TxIn `json:"tx_ins"`
	TxOuts []struct {
		Addr  string `json:"address"`
		Value uint64 `json:"value"`
		Hours uint64 `json:"hours"`
	} `json:"tx_outs"`
}

// CreateRawTx create raw tx base on some utxos.
// mode: POST
// url: /api/v1/create_rawtx?coin_type=[:coin_type]
// request json:
// 		different in bitcoin and skycoin.
func CreateRawTx(se Servicer) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var rlt *pp.EmptyRes
	loop:
		for {
			// get coin type
			cp := r.FormValue("coin_type")

			// get request body
			params := rawTxParams{}
			if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			gw, err := coin.GetGateway(cp)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				break
			}

			var rawtx string
			switch cp {
			case bitcoin.Type:
				outs := make([]bitcoin.TxOut, len(params.TxOuts))
				for i, o := range params.TxOuts {
					outs[i].Addr = o.Addr
					outs[i].Value = o.Value
				}
				rawtx, err = coin.CreateRawTx(params.TxIns, outs)
			case skycoin.Type:
				outs := make([]skycoin.TxOut, len(params.TxOuts))
				for i, o := range params.TxOuts {
					addr, err := cipher.DecodeBase58Address(o.Addr)
					if err != nil {
						logger.Error(err.Error())
						rlt = pp.MakeErrRes(err)
						break loop
					}
					outs[i].Address = addr
					outs[i].Coins = o.Value
					outs[i].Hours = o.Hours
				}
				rawtx, err = gw.CreateRawTx(params.TxIns, outs)
			}
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			res := struct {
				Result *pp.Result `json:"result"`
				Rawtx  string     `json:"rawtx"`
			}{
				Result: pp.MakeResultWithCode(pp.ErrCode_Success),
				Rawtx:  rawtx,
			}
			sendJSON(w, &res)
			return
		}
		sendJSON(w, rlt)
	}
}

// SignRawTx sign transaction.
// mode: POST
// url: /api/v1/signr_awtx?coin_type=[:coin_type]&rawtx=[:rawtx]
// params:
// 		coin_type: skycoin or bitcoin.
// 		rawtx: raw transaction.
func SignRawTx(se Servicer) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var rlt *pp.EmptyRes
		for {
			// check coin type
			cp := r.FormValue("coin_type")

			// get raw tx
			rawtx := r.FormValue("rawtx")
			if rawtx == "" {
				err := errors.New("rawtx is empty")
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			gw, err := coin.GetGateway(cp)
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}

			tx, err := gw.SignRawTx(rawtx, getPrivKey(cp))
			if err != nil {
				logger.Error(err.Error())
				rlt = pp.MakeErrRes(err)
				break
			}
			res := struct {
				Result *pp.Result `json:"result"`
				Rawtx  string     `json:"rawtx"`
			}{
				Result: pp.MakeResultWithCode(pp.ErrCode_Success),
				Rawtx:  tx,
			}
			sendJSON(w, &res)
			return
		}
		sendJSON(w, rlt)
	}
}

func getPrivKey(cp string) coin.GetPrivKey {
	return func(addr string) (string, error) {
		a, err := account.GetActive()
		if err != nil {
			return "", err
		}
		wltID := a.WltIDs[cp]
		if wltID == "" {
			return "", fmt.Errorf("does not have %s wallet", cp)
		}

		_, key, err := wallet.GetKeypair(wltID, addr)
		if err != nil {
			return "", err
		}
		return key, nil
	}
}
