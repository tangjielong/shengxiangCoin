package main

import (
	bc "blockchain"
	"flag"
	"fmt"
	"os"
)

type CLI struct{}

func main(){
	cli := &CLI{}
	cli.Run()
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  reindexutxo - Rebuilds the UTXO set")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		bc.PanicIfError(getBalanceCmd.Parse(os.Args[2:]))
	case "createblockchain":
		bc.PanicIfError(createBlockchainCmd.Parse(os.Args[2:]))
	case "createwallet":
		bc.PanicIfError(createWalletCmd.Parse(os.Args[2:]))
	case "listaddresses":
		bc.PanicIfError(listAddressesCmd.Parse(os.Args[2:]))
	case "printchain":
		bc.PanicIfError(printChainCmd.Parse(os.Args[2:]))
	case "send":
		bc.PanicIfError(sendCmd.Parse(os.Args[2:]))
	case "reindexutxo":
		bc.PanicIfError(reindexUTXOCmd.Parse(os.Args[2:]))
	default:
		cli.printUsage()
		os.Exit(1)
	}
	// 创建区块链
	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		bc.CreateBlockchain(*createBlockchainAddress)
		fmt.Println("成功创建区块链")
	}
	// 打印区块链
	if printChainCmd.Parsed() {
		fmt.Println(bc.PrintChain())
	}
	// 重新索引未消费交易的输出集
	if reindexUTXOCmd.Parsed() {
		bc.ReindexUTXO()
		fmt.Println("成功索引UTXO集")
	}
	// 创建钱包
	if createWalletCmd.Parsed() {
		address := bc.CreateWallet()
		fmt.Printf("新的钱包地址:%s", address)
	}
	// 所有的钱包地址
	if listAddressesCmd.Parsed() {
		fmt.Println(bc.ListAddresses())
	}
	// 查询余额
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		fmt.Printf("您的余额:$%d", bc.GetBalance(*getBalanceAddress))
	}
	// 转账
	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		bc.Send(*sendFrom, *sendTo, *sendAmount)
		fmt.Println("转账完成")
	}
}
