package main

import (
	"fmt"
	"log"

	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"
)

func main() {
	// This is an example mnemonic phrase - replace it with your own
	const mnemonicPhrase = "nut same spike popular already mercy kit board rent light illegal local eight filter tube"

	key, err := walletkeys.XPrivFromMnemonic(mnemonicPhrase)
	if err != nil {
		log.Fatalf("Failed to get xPriv from mnemonic: %v", err)
	}

	fmt.Printf("Extracted xPriv: %s\n", key.String())
}
