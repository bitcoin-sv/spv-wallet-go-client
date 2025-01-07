package main

import (
	"fmt"
	"log"

	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"
)

func main() {
	// This is an example xPriv key - replace it with your own
	const xPriv = "xprv9s21ZrQH143K4VneY3UWCF1o5Kk2tmgGrGtMtsrThCTsHsszEZ6H1iP37ZTwuUBvMwudG68SRkcfTjeu8h3rkayfyqkjKAStFBkuNsBnAkS"

	xPub, err := walletkeys.XPubFromXPriv(xPriv)
	if err != nil {
		log.Fatalf("Failed to get xPriv from mnemonic: %v", err)
	}
	fmt.Printf("Extracted xPub: %s\n", xPub)
}
