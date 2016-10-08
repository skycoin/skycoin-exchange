package mobile

import (
	"github.com/skycoin/skycoin-exchange/src/pp"
	"github.com/skycoin/skycoin-exchange/src/sknet"
	"github.com/skycoin/skycoin/src/cipher"
)

type skyNode struct {
	NodeAddr string
}

func (sn skyNode) GetBalance(addr string) (uint64, error) {
	// get uxout of the address
	_, s := cipher.GenerateKeyPair()
	sknet.SetKey(s.Hex())

	req := pp.GetUtxoReq{
		CoinType:  pp.PtrString("skycoin"),
		Addresses: []string{addr},
	}
	res := pp.GetUtxoRes{}
	if err := sknet.EncryGet(sn.NodeAddr, "/auth/get/utxos", req, &res); err != nil {
		return 0, err
	}
	var bal uint64
	for _, u := range res.SkyUtxos {
		bal += u.GetCoins()
	}

	return bal, nil
}

func (sn skyNode) ValidateAddr(address string) error {
	_, err := cipher.DecodeBase58Address(address)
	return err
}
