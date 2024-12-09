package main

import (
	"context"
	"fmt"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
)

func main() {
	usersAPI, err := wallet.NewUserAPIWithXPriv(exampleutil.NewDefaultConfig(), examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	transactionID := "86cafa5b-fdaa-4629-ae46-78d68d6a180b"
	transaction, err := usersAPI.UpdateTransactionMetadata(context.Background(), &commands.UpdateTransactionMetadata{
		ID:       transactionID,
		Metadata: map[string]any{"new_key": "new_value"},
	})
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print(fmt.Sprintf("[HTTP PATCH] Update transaction metadata - api/v1/transactions/%s", transactionID), transaction)
}
