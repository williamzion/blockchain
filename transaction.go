package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

// subsidy is the amount of reward.
const subsidy = 10

// Transaction represents a Bitcoin transaction.
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// IsCoinbase checks whether the transaction is coinbase.
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// SetID sets ID of a transaction.
func (tx *Transaction) SetID() {
	var (
		encoded bytes.Buffer
		hash    [32]byte
	)

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
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

// NewUTXOTransaction creates a new transaction.
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var (
		inputs  []TXInput
		outputs []TXOutput
	)

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("error: not enough funds")
	}

	// Build a list of inputs.
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{
				Txid:      txID,
				Vout:      out,
				ScriptSig: from,
			}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs.
	// outputs that’s locked with the receiver address. This is the actual transferring of coins to other address.
	outputs = append(outputs, TXOutput{
		Value:        amount,
		ScriptPubKey: to,
	})
	if acc > amount {
		// outputs that’s locked with the sender address. This is a change.
		outputs = append(outputs, TXOutput{
			Value:        acc - amount, // a change
			ScriptPubKey: from,
		})
	}

	tx := Transaction{
		ID:   nil,
		Vin:  inputs,
		Vout: outputs,
	}
	tx.SetID()
	return &tx
}
