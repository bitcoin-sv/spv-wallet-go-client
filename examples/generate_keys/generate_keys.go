package main

import (
	"fmt"
	"log"

	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"
)

func main() {
	keys, err := walletkeys.RandomKeys()
	if err != nil {
		log.Fatalf("Failed to generate random keys: %v", err)
	}

	fmt.Printf("Generated xPub for user: %s\n", keys.XPub())
	fmt.Printf("Generated xPriv for user: %s\n", keys.XPriv())
}
