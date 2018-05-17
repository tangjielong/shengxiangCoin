package blockchain

import (
	"bytes"
	"encoding/hex"
)

type TxInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

func (in *TxInput) UseKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func (in *TxInput) InfoMap() map[string]interface{} {
	info := make(map[string]interface{})
	info["Txid"] = hex.EncodeToString(in.Txid)
	info["Vout"] = in.Vout
	return info
}
