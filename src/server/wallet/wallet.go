package wallet

type CoinType int8
type WalletType int8

const (
	Bitcoin = iota
	Skycoin
	// Shellcoin
	// Ethereum
	// other coins...
)

const (
	Deterministic = iota // default wallet type
)

type AddressEntry struct {
	Coin    CoinType // coin type
	Address string   // address
	Pubkey  string   // publich key
	Seckey  string   // private key
}

type Wallet interface {
	SetID(id string)
	GetID() string
	NewAddresses(coinType CoinType, num int) []AddressEntry
	GetBalance(addr string) (string, error)
}
