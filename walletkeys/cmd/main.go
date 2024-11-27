package main

import (
	"fmt"
	"log"

	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"
)

func main() {
	keys, err := walletkeys.RandomKeysWithMnemonic()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("XPriv: ", keys.Keys.XPriv())
	fmt.Println("XPub: ", keys.Keys.XPub())
	fmt.Println("Mnemonic: ", keys.Mnemonic())
}
