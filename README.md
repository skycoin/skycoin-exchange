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
url: /api/v1/coins?id=[:id]&key=[:key]
params:
	id: account id.
	key: account key.
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
mode: POST
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
url: /api/v1/account/order?id=[:id]&key=[:key]
params:
	id: account id.
	key: account key.
request json:
{
   "type": "bid", // bid or ask
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
url: /api/v1/orders/[:type]?coin_pair=[:coin_pair]&start=[:start]&end=[:end]&id=[:id]&key=[:key]
params:
	id: account id.
	key: account key.
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
#### get utxos
```
mode: GET
url: /api/v1/utxos?coin_type=[:coin_type]&addrs=[:addrs]&id=[:id]&key=[:key]
params:
coin_type: coin type, can be bitcoin or skycoin.
addrs: addresses, joined with `,`.
id: account id.
key: account key.
response json:
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "btc_utxos": [
    {
      "address": "1EknG7EauSW4zxFtSrCQSHe5PJenks55s6",
      "txid": "c5ab911556a4628a5a98ee5386a8a3b465831c66953d288bbfeb4221e95158d8d",
      "vout": 0,
      "amount": 90000
    },
    {
      "address": "1EknG7EauSW4zxFtSrCQSHe5PJenks55s6",
      "txid": "a9a1ef0525b1446232fcb69bb4ef99ef239f78a7046f784b972f22a60348d963",
      "vout": 0,
      "amount": 90000
    }
  ]
}
```
#### get transaction
```
mode: GET
url: /api/v1/tx?id=[:id]&key=[:key]&coin_type=[:coin_type]&txid=[:txid]
params:
	id: account id
	key: account key
	coin_type: bitcoin or skycoin
	txid: transaction id
response json:
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "coin_type": "bitcoin",
  "tx": {
    "value": "{\"txid\":\"5756ff16e2b9f881cd15b8a7e478b4899965f87f553b6210d0f8e5bf5be7df1d\",\"version\":1,\"locktime\":981825022,\"vin\":[{\"coinbase\":\"03a6ab05e4b883e5bda9e7a59ee4bb99e9b1bc76a3a2bb0e9c92f06e4a6349de9ccc8fbe0fad11133ed73c78ee12876334c13c02000000f09f909f2f4249503130302f4d696e65642062792073647a686162636400000000000000000000000000000000\",\"txid\":\"\",\"vout\":0,\"scriptSig\":null,\"sequence\":2765846367}],\"vout\":[{\"value\":\"25.37726812\",\"n\":0,\"scriptPubKey\":{\"asm\":\"OP_DUP OP_HASH160 c825a1ecf2a6830c4401620c3a16f1995057c2ab OP_EQUALVERIFY OP_CHECKSIG\",\"hex\":\"76a914c825a1ecf2a6830c4401620c3a16f1995057c2ab88ac\",\"type\":\"pubkeyhash\",\"addresses\":[\"1KFHE7w8BhaENAswwryaoccDb6qcT6DbYY\"]}}],\"blockhash\":\"0000000000000000027d0985fef71cbc05a5ee5cdbdc4c6baf2307e6c5db8591\",\"confirmations\":54117,\"time\":1440604784,\"blocktime\":1440604784}"
  }
}
```
#### get raw transaction
```
mode: GET
url: /api/v1/rawtx?id=[:id]&key=[:key]&coin_type=[:coin_type]&txid=[:txid]
params:
	id: account id
	key: account key
	coin_type: bitcoin or skycoin
	txid: transaction id
response json:
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "coin_type": "bitcoin",
  "rawtx": "010000000132ea3894c4b2c68bb1255be5d0e8a26bd336fd7a562eca9f7c435c9268199f06020000006b483045022100dd4e1b960726e3d3d205cb5ef4d92b3e04f3839757606800ed662069a841ffdc02203f68723bbbf9800d16555ace1ef2f46e65c2a6341643f3c5bf84158b108e6d5d012103eb8b81f8ebc988c61d3cc4c4ac3d546b02a4994d612725e91d8d69a72045fb18ffffffff019d3b1f020000000017a914bfc03379d17dd1e918a026b76cde472bea7ac7268700000000"
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
