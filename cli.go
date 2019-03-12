package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

// CLI represents command line.
type CLI struct{}

func (cli *CLI) createBlockChain(address string) {
	if !ValidateAddr(address) {
		log.Panic("error: address is not valid")
	}
	bc := CreateBlockChain(address)
	defer bc.db.Close()

	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex()
	
	fmt.Println("Done creating blockchain.")
}

func (cli *CLI) getBalance(address string) {
	if !ValidateAddr(address) {
		log.Panic("error: address is not valid")
	}
	bc := NewBlockChain()
	defer bc.db.Close()

	// The account balance is the sum of values of all unspent transaction outputs locked by the account address.
	balance := 0
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := bc.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %q: %d\n", address, balance)
}

func (cli *CLI) createWallet() {
	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Printf("Your new address: %s\n", address)
}

func (cli *CLI) listAllAddrs() {
	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddrs()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" createblockchain -address ADDRESS: Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println(" createwallet: Generate a new key pair and saves it to the wallet file")
	fmt.Println(" listaddresses: List all addresses from the wallet file")
	fmt.Println(" getbalance -address ADDRESS: Get balance of ADDRESS")
	fmt.Println("	printchain: Print all blocks of the blockchain")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT: Send AMOUNT of coins from FROM address to TO")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printChain() {
	bc := NewBlockChain()
	defer bc.db.Close()

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Previous hash: %x\n", block.PrevBlockHash)
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

func (cli *CLI) send(from, to string, amount int) {
	if !ValidateAddr(from) {
		log.Panic("error: address is not valid")
	}
	if !ValidateAddr(to) {
		log.Panic("error: address is not valid")
	}
	bc := NewBlockChain()
	UTXOSet := UTXOSet{bc}
	defer bc.db.Close()

	tx := NewUTXOTransaction(from, to, amount, &UTXOSet)
	cbTx := NewCoinbaseTX(from, "")
	txs := []*Transaction{cbTx, tx}

	bc.MineBlock(txs)
	fmt.Println("Success!")
}

// Run is an entry point for CLI, it parses command line arguments and process es commands.
func (cli *CLI) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddrsCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddr := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddr := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Receiver wallet address.")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddrsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddr == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddr)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddr == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockChain(*createBlockchainAddr)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if listAddrsCmd.Parsed() {
		cli.listAllAddrs()
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}
