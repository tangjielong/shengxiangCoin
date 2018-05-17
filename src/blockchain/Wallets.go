package blockchain

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

type Wallets struct {
	Wallets map[string]*Wallet
}

func GetWallets() (*Wallets, error) {
	wallets := Wallets{make(map[string]*Wallet)}
	err := wallets.LoadFromFile()
	return &wallets, err
}

func (ws Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := wallet.GetAddress()
	ws.Wallets[address] = wallet
	return address
}

func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws Wallets) GetAddresses() []string {
	var addresses []string
	for k := range ws.Wallets {
		addresses = append(addresses, k)
	}
	return addresses
}

func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	PanicIfError(err)
	ws.Wallets = wallets.Wallets
	return nil
}

func (ws Wallets) SaveToFile() {
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	PanicIfError(encoder.Encode(ws))
	PanicIfError(ioutil.WriteFile(walletFile, content.Bytes(), 0644))
}
