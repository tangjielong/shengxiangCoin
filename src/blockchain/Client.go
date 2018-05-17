package blockchain

type Client struct {
}

func GetBalance(address string) int {
	bc := GetBlockChain()
	defer bc.Close()
	UTXOSet := UTXOSet{bc}
	balance := 0
	UTXOs := UTXOSet.FindUTXO(GetPubKeyHash(address))
	for _, out := range UTXOs {
		balance += out.Value
	}
	return balance
}

func Send(from, to string, amount int) {
	bc := GetBlockChain()
	defer bc.Close()
	UTXOSet := UTXOSet{bc}
	tx := NewUTXOTransaction(from, to, amount, bc)
	cbTx := NewCoinbaseTx(from)
	newBlock := bc.MineBlock([]*Transaction{tx, cbTx})
	UTXOSet.Update(newBlock)
}

func CreateWallet() string {
	wallets, _ := GetWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()
	return address
}

func ListAddresses() []string {
	wallets, _ := GetWallets()
	return wallets.GetAddresses()
}
