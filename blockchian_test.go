package main

import (
	"fmt"
	"strconv"
	"testing"
)

func TestNewBlockchain(t *testing.T) {
	blockchain := NewBlockChain()
	defer blockchain.db.Close()

	blockchain.AddBlock("transacton 1")
	blockchain.AddBlock("transacton 2")
	blockchain.AddBlock("transacton 3")

	it := blockchain.Iterator()

	for {
		block := it.Next()

		fmt.Printf("Previous hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %s\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		// When there is no more block.
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
