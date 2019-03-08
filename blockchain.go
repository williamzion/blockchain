package main

import (
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

const (
	dbFile       = "blockchain.db"
	blocksBucket = "blocks"
)

// Blockchain represents a chain of blocks.
// It’s an ordered, back-linked list. Which means that blocks are stored in the insertion order and that each block is linked to the previous one.
// For simplicity, the Blockchain only keeps a slice of blocks, no map implementation.
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// AddBlock saves provided data as a block adding to the blockchain.
func (bc *Blockchain) AddBlock(data string) {
	prev := bc.blocks[len(bc.blocks)-1]
	new := NewBlock(data, prev.Hash)
	bc.blocks = append(bc.blocks, new)
}

// NewBlockChain returns a new blockchain with genesis block.
func NewBlockChain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			fmt.Println("No existing block found. Creating a new one...")
			genesis := NewGenesisBlock()

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			// Store serialized Block structure.
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			// Store the hash of the last block in a chain.
			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	return &Blockchain{tip, db}
}
