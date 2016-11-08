# Skycoin Exchange

## Installation

``` bash
go get github.com/skycoin/skycoin-exchange
```

## Sync the dependencies

```bash
go get github.com/robfig/glock
glock sync github.com/skycoin/skycoin-exchange
```

## Running server

``` bash
cd skycoin-exchange/cmd/server
go run main.go -seed=$seed -skycoin_node_addr=$skycoin_node_address
```

The `seed` flag must be specificed, server will generate wallet base on it.
The default server port is 8080, and you can use the `port` flag to change it.

As the exchange will comunicate with skycoin node, we use the `skycoin_node_addr`
flag to make it configurable. The default value of `127.0.0.1:6420` will be used
if it's not set.

## Setup admin in server <a id="setup-admin"></a>

As some apis need admin privilege, the server do not have admin account by default，use the following command to set up admin accounts.

``` bash
go run main.go -admins="0311ff3ed447e3ebe176e929017556e2d2be7c52b1f241dd80df98635ea9f53b22"
```

Multiple admin accounts are connected by commas.

## Help

For more usage, run the help command:

``` bash
go run main.go --help
```

## API client

### Running the client

``` bash
./client.sh
```

Default client port is 6060, use the `p` flag to change it.

### Create account

* mode: POST
* url: /api/v1/accounts

response json:

``` json
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "pubkey": "03c9852d11d2d84d23356df6390f10a53dbd8026ed4bc58f9ccddd7a7c69e00715",
  "created_at": 1470188576
}
```

Once the new account is created, this account will become the `active account`, that means the following apis calls are all base on this account.

### Get account

This api is used to get account in client, you can use it to list all the acccounts, or to get the active account.

* mode: GET
* url: /api/v1/account?active=[:active]
* params:
  * active: optional condition, if not set, return all accounts, otherwise the active must be 1 and return the active account.

response json:

``` json
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "accounts": [
    {
      "pubkey": "02c9656e65f70753f021832a7a1874c966917974b242b11b2d73d04bcaaea21a4d",
      "wallet_ids": {
        "bitcoin": "bitcoin_myf"
      }
    }
  ]
}
```

### Switch account

* mode: PUT
* url: /api/v1/account/state?pubkey=[:pubkey]

response json:

``` json
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  }
}
```

### Get supported coins

* mode:GET
* url: /api/v1/coins

response json:

``` json
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

### Get deposit address

* mode: POST
* url: /api/v1/account/deposit_address?coin_type=[:coin_type]
* params:
  * coin_type: can be bitcoin, skycoin, etc.

response json:

``` json
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "coin_type": "bitcoin",
  "address": "1HBuSp1G151xTqLpMT3mBDXskC5iVNTAwx"
}
```

### Get account balance

* mode: GET
* url: /api/v1/account/balance?coin_type=[:coin_type]
* params:
  * coin_type: can be bitcoin, skycoin.

response json:

``` json
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "coin_type": "bitcoin",
  "balance": 480000
}
```

### Withdraw coins

* mdoe: POST
* url: /api/v1/account/withdrawal?coin_type=[:type]&amount=[:amt]&toaddr=[:toaddr]
* params:
  * coin_type: can be bitcoin, skycoin, etc.
  * amount: the coin number you want to withdrawal, btc in satoshis, sky in drops.
  * toaddr: address you want to receive the coins.

response json:

``` json
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "new_txid": "21b1a9c59a3a631f14b7f91c9b886f6e379c36dd357f7628964107c4d953ea5a"
}
```

### Create order

* mode: POST
* url: /api/v1/account/order?coin_pair=[:coin_pair]&type=[:type]&price=[:price]&amt=[:amt]
* params:
  * coin_pair: coin pair, like bitcoin/skycoin.
  * type: order type, can be bid or ask
  * price: price
  * amt: amount

response json:

``` json
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "order_id": 8
}
```

### Get orders

* mode: GET
* url: /api/v1/orders/[:type]?coin_pair=[:coin_pair]&start=[:start]&end=[:end]
* params:
  * type: order type, can be bid or ask.
  * coin_pair: coin pair, joined by '/', like: bitcoin/skycoin.
  * start: start index of the orders.
  * end: end index of the orders.

response json:

``` json
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

### Get utxos

* mode: GET
* url: /api/v1/utxos?coin_type=[:coin_type]&addrs=[:addrs]
* params:
  * coin_type: coin type, can be bitcoin or skycoin.
  * addrs: addresses, connected with comma.

response json:

