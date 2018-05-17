package blockchain

import (
	"log"
)

func PanicIfError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
