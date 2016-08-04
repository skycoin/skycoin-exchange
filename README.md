# Skycoin Exchange

## Installation
```
go get github.com/skycoin/skycoin-exchange
```

## Running server
```
cd skycoin-exchange/cmd/server
go run main.go -seed="wlt seed name"
```
Default server port is 8080, must specify the seed name.
For more usage, run the help command as below:

```
go run main.go --help
```
## Running rpclient
```
cd skycoin-exchange/cmd/client
go run client.go
```
Default rpclient port is 6060.

### Rpc client APIs
#### create account:
```
mode: POST
url: /api/v1/accounts
response json:
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "account_id": "03c9852d11d2d84d23356df6390f10a53dbd8026ed4bc58f9ccddd7a7c69e00715",
  "key": "36797b7fe0f8bcc96ac4f4110eaab0d1cb32fbe8961a77bf132c0efa02e760a7",
  "created_at": 1470188576
}
```

The rpc client will use the account id and key to communicate with exchange server, most of the following APIs will use the id and key.

#### get supported coins
```
mode:GET
url: /api/v1/coins
response json:
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "coins": [
    "BTC",
    "SKY"
  ]
}
```

#### get deposit address
```
mode: GET
url: /api/v1/account/deposit_address?id=[:id]&key=[:key]&coin_type=[:type]
params:
	id: account id.
	key: account key.
	coin_type: can be bitcoin, skycoin, etc.
response json:
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "account_id": "02169842b50a2f452039d18d7b885e1b99801475489368ddcd58365f135784585c",
  "coin_type": "bitcoin",
  "address": "1HBuSp1G151xTqLpMT3mBDXskC5iVNTAwx"
}
```

#### get account balance
```
mode: GET
url: /api/v1/account/balance?id=[:id]&key=[:key]coin_type=[:type]
params:
	id: account id.
	key: account key.
	coin_type: can be bitcoin, skycoin, etc.
response json:
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "account_id": "02169842b50a2f452039d18d7b885e1b99801475489368ddcd58365f135784585c",
  "coin_type": "bitcoin",
  "balance": 480000
}
```

#### withdraw coins
```
mdoe: POST
url: /api/v1/account/withdrawal?id=[:id]&key=[:key]&coin_type=[:type]&amount=[:amt]&toaddr=[:toaddr]
params:
	id: account id.
	key: account key.
	coin_type: can be bitcoin, skycoin, etc.
	amount: the coin number you want to withdrawal, btc in satoshis, sky in drops.
	toaddr: address you want to receive the coins.
response json:
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "account_id": "02169842b50a2f452039d18d7b885e1b99801475489368ddcd58365f135784585c",
  "new_txid": "21b1a9c59a3a631f14b7f91c9b886f6e379c36dd357f7628964107c4d953ea5a"
}

```

#### create order
```
mode: POST
url: /api/v1/account/order/[:type]?id=[:id]&key=[:key]
params:
	type: order type, can be bid or ask.
	id: account id.
	key: account key.
request json:
{
   "coin_pair":"bitcoin/skycoin",
   "amount":90000,
   "price":25
}
response json:
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "account_id": "02169842b50a2f452039d18d7b885e1b99801475489368ddcd58365f135784585c",
  "order_id": 8
}
```

#### get orders
```
mode: GET
url: /api/v1/orders/[:type]?coin_pair=[:coin_pair]&start=[:start]&end=[:end]
params:
	type: order type, can be bid or ask.
	coin_pair: coin pair, joined by '/', like: bitcoin/skycoin.
	start: start index of the orders.
	end: end index of the orders.
response json:
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "coin_pair": "bitcoin/skycoin",
  "type": "bid",
  "orders": [
    {
      "id": 8,
      "type": "bid",
      "price": 25,
      "amount": 90000,
      "rest_amt": 90000,
      "created_at": 1470193222
    },
    {
      "id": 3,
      "type": "bid",
      "price": 25,
      "amount": 90000,
      "rest_amt": 90000,
      "created_at": 1470152057
    }
  ]
}
```
Dependencies
---

```
go get github.com/robfig/glock
glock sync github.com/skycoin/skycoin-exchange
```

To update dependencies
```
glock save github.com/skycoin/skycoin-exchange ???
```
