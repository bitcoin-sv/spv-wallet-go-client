package main

import (
	"context"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

func main() {
	usersAPI, err := wallet.NewUserAPIWithXPriv(exampleutil.ExampleConfig, examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	transaction, err := usersAPI.DraftTransaction(context.Background(), &commands.DraftTransaction{
		Config: response.TransactionConfig{
			Outputs: []*response.TransactionOutput{
				{
					To:       "receiver@example.com",
					Satoshis: 1,
				},
			},
		},
		Metadata: map[string]any{"key": "value"},
	})
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print("[HTTP POST] Draft transaction - api/v1/transactions", transaction)
}
