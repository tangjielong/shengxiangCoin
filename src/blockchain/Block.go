package blockchain

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"time"
)

type Block struct {
	Timestamp    int64
	Transactions []*Transaction
	PreBlockHash []byte
	Hash         []byte
	Nonce        int
}

func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func (b *Block) HashTransactions() []byte {
	var transactions [][]byte

	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := NewMerkleTree(transactions)
	return mTree.RootNode.Data
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		panic(err)
	}
	return &block
}
func (b *Block) InfoMap() map[string]interface{} {
	info := make(map[string]interface{})
	info["Hash"] = hex.EncodeToString(b.Hash)
	info["PreBlockHash"] = hex.EncodeToString(b.PreBlockHash)
	info["Timestamp"] = time.Unix(b.Timestamp, 0)
	var Transactions []map[string]interface{}
	for _, tx := range b.Transactions {
		Transactions = append(Transactions, tx.InfoMap())
	}
	info["Transactions"] = Transactions
	return info
}
