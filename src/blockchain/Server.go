package blockchain

import (
	"encoding/json"
)

func CreateBlockchain(address string) {
	bc := CreateBlockChain(address)
	defer bc.Close()

	UTXOSet := UTXOSet{bc}
	UTXOSet.ReIndex()
}

func PrintChain() string {
	bc := GetBlockChain()
	defer bc.Close()
	info := bc.InfoMap()
	b, err := json.Marshal(info)
	PanicIfError(err)
	return string(b)
}

func ReindexUTXO() {
	bc := GetBlockChain()
	UTXOSet := UTXOSet{bc}
	UTXOSet.ReIndex()
}
