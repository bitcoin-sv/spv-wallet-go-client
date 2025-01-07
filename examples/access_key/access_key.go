package main

import (
	"context"
	"fmt"
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
	generated, err := usersAPI.GenerateAccessKey(ctx, &commands.GenerateAccessKey{
		Metadata: queryparams.Metadata{"key": "value"},
	})
	if err != nil {
		log.Fatalf("Failed to generate access key: %v", err)
	}
	exampleutil.PrettyPrint("Generated access key", generated)

	fetched, err := usersAPI.AccessKey(ctx, generated.ID)
	if err != nil {
		log.Fatalf("Failed to fetch access key: %v", err)
	}
	exampleutil.PrettyPrint("Fetched access key", fetched)

	err = usersAPI.RevokeAccessKey(ctx, generated.ID)
	if err != nil {
		log.Fatalf("Failed to revoke access key: %v", err)
	}
	fmt.Printf("Revoke access key: %s\n", generated.ID)
}
