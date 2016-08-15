package model

import "github.com/skycoin/skycoin/src/cipher"

type Client struct {
	ServApiRoot string        // api root
	ServPubkey  cipher.PubKey // server's pubkey.
}
