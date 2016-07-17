# Skycoin Exchange

## Installation
```
go get github.com/skycoin/skycoin-exchange
```
## Running server
```
cd skycoin-exchange/cmd/server
go run main.go
```
Default server port is 8080, run the following command to change the port to 8081

```
go run main.go -port=8081
```
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
