package main

import "crypto/sha256"

// MerkleTree represent a Merkle tree.
type MerkleTree struct {
	RootNode *MerkleNode
}

// MerkleNode represent a Merkle tree node.
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}


