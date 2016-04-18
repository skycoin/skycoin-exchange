package skyclient

import (
	"github.com/skycoin/skycoin/src/wallet"
)

type Client struct {
	Wallets wallet.Wallets
}

func New() (*Client, error) {
	client := Client{}
	var err error
	client.Wallets, err = wallet.LoadWallets("wallets/")
	return &client, err
}

func (self *Client) CreateNewAccount(accountName string) (string, error) {
	return accountName, nil
}

func (self *Client) GetAccountAddress(accountName string) (string, error) {
	return accountName, nil
}

func (self *Client) SendFrom(accountName string, addressDestination string) (string, error) {
	return accountName, nil
}

func (self *Client) GetBalance(accountName string) (int, error) {
	return 0, nil
}
