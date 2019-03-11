package main

// TXOutput represents a transaction output.
type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

// Lock signs the output.
func (out *TXOutput) Lock(address []byte)  {
	pubKeyHash: = Base58Decode(address)
	out.PubKeyHash = pubKeyHash[1:len(pubKeyHash)-4]
}

// IsLockedWithKey checks if the output can be used by the owner of the pubkey.
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}