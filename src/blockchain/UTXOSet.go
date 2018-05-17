package blockchain

import (
	"encoding/hex"
	"github.com/boltdb/bolt"
)

const utxoBucket = "utxo"

type UTXOSet struct {
	BlockChain *BlockChain
}

func (set UTXOSet) ReIndex() {
	db := set.BlockChain.db
	bucketName := []byte(utxoBucket)
	// 创建utxo数据库
	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			return err
		}
		_, err = tx.CreateBucket(bucketName)
		return err
	})
	PanicIfError(err)
	// 遍历blockchain，找到那些没有消费的交易的输出（余额）
	UTXO := set.BlockChain.FindUTXO()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				return err
			}
			// 保存 txid->UTXO
			err = b.Put(key, outs.Serialize())
			return err
		}
		return nil
	})
	PanicIfError(err)
}

func (set UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := set.BlockChain.db
	bucketName := []byte(utxoBucket)

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		c := b.Cursor()
		// 遍历UTXO集
	Loop:
		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := DeserializeOutputs(v)
			for outIdx, out := range outs.Outputs {
				// 找出相同地址的输出
				if out.IsLockedWithKey(pubKeyHash) {
					accumulated += out.Value
					if ids, ok := unspentOutputs[txID]; ok {
						unspentOutputs[txID] = append(ids, outIdx)
					} else {
						unspentOutputs[txID] = []int{outIdx}
					}
					// 足够
					if accumulated >= amount {
						break Loop
					}
				}
			}
		}
		return nil
	})
	PanicIfError(err)
	return accumulated, unspentOutputs
}

func (set UTXOSet) FindUTXO(pubKeyHash []byte) []TxOutput {
	var UTXOs []TxOutput
	db := set.BlockChain.db
	bucketName := []byte(utxoBucket)
	// 遍历UTXO集
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := DeserializeOutputs(v)
			for _, out := range outs.Outputs {
				// 找出相同地址的输出
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}
		return nil
	})
	PanicIfError(err)
	return UTXOs
}

func (set UTXOSet) Update(block *Block) {
	db := set.BlockChain.db
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		for _, tx := range block.Transactions {
			if !tx.IsCoinbase() {
				for _, vin := range tx.Vin {
					outBytes := b.Get(vin.Txid)
					outs := DeserializeOutputs(outBytes)
					delete(outs.Outputs, vin.Vout)
					if len(outs.Outputs) == 0 {
						// 这个交易已经没有未消费的输出
						err := b.Delete(vin.Txid)
						if err != nil {
							return err
						}
					} else {
						// 更新
						err := b.Put(vin.Txid, outs.Serialize())
						if err != nil {
							return err
						}
					}
				}
			}
			// 当前交易的所有输出都是未消费输出
			newOutputs := NewTxOutputs()
			for outIdx, out := range tx.Vout {
				newOutputs.Outputs[outIdx] = out
			}
			err := b.Put(tx.ID, newOutputs.Serialize())
			return err
		}
		return nil
	})
	PanicIfError(err)
}
