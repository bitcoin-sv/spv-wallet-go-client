package main

import (
	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
)

func main() {
	// Replace with created access key
	exampleAccessKey := "some_generated_access_key"

	// Create a client
	_, _ = walletclient.NewWithAccessKey(exampleAccessKey, "http://localhost:3003/v1")

}
