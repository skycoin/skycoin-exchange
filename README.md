

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