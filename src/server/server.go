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
	GetServPrivKey() cipher.SecKey
	GetPrivKey(ct wallet.CoinType, addr string) (string, error)
	GetNewAddress(coinType wallet.CoinType) string
	ChooseUtxos(coinType wallet.CoinType, amount uint64) ([]bitcoin.Utxo, error)
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
	UtxosManager
	cfg    Config
	wallet wallet.Wallet
	wltMtx sync.RWMutex // mutex for protecting the wallet.
}

type UtxosManager interface {
	ChooseUtxos(coinType wallet.CoinType, amount uint64) ([]bitcoin.Utxo, error) // choose appropriate utxos and mark them as spending.
	// AddWatchAddress(coinType wallet.CoinType, addr string)                     // add address which has unspent outputs in he server.
	UpdateUtxos() // update the server's utxos.
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

// GetServPrivKey returnt he sever's private key.
func (self ExchangeServer) GetServPrivKey() cipher.SecKey {
	return self.cfg.Seckey
}

// GetPrivKey return the private key of specific address.
func (self ExchangeServer) GetPrivKey(ct wallet.CoinType, addr string) (string, error) {
	entry, err := self.wallet.GetAddressEntry(ct, addr)
	if err != nil {
		return "", err
	}

	return entry.Secret, nil
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

// func (self *ExchangeServer) ChooseUtxos(coinType wallet.CoinType, amount uint64) ([]bitcoin.UtxoWithkey, error) {

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

// 	return []bitcoin.UtxoWithkey{}, nil
// }

func GenerateWithdrawlTx(svr Server, act account.Accounter, coinType wallet.CoinType, amount uint64, toAddr string) ([]byte, error) {
	bal := act.GetBalance(coinType)
	fee := svr.GetFee()
	if bal < amount+fee {
		return []byte{}, errors.New("balance is not sufficient")
	}

	utxos, err := svr.ChooseUtxos(coinType, amount)
	if err != nil {
		return []byte{}, err
	}

	var totalAmounts uint64
	for _, u := range utxos {
		totalAmounts += u.GetAmount()
	}

	// generate a change address
	chgAddr := svr.GetNewAddress(coinType)
	chgAmt := totalAmounts - fee - amount

	outAddrs := []bitcoin.UtxoOut{
		bitcoin.UtxoOut{Addr: toAddr, Value: amount},
		bitcoin.UtxoOut{Addr: chgAddr, Value: chgAmt},
	}

	tx, err := bitcoin.NewTransaction(utxos, outAddrs)
	if err != nil {
		return []byte{}, err
	}

	return bitcoin.DumpTxBytes(tx), nil
}
