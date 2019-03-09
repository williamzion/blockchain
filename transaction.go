package main

import "fmt"

// subsidy is the amount of reward.
const subsidy = 10

// Transaction represents a Bitcoin transaction.
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// TXOutput represents a transaction output.
type TXOutput struct {
	Value        int
	ScriptPubKey string
}

// TXInput represents a transaction input.
type TXInput struct {
	// references previous output.
	Txid []byte
	// index of an output in the transaction.
	Vout int
	// script providing data to be used in an output’s ScriptPubKey
	ScriptSig string
}

// NewCoinbaseTX creates a new coinbase transaction.
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to %q", to)
	}

	txin := TXInput{
		Txid:      []byte{},
		Vout:      -1,
		ScriptSig: data,
	}
	txout := TXOutput{
		Value:        subsidy,
		ScriptPubKey: to,
	}
	tx := Transaction{
		ID:   nil,
		Vin:  []TXInput{txin},
		Vout: []TXOutput{txout},
	}

	tx.SetID()
	return &tx
}
