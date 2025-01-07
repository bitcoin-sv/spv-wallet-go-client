package main

import (
	"context"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/queryparams"
)

func main() {
	usersAPI, err := wallet.NewUserAPIWithXPriv(exampleutil.NewDefaultConfig(), examples.UserXPriv)
	if err != nil {
		log.Fatalf("Failed to initialize user API with XPriv: %v", err)
	}

	ctx := context.Background()
	created, err := usersAPI.SendToRecipients(ctx, &commands.SendToRecipients{
		Recipients: []*commands.Recipients{
			{
				Satoshis: 1,
				To:       "alice@example.com",
			},
		},
		Metadata: queryparams.Metadata{"key": "value"},
	})
	if err != nil {
		log.Fatalf("Failed to create transaction: %v", err)
	}
	exampleutil.PrettyPrint("Created transaction", created)

	fetch, err := usersAPI.Transaction(ctx, created.ID)
	if err != nil {
		log.Fatalf("Failed to fetch transaction: %v", err)
	}

	exampleutil.PrettyPrint("Fetched transaction", fetch)
}
