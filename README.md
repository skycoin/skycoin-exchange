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
Default server port is 8080, must specific the seed name.
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
#### create account
```
mode: POST
url: /api/v1/accounts
```

#### get deposit address
```
mode: GET
url: /api/v1/deposit_address?cointype=[:type]
params:
	type: can be bitcoin, skycoin, etc.
```

#### get account balance
```
mode: GET
url: /api/v1/account/balance?cointype=[:type]
params:
	type: can be bitcoin, skycoin, etc.
```

#### withdraw coins
```
mdoe: POST
url: /api/v1/account/withdrawal?cointype=[:type]&amount=[:amt]&toaddr=[:toaddr]
params:
	type: can be bitcoin, skycoin, etc
	amount: the coin number you want to withdrawal, in satoshis.
	toaddr: address you want to receive the coins.
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
