package skycoin_exchange

import (
	"github.com/skycoin/skycoin/src/aether/encoder"
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
	Signature cipher.Sig    //set to zero
	Type      uint16
	Seq       uint64 //nonce, unix time, nanoseconds?
	Msg       []byte
}

func (self MsgAuth) CalcHash() cipher.SHA256 {
	self.Signature = cipher.Signature{} //zero out
	b1 := encoder.Serialize(self)       //body
	return cipher.SumSHA256(b1)
}

//check user authentication for request
func CheckMsgAuth(a MsgAuth, msg []byte) error {
	hash := a.CalcHash()
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
//func CreateMsgAuth(seckey cipher.SecKey, msg []byte) (MsgAuth, error) {
func CreateMsgAuth(seckey cipher.SecKey, a MsgAuth) (MsgAuth, error) {
	a.Address = cipher.AddressFromSecKey(seckey)
	a.Hash = cipher.CalcHash(msg)
	A.Signature = cipher.SignHash(hash, seckey)

	err := CheckMsgAuth(auth)
	if err != nil {
		return MsgAuth{}, err
	}

	return auth
}
