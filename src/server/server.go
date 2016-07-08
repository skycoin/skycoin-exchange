package server

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/skycoin/skycoin-exchange/src/server/account"
	bitcoin "github.com/skycoin/skycoin-exchange/src/server/coin_interface/bitcoin"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
)

type Server interface {
	Run()
	CreateAccountWithPubkey(pubkey cipher.PubKey) (account.Accounter, error)
	GetAccount(id account.AccountID) (account.Accounter, error)
	GetFee() uint64
	Encrypt(data []byte, pubkey cipher.PubKey, nonce []byte) ([]byte, error)
	Decrypt(data []byte, pubkey cipher.PubKey, nonce []byte) ([]byte, error)
	GetNewAddress(coinType wallet.CoinType) string
}

// Config store server's configuration.
type Config struct {
	Port       int
	Fee        int
	DataDir    string
	WalletName string
	Seed       string
	Seckey     cipher.SecKey
}

/*
	The server gets events from the client and processes them
	- get balance/status
	- get deposit addresses
	- withdrawl bitcoin
	- withdrawl skycoin
	- add bid
	- add ask
	- get order book
*/

type ExchangeServer struct {
	account.AccountManager
	cfg    Config
	wallet wallet.Wallet
	wltMtx sync.RWMutex // mutex for protecting the wallet.
}

// New create new server
func New(cfg Config) Server {
	// init the data dir
	path := util.InitDataDir(cfg.DataDir)

	// set the wallet dir.
	wallet.WltDir = filepath.Join(path, "wallets")

	// get wallet
	var wlt wallet.Wallet
	var err error
	if wallet.IsExist(cfg.WalletName) {
		wlt, err = wallet.Load(cfg.WalletName)
		if err != nil {
			panic("server load walle failed")
		}
	} else {
		wlt, err = wallet.New(cfg.WalletName, wallet.Deterministic, cfg.Seed)
		if err != nil {
			panic("server create wallet failed")
		}
	}
	s := &ExchangeServer{
		cfg:            cfg,
		wallet:         wlt,
		AccountManager: account.NewExchangeAccountManager(),
	}
	return s
}

func (self *ExchangeServer) Run() {
	log.Println("skycoin-exchange server started, port:", self.cfg.Port)

	r := NewRouter(self)

	r.Run(fmt.Sprintf(":%d", self.cfg.Port))
}

func (self ExchangeServer) GetFee() uint64 {
	return uint64(self.cfg.Fee)
}

func (self *ExchangeServer) GetNewAddress(coinType wallet.CoinType) string {
	self.wltMtx.Lock()
	defer self.wltMtx.Unlock()
	addrEntry, err := self.wallet.NewAddresses(coinType, 1)
	if err != nil {
		panic("server get new address failed")
	}
	return addrEntry[0].Address
}

func (self ExchangeServer) Decrypt(data []byte, pubkey cipher.PubKey, nonce []byte) (d []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("server decrypt faild")
		}
	}()

	key := cipher.ECDH(pubkey, self.cfg.Seckey)
	d, err = Encrypt(data, key, nonce)
	return
}

func (self ExchangeServer) Encrypt(data []byte, pubkey cipher.PubKey, nonce []byte) (d []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("server encrypt faild")
		}
	}()

	key := cipher.ECDH(pubkey, self.cfg.Seckey)
	d, err = Decrypt(data, key, nonce)
	return
}

func GenerateWithdrawlTx(svr Server, id account.AccountID, coinType wallet.CoinType, amount uint64, toAddr string) ([]byte, error) {
	return []byte{}, nil
	// check if balance sufficient
	// bla := self.GetBalance(coinType)

	// if bla < (amount + fee) {
	// 	return []byte{}, errors.New("balance is not sufficient")
	// }

	// choose the appropriate utxos。
	// utxos, err := chooseUtxos(svr, coinType, amount)
	// if err != nil {
	// 	return []byte{}, err
	// }
	//
	// // create change address
	// changeAddr := svr.GetNewAddress(coinType)
	//
	// outAddrs := []bitcoin.UtxoOut{
	// 	bitcoin.UtxoOut{Addr: toAddr, Value: amount},
	// 	bitcoin.UtxoOut{Addr: changeAddr, Value: bla - amount - fee},
	// }
	//
	// switch coinType {
	// case wallet.Bitcoin:
	// 	tx, err := bitcoin.NewTransaction(utxos, outAddrs)
	// 	if err != nil {
	// 		return []byte{}, errors.New("create bitcoin transaction failed")
	// 	}
	// 	return bitcoin.DumpTxBytes(tx), nil
	// default:
	// 	return []byte{}, errors.New("unknow coin type")
	// }

	// get utxos
	// utxos := []bitcoin.UtxoWithkey{}
	// switch coinType {
	// case wallet.Bitcoin:
	// 	for _, addrEntry := range addrs {
	// 		us := bitcoin.GetUnspentOutputs(addrEntry.Address)
	// 		usks := make([]bitcoin.UtxoWithkey, len(us))
	// 		for i, u := range us {
	// 			usk := bitcoin.NewUtxoWithKey(u, addrEntry.Secret)
	// 			usks[i] = usk
	// 		}
	// 		utxos = append(utxos, usks...)
	// 	}
	// 	msgTx, err := bitcoin.NewTransaction(utxos, outAddrs)
	// 	if err != nil {
	// 		return []byte{}, errors.New("create bitcoin transaction faild")
	// 	}
	// 	return bitcoin.DumpTxBytes(msgTx), nil
	// default:
	// 	return []byte{}, errors.New("unknow coin type")
	// }
}

func chooseUtxos(svr Server, coinType wallet.CoinType, amount uint64) ([]bitcoin.UtxoWithkey, error) {
	// addrEntries, err := a.GetAddressEntries(coinType)
	// utxoks := []bitcoin.UtxoWithkey{}
	// if err != nil {
	// 	return utxoks, errors.New("get account addresses failed")
	// }
	//
	// addrBals := map[string]uint64{} // key: address, value: balance
	// addrKeys := map[string]string{} // key: address, value: private key
	// balList := []addrBalance{}
	//
	// for _, addrEntry := range addrEntries {
	// 	// get the balance of addr
	// 	b, err := a.GetAddressBalance(addrEntry.Address)
	// 	if err != nil {
	// 		return utxoks, err
	// 	}
	// 	addrBals[addrEntry.Address] = b
	// 	addrKeys[addrEntry.Address] = addrEntry.Secret
	// 	balList = append(balList, addrBalance{Addr: addrEntry.Address, Balance: b})
	// }
	//
	// // sort the bals list
	// sort.Sort(byBalance(balList))

	return []bitcoin.UtxoWithkey{}, nil
}
