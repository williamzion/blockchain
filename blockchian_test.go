package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestNewBlockchain(t *testing.T) {
	blockchain := NewBlockChain()

	blockchain.AddBlock("transacton 1")
	blockchain.AddBlock("transacton 2")
	blockchain.AddBlock("transacton 3")

	for i, block := range blockchain.blocks {
		fmt.Printf("Prev Hash: %x\n", block.PrevBlockHash)
		fmt.Println("Data:", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()

		if i > 0 {
			if !bytes.Equal(block.PrevBlockHash, blockchain.blocks[i-1].Hash) {
				t.Fatalf("prev hash (%x) != current hash (%x), want equal", blockchain.blocks[i-1].Hash, block.Hash)
			}
		}
	}
}
