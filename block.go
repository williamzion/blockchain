package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

// A Block is composed of headers and transactons(Data in this case). This is a
// simplified and mixed datastructure.
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

// SetHash calculates and sets block hash.
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

// NewBlock creates a block with block data and previous block hash and returns it.
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{}, // hash will be calculated block itself.
	}
	block.SetHash()

	return block
}

// NewGenesisBlock creates and returns genesis Block.
func NewGenesisBlock() *Block {
	return NewBlock("Genesis block", []byte{})
}
