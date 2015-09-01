package skycoin_exchange

import (
	"github.com/skycoin/skycoin/src/cipher"
)

/*
When sending request, client includes
- their address (identifier)
- hash of message
- signs hash with their private key
*/

//send this with every request
//
type MsgAuth struct {
	Address   cipher.Address
	Hash      cipher.SHA256 //hash of the JSON message
	Signature cipher.Sig
}

//check user authentication for request
func CheckMsgAuth(a MsgAuth, msg []byte) error {
	hash := cipher.SumSHA256(msg)
	if hash != a.Hash {
		return errors.error("hash does not match")
	}
	//check if pubkey can be recovered from the signature
	err := cipher.VerifySignedHash(a.Signature, a.Hash)
	if err != nil {
		return err
	}
	//check signature
	err = cipher.ChkSig(a.Address, a.Hash, a.Signatures)
	if err != nil {
		return nil
	}
}

//creates user authentication for json request
func CreateMsgAuth(seckey cipher.SecKey, msg []byte) (MsgAuth, error) {
	hash := cipher.SumSHA256(msg)
	//func SignHash(hash SHA256, sec SecKey) Sig {
	sig := cipher.SignHash(hash, seckey)
	addr := cipher.AddressFromSecKey(seckey)

	auth := MsgAuth{Address: addr, Signature: sig, Hash: hash}

	err := CheckMsgAuth(auth)
	if err != nil {
		return MsgAuth{}, err
	}

	return auth
}
