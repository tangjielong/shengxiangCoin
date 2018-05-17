package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

const blockBucket = "blocks"
const dbFile = "hello.db"

type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

func CreateBlockChain(address string) *BlockChain {
	if !ValidateAddress(address) {
		log.Panic("Address is not valid")
	}
	if dbExists() {
		log.Panic("Blockchain already exists.")
	}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	PanicIfError(err)
	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTx(address)
		genesis := NewBlock([]*Transaction{cbtx}, []byte{})
		b, err := tx.CreateBucket([]byte(blockBucket))
		if err != nil {
			return err
		}
		if err = b.Put(genesis.Hash, genesis.Serialize()); err != nil {
			return err
		}
		if err = b.Put([]byte("l"), genesis.Hash); err != nil {
			return err
		}
		tip = genesis.Hash
		return nil
	})
	return &BlockChain{tip, db}
}
func GetBlockChain() *BlockChain {
	if !dbExists() {
		log.Panic("No existing blockchain found. Create one first.")
	}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	PanicIfError(err)
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		tip = b.Get([]byte("l"))
		return nil
	})
	PanicIfError(err)
	return &BlockChain{tip, db}
}

func (bc *BlockChain) MineBlock(transactions []*Transaction) *Block {
	var lastHash []byte
	for _, tx := range transactions {
		if !bc.VerifyTransaction(tx) {
			log.Panic("invalid transaction")
		}
	}
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})
	PanicIfError(err)
	newBlock := NewBlock(transactions, lastHash)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		PanicIfError(err)
		err = b.Put([]byte("l"), newBlock.Hash)
		PanicIfError(err)
		bc.tip = newBlock.Hash
		return nil
	})
	PanicIfError(err)
	return newBlock
}

func (bc *BlockChain) FindUTXO() map[string]TxOutputs {
	UTXOs := make(map[string]TxOutputs)
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()
	for bci.HasNext() {
		block := bci.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		Outputs:
			for outIdx, out := range tx.Vout {
				if spentTXOs[txID] != nil {
					// spentOut 已经消费的输出的索引
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				// 输出索引->输出
				if outs, ok := UTXOs[txID]; ok {
					outs.Outputs[outIdx] = out
				} else {
					outs := TxOutputs{map[int]TxOutput{outIdx: out}}
					UTXOs[txID] = outs
				}
			}
			// spentTXOs 已经消费的输出
			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					inTXID := hex.EncodeToString(in.Txid)
					spentTXOs[inTXID] = append(spentTXOs[inTXID], in.Vout)
				}
			}
		}
	}
	return UTXOs
}

func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()
	for bci.HasNext() {
		block := bci.Next()
		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}
	}
	return Transaction{}, errors.New("transaction is not found")
}

func (bc *BlockChain) SignTransaction(tx *Transaction, priKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)
	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		PanicIfError(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	tx.Sign(priKey, prevTXs)
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	preTXs := make(map[string]Transaction)
	for _, vin := range tx.Vin {
		preTX, err := bc.FindTransaction(vin.Txid)
		PanicIfError(err)
		preTXs[hex.EncodeToString(preTX.ID)] = preTX
	}
	return tx.Verify(preTXs)
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.tip, bc.db}
}

func (bc *BlockChain) InfoMap() map[string]interface{} {
	info := make(map[string]interface{})
	var blocks []map[string]interface{}
	bci := bc.Iterator()
	for bci.HasNext() {
		block := bci.Next()
		blocks = append(blocks, block.InfoMap())
	}
	info["blocks"] = blocks
	return info
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func (bc *BlockChain) Close() {
	bc.db.Close()
}

type BlockChainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (i *BlockChainIterator) Next() *Block {
	var block *Block
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	PanicIfError(err)
	i.currentHash = block.PreBlockHash
	return block
}

func (i *BlockChainIterator) HasNext() bool {
	return len(i.currentHash) != 0
}
