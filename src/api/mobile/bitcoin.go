package mobile

import (
	"github.com/skycoin/skycoin-exchange/src/coin"
	bitcoin "github.com/skycoin/skycoin-exchange/src/coin/bitcoin"
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/sknet"
	"github.com/skycoin/skycoin/src/cipher"
)

type btcNode struct {
	NodeAddr string
}

func (bn btcNode) ValidateAddr(address string) error {
	_, err := cipher.BitcoinDecodeBase58Address(address)
	return err
}

func (bn btcNode) GetBalance(addrs []string) (uint64, error) {
	// get uxout of the address
	_, s := cipher.GenerateKeyPair()
	sknet.SetKey(s.Hex())

	req := pp.GetUtxoReq{
		CoinType:  pp.PtrString("bitcoin"),
		Addresses: addrs,
	}
	res := pp.GetUtxoRes{}
	if err := sknet.EncryGet(bn.NodeAddr, "/auth/get/utxos", req, &res); err != nil {
		return 0, err
	}
	var bal uint64
	for _, u := range res.BtcUtxos {
		bal += u.GetAmount()
	}

	return bal, nil
}

func (bn btcNode) CreateRawTx(txIns []coin.TxIn, keys []cipher.SecKey, txOuts interface{}) (string, error) {
	gw := bitcoin.Gateway{}
	return gw.CreateRawTx(txIns, txOuts)
}

func (bn btcNode) BroadcastTx(rawtx string) (string, error) {
	gw := bitcoin.Gateway{}
	return gw.InjectTx(rawtx)
}

func (bn btcNode) PrepareTx(addrs []string, toAddr string, amt uint64) ([]coin.TxIn, []string, interface{}, error) {
	return nil, nil, nil, nil
}
