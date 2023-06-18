package main

import (
	"reflect"
	"testing"
	"time"
)

func TestBlockchain_AddBlock(t *testing.T) {
	blockchain := Blockchain{
		Chain: []Block{
			{
				Index:        0,
				Timestamp:    time.Now().String(),
				Transactions: []Transaction{},
				Hash:         "",
				PrevHash:     "",
			},
		},
		TransactionPool:     []Transaction{},
		ValidatedPool:       []Transaction{},
		TransactionID:       1,
		ConflictResolutions: 0,
	}

	transactions := []Transaction{
		{
			Sender:   "Alice",
			Receiver: "Bob",
			Message:  "Hello!",
		},
	}
	blockchain.addBlock(transactions)

	expectedChain := []Block{
		{
			Index:        0,
			Timestamp:    "", 
			Transactions: []Transaction{},
			Hash:         "",
			PrevHash:     "",
		},
		{
			Index:        1,
			Timestamp:    "", 
			Transactions: transactions,
			Hash:         "",
			PrevHash:     "",
		},
	}

	for i := range blockchain.Chain {
		blockchain.Chain[i].Timestamp = ""
		blockchain.Chain[i].Hash = ""
		blockchain.Chain[i].PrevHash = ""
	}

	if !reflect.DeepEqual(blockchain.Chain, expectedChain) {
		t.Errorf("addBlock failed, expected chain: %v, got: %v", expectedChain, blockchain.Chain)
	}
}

func TestBlockchain_Validate(t *testing.T) {
	invalidBlockchain := Blockchain{
		Chain: []Block{
			{
				Index:        0,
				Timestamp:    time.Now().String(),
				Transactions: []Transaction{},
				Hash:         "",
				PrevHash:     "",
			},
			{
				Index:     1,
				Timestamp: time.Now().String(),
				Transactions: []Transaction{
					{
						Sender:   "Alice",
						Receiver: "Bob",
						Message:  "Hello!",
					},
				},
				Hash:     "invalid_hash",
				PrevHash: "",             
			},
		},
		TransactionPool:     []Transaction{},
		ValidatedPool:       []Transaction{},
		TransactionID:       1,
		ConflictResolutions: 0,
	}

	isValid := invalidBlockchain.validate()

	if isValid {
		t.Error("validate failed, expected blockchain to be invalid")
	}
}

func TestCalculateHash(t *testing.T) {
	index := 1
	timestamp := "2023-06-12 15:58:15"
	transactions := []Transaction{
		{
			Sender:   "Alice",
			Receiver: "Bob",
			Message:  "Hello!",
		},
	}
	prevHash := "prev_hash"

	expectedHash := "313173440e3ba65b25bdc10b3f2d4240a2053d23c2bffb050835264544b142bd"

	hash := CalculateHash(index, timestamp, transactions, prevHash)

	if hash != expectedHash {
		t.Errorf("CalculateHash failed, expected: %s, got: %s", expectedHash, hash)
	}
}
