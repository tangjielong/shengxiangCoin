package blockchain

import (
	"testing"
)

const foo = "1Pg9f5myjVpxHFdg7R7egzyExNPzGEdscS"
const bar = "14dVjU2fw8ZdZ3p2sszvZFW18PTH7Zw13p"

func TestCreateWallet(t *testing.T) {
	address := CreateWallet()
	t.Log(address)
}

func TestListAddresses(t *testing.T) {
	t.Log(ListAddresses())
}

func TestSend(t *testing.T) {
	Send(bar, foo, 1)
	t.Log("Success")
}

func TestGetBalance(t *testing.T) {
	balance := GetBalance(foo)
	t.Log(balance)
}
