package mobile

import (
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

func (bn btcNode) GetBalance(addr string) (uint64, error) {
	// get uxout of the address
	_, s := cipher.GenerateKeyPair()
	sknet.SetKey(s.Hex())

	req := pp.GetUtxoReq{
		CoinType:  pp.PtrString("bitcoin"),
		Addresses: []string{addr},
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
