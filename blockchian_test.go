package main

import (
	"fmt"
	"strconv"
	"testing"
)

func TestNewBlockchain(t *testing.T) {
	blockchain := NewBlockChain()

	blockchain.AddBlock("transacton 1")
	blockchain.AddBlock("transacton 2")
	blockchain.AddBlock("transacton 3")

	for _, block := range blockchain.blocks {
		fmt.Printf("Prev Hash: %x\n", block.PrevBlockHash)
		fmt.Println("Data:", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()

		pow := NewProofOfWork(block)
		fmt.Printf("pow:%s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
