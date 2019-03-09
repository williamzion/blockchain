package main

import (
	"log"

	bolt "go.etcd.io/bbolt"
)

const (
	dbFile              = "blockchain.db"
	blocksBucket        = "blocks"
	genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

// Blockchain represents a chain of blocks.
// It’s an ordered, back-linked list. Which means that blocks are stored in the insertion order and that each block is linked to the previous one.
// For simplicity, the Blockchain only keeps a slice of blocks, no map implementation.
type Blockchain struct {
	// Hash of last block in blockchain. It is an identifier of a blockchain.
	tip []byte
	// persistent db connection.
	db *bolt.DB
}

// BlockchainIterator represents an iterator over blockchain blocks.
// It contains currentHash cursor and a db connection and returns the next
// block from a blockchain.
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// AddBlock saves provided data as a block adding to the blockchain.
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		// Update bockchain tip once a new block was added.
		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash

		return nil
	})
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

// NewBlockChain returns a new blockchain with genesis block.
// It takes an address which will receive the reward for mining the genesis block.
// A db connection included in the returned value is intended to be reused.
func NewBlockChain(address string) *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)

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

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db}
}
