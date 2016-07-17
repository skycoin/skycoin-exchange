

This is prototype exchange for Skycoin/Bitcoin

- Username is public key hashes (Addresses)
- All post requests are signed by public key

URLS
- deposit bitcoin
- withdrawl bitcoin
- deposit skycoin
- withdrawl skycoin

- Get Balance
- New Bid
- New Ask
- Order Book
- Cancel Order

States
- client balance
- bid/ask book

All responses from the server are json
- the client is run on local host and generates the webpage
- there is a configuration file that stores the account information

Each Bid/Ask has a unique id that is incremented
Bid/Asks are each cleared each tick, which is incrememented

Components
- Server
-- runs order book, takes in events, responds to events, exposes RPC
- Client
-- polls server, exposes JSON/html interface, has local web-interface
- Accounts Manager
-- handles coin withdrawls and deposits. Credits and debits accounts

=== Authentication

Must have public key to identify server
- self-describing RPC

Should operate over IRC like network

=== Statistics

Keep track of
- exchange balance net, in/out
- how many Bitcoin/Skycoin flow into/out of the exchange (total in, total out)
- how many coins flow between each other (capital flows)


## message struct

#### create account
request
{
  pubkey string
}

response
{
  success bool
  accountID string
  createdAt string
}

#### negotiate nonce key
request
{
  pubkey string
}

response
{
  pubkey string
}

#### transfer data
request
{
  accountID string   // account id.
  data []byte // encrypted data with nonce key.
}

response
{
  success bool
  data []byte
}

##### get deposit address
request
{
  accountID string
  {
    coinType string
  }
}

response
{
  success bool
  {
    accountID string
    depositeAddr string
  }
}
