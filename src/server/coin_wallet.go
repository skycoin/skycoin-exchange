package skycoin_exchange

// import (
// 	"fmt"
// 	"os/exec"
//
// 	"github.com/skycoin/skycoin/src/cipher"
// )
//
// type CoinType int8
//
// const (
// 	Bitcoin = iota
// 	Skycoin
// 	// Shellcoin
// 	// Ethereum
// 	// other coins...
// )
//
// // Wallet, generate and store addresses for various coin types.
// type Wallet struct {
// 	ID          string                      // wallet id
// 	Seed        string                      // used to generate address
// 	Type        string                      // default: deterministic
// 	CoinAddress map[CoinType][]AddressEntry // key is coin type, value is address list.
// }
//
// type AddressEntry struct {
// 	Coin    CoinType // coin type
// 	Address string   // address
// 	Pubkey  string   // publich key
// 	Seckey  string   // private key
// }
//
// func CreateWallet() Wallet {
// 	return Wallet{
// 		Type:        "deterministic",
// 		Seed:        cipher.SumSHA256(cipher.RandByte(1024)).Hex(),
// 		CoinAddress: make(map[CoinType][]AddressEntry)}
// }
//
// // GenerateAddress, generate new address base on the seed and coin type, and then store the address.
// func (self *Wallet) NewAddresses(seed string, coinType CoinType, num int32) error {
// 	// use cmd line tool to generate addresses.
// 	switch coinType {
// 	case Bitcoin:
// 		cmd := exec.Command("address_gen", "-b", fmt.Sprintf("-n=%d", num), fmt.Sprintf("-seed=\"%s\"", seed))
// 	case Skycoin:
// 		cmd := exec.Command("address_gen", fmt.Sprintf("-n=%d", num), fmt.Sprintf("-seed=\"%s\"", seed))
// 	}
// 	return nil
// }
