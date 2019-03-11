package main

import (
	"log"

	bolt "go.etcd.io/bbolt"
)

// BlockchainIterator represents an iterator over blockchain blocks.
// It contains currentHash cursor and a db connection and returns the next
// block from a blockchain.
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Iterator returns an iterator of a blockchain at current state starting from the tip.
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}

// Next represents a iterator cursor returns te next block in blockchain starting from top to bottom, from newest to oldest.
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash
	return block
}
