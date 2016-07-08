package account

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	bitcoin "github.com/skycoin/skycoin-exchange/src/server/coin_interface/bitcoin"
	"github.com/skycoin/skycoin-exchange/src/server/wallet"
	"github.com/skycoin/skycoin/src/cipher"
)

type AccountID cipher.PubKey

type Accounter interface {
	GetWalletID() string                                                                             // return the wallet id.
	GetID() AccountID                                                                                // return the account id.
	GetNewAddress(ct wallet.CoinType) string                                                         // return new address for receiveing coins.
	GetBalance(ct wallet.CoinType) uint64                                                            // return the account's balance.
	GetAddressBalance(addr string) (uint64, error)                                                   // return the address's balance
	GenerateWithdrawlTx(ct wallet.CoinType, Amount uint64, toAdd string, fee uint64) ([]byte, error) // account generate withdraw transaction.
	GetAddressEntries(coinType wallet.CoinType) ([]wallet.AddressEntry, error)
}

// ExchangeAccount maintains the account state
type ExchangeAccount struct {
	ID          AccountID                  // account id
	balance     map[wallet.CoinType]uint64 // the balance should not be accessed directly.
	wltID       string                     // wallet used to maintain the address, UTXOs, balance, etc.
	balance_mtx sync.RWMutex               // mutex used to protect the balance's concurrent read and write.
	wlt_mtx     sync.Mutex                 // mutex used to protect the wallet's conncurrent read and write.
}

type addrBalance struct {
	Addr    string
	Balance uint64
}

type byBalance []addrBalance

func (bb byBalance) Len() int           { return len(bb) }
func (bb byBalance) Swap(i, j int)      { bb[i], bb[j] = bb[j], bb[i] }
func (bb byBalance) Less(i, j int) bool { return bb[i].Balance < bb[j].Balance }

// newExchangeAccount helper function for generating and initialize ExchangeAccount
func newExchangeAccount(id AccountID, wltID string) ExchangeAccount {
	return ExchangeAccount{
		ID:    id,
		wltID: wltID,
		balance: map[wallet.CoinType]uint64{
			wallet.Bitcoin: 0,
			wallet.Skycoin: 0,
		}}
}

func (self ExchangeAccount) GetID() AccountID {
	return self.ID
}

func (self ExchangeAccount) GetWalletID() string {
	return self.wltID
}

// GetNewAddress generate new address for this account.
func (self *ExchangeAccount) GetNewAddress(ct wallet.CoinType) string {
	// get the wallet.
	wlt, err := wallet.GetWallet(self.wltID)
	if err != nil {
		panic(fmt.Sprintf("account get wallet faild, wallet id:%s", self.wltID))
	}

	self.wlt_mtx.Lock()
	defer self.wlt_mtx.Unlock()
	addr, err := wlt.NewAddresses(ct, 1)
	if err != nil {
		panic(err)
	}
	return addr[0].Address
}

// Get the current recored balance.
func (self *ExchangeAccount) GetBalance(coinType wallet.CoinType) uint64 {
	self.balance_mtx.RLock()
	defer self.balance_mtx.RUnlock()
	return self.balance[coinType]
}

func (self ExchangeAccount) GetAddressBalance(addr string) (uint64, error) {
	return bitcoin.GetBalance([]string{addr})
}

// SetBalance update the balanace of specific coin.
func (self *ExchangeAccount) setBalance(coinType wallet.CoinType, balance uint64) error {
	self.balance_mtx.Lock()
	defer self.balance_mtx.Unlock()
	if _, ok := self.balance[coinType]; !ok {
		return fmt.Errorf("the account does not have %s", coinType)
	}
	self.balance[coinType] = balance
	return nil
}

// GenerateWithdrawTx
func (self ExchangeAccount) GenerateWithdrawlTx(coinType wallet.CoinType, amount uint64, toAddr string, fee uint64) ([]byte, error) {
	// check if balance sufficient
	bla := self.GetBalance(coinType)

	if bla < (amount + fee) {
		return []byte{}, errors.New("balance is not sufficient")
	}

	// choose the appropriate utxos。
	utxos, err := chooseUtxos(&self, coinType, amount)
	if err != nil {
		return []byte{}, err
	}

	// create change address
	changeAddr := self.GetNewAddress(coinType)

	outAddrs := []bitcoin.UtxoOut{
		bitcoin.UtxoOut{Addr: toAddr, Value: amount},
		bitcoin.UtxoOut{Addr: changeAddr, Value: bla - amount - fee},
	}

	switch coinType {
	case wallet.Bitcoin:
		tx, err := bitcoin.NewTransaction(utxos, outAddrs)
		if err != nil {
			return []byte{}, errors.New("create bitcoin transaction failed")
		}
		return bitcoin.DumpTxBytes(tx), nil
	default:
		return []byte{}, errors.New("unknow coin type")
	}

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

func chooseUtxos(a Accounter, coinType wallet.CoinType, amount uint64) ([]bitcoin.UtxoWithkey, error) {
	addrEntries, err := a.GetAddressEntries(coinType)
	utxoks := []bitcoin.UtxoWithkey{}
	if err != nil {
		return utxoks, errors.New("get account addresses failed")
	}

	addrBals := map[string]uint64{} // key: address, value: balance
	addrKeys := map[string]string{} // key: address, value: private key
	balList := []addrBalance{}

	for _, addrEntry := range addrEntries {
		// get the balance of addr
		b, err := a.GetAddressBalance(addrEntry.Address)
		if err != nil {
			return utxoks, err
		}
		addrBals[addrEntry.Address] = b
		addrKeys[addrEntry.Address] = addrEntry.Secret
		balList = append(balList, addrBalance{Addr: addrEntry.Address, Balance: b})
	}

	// sort the bals list
	sort.Sort(byBalance(balList))

	return []bitcoin.UtxoWithkey{}, nil

}

// chooseUtxos choose utxos that can satisify the amount, and set the private key as well.
// func (self ExchangeAccount) chooseUtxos(coinType wallet.CoinType, amount uint64) ([]bitcoin.UtxoWithkey, error) {
// }

func (self ExchangeAccount) GetAddressEntries(coinType wallet.CoinType) ([]wallet.AddressEntry, error) {
	// get address list of this account
	wlt, err := wallet.GetWallet(self.wltID)
	if err != nil {
		return []wallet.AddressEntry{}, fmt.Errorf("account get wallet faild, wallet id:%s", self.wltID)
	}
	addresses := wlt.GetAddressEntries(coinType)
	return addresses, nil
}
