package blockchain

import (
	"testing"
)

const lin = "1Bknnztugm7d8s813CFvUcJHhzezUUDPYM"

func TestCreateBlockChain(t *testing.T) {
	CreateBlockChain(lin)
	t.Log("Done")
}

func TestPrintChain(t *testing.T) {
	t.Log(PrintChain())
}

func TestReindexUTXO(t *testing.T) {
	ReindexUTXO()
}
