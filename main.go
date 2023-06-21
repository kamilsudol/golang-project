package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Represents a transaction in the blockchain
type Transaction struct {
	Sender   string
	Receiver string
	Message  string
}

// Represents a block in the blockchain
type Block struct {
	Index        int
	Timestamp    string
	Transactions []Transaction
	Hash         string
	PrevHash     string
}

// Represents the blockchain
type Blockchain struct {
	Chain []Block
}

// Calculates the hash of a block
func CalculateHash(index int, timestamp string, transactions []Transaction, prevHash string) string {
	data := strconv.Itoa(index) + timestamp + fmt.Sprintf("%v", transactions) + prevHash
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Adds a new block to the blockchain
func (bc *Blockchain) addBlock(transactions []Transaction) {
	prevBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		Hash:         "",
		PrevHash:     prevBlock.Hash,
	}
	newBlock.Hash = CalculateHash(newBlock.Index, newBlock.Timestamp, newBlock.Transactions, newBlock.PrevHash)
	bc.Chain = append(bc.Chain, newBlock)
}

// Handles the creation of a new transaction and adds it to the blockchain
func handleMessage(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside handle message...", r)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sender := r.FormValue("sender")
	receiver := r.FormValue("receiver")
	message := r.FormValue("message")

	if sender == "" || receiver == "" || message == "" {
		renderMessengerPage(w, blockchain, "All fields are required", sender, receiver)
		return
	}

	transaction := Transaction{
		Sender:   sender,
		Receiver: receiver,
		Message:  message,
	}
	log.Println("add new block, transaction: ", transaction)
	log.Println("blockchain: ", blockchain)
	blockchain.addBlock([]Transaction{transaction})
	log.Println("end request, transaction pool: ", transaction)
	renderMessengerPage(w, blockchain, "", "", "")
}

// Handles the validation of the blockchain
func handleValidation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	isValid := blockchain.validate()

	if isValid {
		message := "Blockchain is valid."
		renderMessengerPage(w, blockchain, message, "", "")
	} else {
		message := "Blockchain is not valid."
		renderMessengerPage(w, blockchain, message, "", "")
	}
}

// Handles the conflict resolution of the blockchain
func HandleConflictResolution(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	blockchain.resolveConflicts()

	message := "Conflict resolution completed."
	renderMessengerPage(w, blockchain, message, "", "")
}

// Checks the validity of the blockchain
func (bc *Blockchain) validate() bool {
	for i := 1; i < len(bc.Chain); i++ {
		currentBlock := bc.Chain[i]
		prevBlock := bc.Chain[i-1]

		// Verify the block hash
		if currentBlock.Hash != CalculateHash(currentBlock.Index, currentBlock.Timestamp, currentBlock.Transactions, currentBlock.PrevHash) {
			return false
		}

		// Verify the previous block hash
		if currentBlock.PrevHash != prevBlock.Hash {
			return false
		}
	}

	return true
}

// Resolves conflicts in the blockchain by choosing the longest chain
func (bc *Blockchain) resolveConflicts() {
	maxLength := len(bc.Chain)
	longestChain := bc.Chain

	for _, block := range bc.Chain {
		if len(block.Transactions) > maxLength {
			maxLength = len(block.Transactions)
			longestChain = append([]Block{}, block)
		} else if len(block.Transactions) == maxLength {
			longestChain = append(longestChain, block)
		}
	}

	bc.Chain = longestChain
}

// Renders the messenger page with the provided data
func renderMessengerPage(w http.ResponseWriter, bc Blockchain, message string, sender string, receiver string) {
	tmpl := template.Must(template.ParseFiles("messenger.html"))

	blockData := make([]BlockData, len(bc.Chain))
	for i, block := range bc.Chain {
		blockData[i] = BlockData{
			Index:        block.Index,
			Timestamp:    block.Timestamp,
			Transactions: block.Transactions,
			Hash:         block.Hash,
			PrevHash:     block.PrevHash,
		}
	}

	data := MessengerPage{
		Blockchain: bc,
		Message:    message,
		Sender:     sender,
		Receiver:   receiver,
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Function for the messenger page
func messengerHandler(w http.ResponseWriter, r *http.Request) {
	renderMessengerPage(w, blockchain, "", "", "")
}

type MessengerPage struct {
	Blockchain Blockchain
	Message    string
	Sender     string
	Receiver   string
}

type BlockData struct {
	Index        int
	Timestamp    string
	Transactions []Transaction
	Hash         string
	PrevHash     string
}

var blockchain Blockchain

func main() {
	// Create the initial block in the blockchain
	initializeBlockChain()

	http.HandleFunc("/", messengerHandler)
	http.HandleFunc("/message", handleMessage)
	http.HandleFunc("/validate", handleValidation)
	http.HandleFunc("/resolve", HandleConflictResolution)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initializeBlockChain() {
	genesisBlock := Block{
		Index:        0,
		Timestamp:    time.Now().String(),
		Transactions: []Transaction{},
		Hash:         "",
		PrevHash:     "",
	}
	genesisBlock.Hash = CalculateHash(genesisBlock.Index, genesisBlock.Timestamp, genesisBlock.Transactions, genesisBlock.PrevHash)

	blockchain = Blockchain{
		Chain: []Block{genesisBlock},
	}
}
