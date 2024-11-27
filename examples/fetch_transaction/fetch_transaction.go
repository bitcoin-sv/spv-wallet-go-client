package main

import (
	"context"
	"fmt"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
)

func main() {
	usersAPI, err := wallet.NewUserAPIWithXPriv(exampleutil.ExampleConfig, examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	transactionID := "d291ac4c-04b3-48cc-a7e3-1338ac42810b"
	transaction, err := usersAPI.Transaction(context.Background(), transactionID)
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print(fmt.Sprintf("[HTTP GET] Transaction - api/v1/transactions/%s", transactionID), transaction)
}
