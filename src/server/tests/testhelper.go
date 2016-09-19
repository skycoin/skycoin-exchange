package tests

// import (
// 	"errors"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/skycoin/skycoin-exchange/src/server"
// 	"github.com/skycoin/skycoin-exchange/src/server/account"
// 	bitcoin "github.com/skycoin/skycoin-exchange/src/server/coin/bitcoin"
// 	"github.com/skycoin/skycoin-exchange/src/server/engine"
// 	"github.com/skycoin/skycoin-exchange/src/server/wallet"
// 	"github.com/skycoin/skycoin/src/cipher"
// )

// // CaseHandler represents one test case, which will be invoked by MockServer.
// type CaseHandler func() (*httptest.ResponseRecorder, *http.Request)

// // MockServer mock server state for various test cases,
// // people can fake the server's state by fullfill the Server interface, and
// // define various request cases by defining functions that match the signature of
// // CaseHandler.
// func MockServer(egn engine.Exchange, fs CaseHandler) *httptest.ResponseRecorder {
// 	gin.SetMode(gin.TestMode)
// 	router := server.NewRouter(egn)
// 	w, r := fs()
// 	router.ServeHTTP(w, r)
// 	return w
// }

// // HttpRequestCase is used to create REST api test cases.
// func HttpRequestCase(method string, url string, body io.Reader) CaseHandler {
// 	return func() (*httptest.ResponseRecorder, *http.Request) {
// 		w := httptest.NewRecorder()
// 		r, err := http.NewRequest(method, url, body)
// 		if err != nil {
// 			panic(err)
// 		}
// 		switch method {
// 		case "POST":
// 			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 		}
// 		return w, r
// 	}
// }

// // FakeAccount for mocking various account state.
// type FakeAccount struct {
// 	ID      string
// 	Addr    string
// 	Balance uint64
// }

// // FakeServer for mocking various server status.
// type FakeServer struct {
// 	A      account.Accounter
// 	Seckey cipher.SecKey
// 	Fee    uint64
// }

// func (fa FakeAccount) GetID() account.AccountID {
// 	d, err := cipher.PubKeyFromHex(fa.ID)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return account.AccountID(d)
// }

// func (fa FakeAccount) GetNewAddress(ct coin.Type) string {
// 	return fa.Addr
// }

// func (fa FakeAccount) GetBalance(ct coin.Type) uint64 {
// 	return fa.Balance
// }

// func (fa FakeAccount) GenerateWithdrawlTx(ct coin.Type, Amount uint64, toAdd string, fee uint64) ([]byte, error) {
// 	return []byte{}, nil
// }

// func (fa FakeAccount) GetAddressBalance(addr string) (uint64, error) {
// 	return uint64(0), nil
// }

// func (fa FakeAccount) GetAddressEntries(coinType coin.Type) ([]wallet.AddressEntry, error) {
// 	return []wallet.AddressEntry{}, nil
// }

// func (fa *FakeAccount) AddDepositAddress(ct coin.Type, addr string) {

// }

// func (fa *FakeAccount) DecreaseBalance(ct coin.Type, amt uint64) error {
// 	return nil
// }

// func (fa *FakeAccount) IncreaseBalance(ct coin.Type, amt uint64) error {
// 	return nil
// }

// func (fs *FakeServer) CreateAccountWithPubkey(pk cipher.PubKey) (account.Accounter, error) {
// 	fs.A = &FakeAccount{ID: pk.Hex()}
// 	return fs.A, nil
// }

// func (fs *FakeServer) GetAccount(id account.AccountID) (account.Accounter, error) {
// 	if fs.A != nil && fs.A.GetID() == id {
// 		return fs.A, nil
// 	}
// 	return nil, errors.New("account not found")
// }

// func (fs *FakeServer) Run() {

// }

// func (fs *FakeServer) PutUtxos(ct coin.Type, utxos []bitcoin.UtxoWithkey) {

// }

// func (fs FakeServer) GetBtcFee() uint64 {
// 	return fs.Fee
// }

// func (fs FakeServer) GetServPrivKey() cipher.SecKey {
// 	return fs.Seckey
// }
// func (fs *FakeServer) AddWatchAddress(ct coin.Type, addr string) {

// }

// func (fs *FakeServer) ChooseUtxos(coinType coin.Type, amount uint64, tm time.Duration) (interface{}, error) {
// 	return []bitcoin.UtxoWithkey{}, nil
// }

// func (fs *FakeServer) GetNewAddress(ct coin.Type) string {
// 	return ""
// }
