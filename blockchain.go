package main

// Blockchain represents a chain of blocks.
// It’s an ordered, back-linked list. Which means that blocks are stored in the insertion order and that each block is linked to the previous one.
// For simplicity, the Blockchain only keeps a slice of blocks, no map implementation.
type Blockchain struct {
	blocks []*Block
}

// AddBlock saves provided data as a block adding to the blockchain.
func (bc *Blockchain) AddBlock(data string) {
	prev := bc.blocks[len(bc.blocks)-1]
	new := NewBlock(data, prev.Hash)
	bc.blocks = append(bc.blocks, new)
}

// NewBlockChain returns a new blockchain with genesis block.
func NewBlockChain() *Blockchain {
	return &Blockchain{
		blocks: []*Block{NewGenesisBlock()},
	}
}