``` json
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

NOTE: the utxos response json struct of bitcoin is different from skycoin.

### Get transaction

This api is used to get tx verbose info by txid.

* mode: GET
* url: /api/v1/tx?coin_type=[:coin_type]&txid=[:txid]
* params:
  * coin_type: bitcoin or skycoin
  * txid: transaction id

response json:

``` json
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "coin_type": "bitcoin",
  "tx": {
    "btc": {
      "txid": "5756ff16e2b9f881cd15b8a7e478b4899965f87f553b6210d0f8e5bf5be7df1d",
      "version": 1,
      "locktime": 981825022,
      "vin": [
        {
          "coinbase": "03a6ab05e4b883e5bda9e7a59ee4bb99e9b1bc76a3a2bb0e9c92f06e4a6349de9ccc8fbe0fad11133ed73c78ee12876334c13c02000000f09f909f2f4249503130302f4d696e65642062792073647a686162636400000000000000000000000000000000",
          "sequence": 2765846367
        }
      ],
      "vout": [
        {
          "value": "25.37726812",
          "n": 0,
          "scriptPubkey": {
            "asm": "OP_DUP OP_HASH160 c825a1ecf2a6830c4401620c3a16f1995057c2ab OP_EQUALVERIFY OP_CHECKSIG",
            "hex": "76a914c825a1ecf2a6830c4401620c3a16f1995057c2ab88ac",
            "type": "pubkeyhash",
            "addresses": [
              "1KFHE7w8BhaENAswwryaoccDb6qcT6DbYY"
            ]
          }
        }
      ],
      "blockhash": "0000000000000000027d0985fef71cbc05a5ee5cdbdc4c6baf2307e6c5db8591",
      "confirmations": 54245,
      "time": 1440604784,
      "blocktime": 1440604784
    }
  }
}
```

NOTE: the bitcoin transaction response struct is differen from skycoin.

### Get raw transaction

* mode: GET
* url: /api/v1/rawtx?coin_type=[:coin_type]&txid=[:txid]
* params:
  * coin_type: bitcoin or skycoin
  * txid: transaction id

response json:

``` json
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

### Update balance

This api is used to update balance of specific account, need admin privilege, to acquire admin privilege, see [setup admin in server](#setup-admin).

* mode: PUT
* url: /api/v1/admin/account/balance?coin_type=[:coin_type]&dst=[:dst]&amt=[:amt]
* params:
  * coin_type: bitcoin or skycoin
  * dst: account pubkey, the account whose balance is going to be updated.
  * amt: balance amount

response json:

``` json
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  }
}
```

### Create wallet

* mode: POST
* url: /api/v1/wallet?type=[:type]&seed=[:seed]
* params:
  * type: wallet type, can be bitcoin or skycoin
  * seed: wallet seed

response json:

``` json
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "id": "bitcoin_sd1101"
}
```

### Generate new address

* mode: POST
* url: /api/v1/wallet/address?id=[:id]
* params:
  * id: wallet id

response json:

``` json
{
  "Result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "address": "1FhzNKvpStmS4ZwpiZwRNVhTQvZBa39VNA"
}
```

### Get public and secret key pair

* mode: GET
* url: /api/v1/wallet/address/key?id=[:id]&address=[:address]
* params:
  * id: wallet id
  * address: coin address

response json:

``` json
{
  "Result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "pubkey": "03712c6bf0601f7a663ad7812f8d031e3e3f07f3f7ed03ad165dd7ee28120e7102",
  "seckey": "L4Y4E6UdCFqmmcYNGBFAXBZif6pfoox5CQNReBsnX8aimgBcRYeX"
}
```

### Get wallet balance

* mode: GET
* url: /api/v1/wallet/balance?id=[:id]

``` json

response json:
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "balance": {
    "amount": 11000000,
    "hours": 1518046
  }
}
```

The above response json is a skycoin wallet's balance result, you can see the balance contains `hours` field, while if bitcoin, the `hours` field will be omited.

### Create raw transaction

* mode: POST
* url: /api/v1/create_rawtx?coin_type=[:coin_type]
* params:
  * coin_type: skycoin or bitcoin

Specify the tx in and outs in json request body as below:

bitcoin request:

``` json
{
    "tx_ins":
    [
        {
            "txid":"44051e627966fad80f5e97890da1c67148f312cd2b617cf788f274f89abf16a4",
            "vout": 0
        }],
    "tx_outs":
    [
        {
          "address":"1GpMHk4GcWRBZfmBBd2iCt4PinhNs4SPrn",
          "value":4000
        },
        {
            "address":"1EknG7EauSW4zxFtSrCQSHe5PJenkn55s6",
            "value": 4000
        }]
}
```

skycoin request:

``` json
{
    "tx_ins":
    [
        {
            "txid":"11ad2877281d541e68a5e3004cccd166d85c9edf252cabfa5bb540648380cea9"
        }],
    "tx_outs":
    [
        {
          "address":"5iUowMCu1VSi565Q5DfLNXz26PwVNgN1jd",
          "value":1000000,
          "hours":2
        }]
}
```

### Sign raw transaction

* mode: POST
* url: /api/v1/sign_rawtx?coin_type=[:coin_type]&raw_tx=[:rawtx]
* params:
  * coin_type: skycoin or bitcoin
  * raw_tx: raw transaction that's going to be signed.

response json:

``` json
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "rawtx": "0100000001a416bf9af874f288f77c612bcd12f34871c6a10d89975e0fd8fa6679621e05440000000000ffffffff02a00f0000000000001976a914ad7e5f825191df239d43376d182cf85d3e9ac8a188aca00f0000000000001976a91496e14d971c0a482f37a06ba23094e0cc779676ff88ac00000000"
}
```

Make sure current active account owns the private key to sign the transaction.

### Inject raw transaction

* mode: POST
* url: /api/v1/inject_rawtx?coin_type=[:coin_type]&rawtx=[:rawtx]
* params:
  * coin_type: skycoin or bitcoin
  * rawtx: raw transaction that's going to be injected.

response json:

``` json
{
  "result": {
    "success": true,
    "errcode": 0,
    "reason": "Success"
  },
  "txid": "11ad2877281d541e68a5e3004cccd166d85c9edf252cabfa5bb540648380cea9"
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
