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
  "account_id": "AyzApEpMdFeInmW1jAJ9Yw8a+9fJh2Qeab1zwqAq9euX",
  "created_at": 1469336497
}
```

#### get deposit address
```
mode: GET
url: /api/v1/deposit_address?cointype=[:type]
params:
	type: can be bitcoin, skycoin, etc.
response json:
{
  "account_id": "AyzApEpMdFeInmW1jAJ9Yw8a+9fJh2Qeab1zwqAq9euX",
  "coin_type": "skycoin",
  "address": "mmKNxRvQ6qm78njpT2W9JiRjC26rgd8xzG"
}
```

#### get account balance
```
mode: GET
url: /api/v1/account/balance?cointype=[:type]
params:
	type: can be bitcoin, skycoin, etc.
response json:
{
  "account_id": "AyzApEpMdFeInmW1jAJ9Yw8a+9fJh2Qeab1zwqAq9euX",
  "coin_type": "skycoin",
  "balance": 9000000
}
```

#### withdraw coins
```
mdoe: POST
url: /api/v1/account/withdrawal?cointype=[:type]&amount=[:amt]&toaddr=[:toaddr]
params:
	type: can be bitcoin, skycoin, etc.
	amount: the coin number you want to withdrawal, in satoshis.
	toaddr: address you want to receive the coinsns.
response json:
{
  "account_id": "AyzApEpMdFeInmW1jAJ9Yw8a+9fJh2Qeab1zwqAq9euX",
  "new_txid": "21b1a9c59a3a631f14b7f91c9b886f6e379c36dd357f7628964107c4d953ea5a"
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
