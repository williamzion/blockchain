package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

// A Block is composed of headers and transactons(Data in this case). This is a
// simplified and mixed datastructure.
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// NewBlock creates a block with block data and previous block hash and returns it.
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{}, // hash will be calculated block itself.
		Nonce:         0,
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock creates and returns genesis Block.
func NewGenesisBlock() *Block {
	return NewBlock("Genesis block", []byte{})
}

// Serialize encodes a block struct into gob data.
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// DeserializeBlock decodes a block in gob encoding into block struct.
func DeserializeBlock(d []byte) *Block {
	var b Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&b)
	if err != nil {
		log.Panic(err)
	}
	return &b
}
