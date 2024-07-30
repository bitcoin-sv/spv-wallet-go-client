/*
Package main - access_key example
*/
package main

import (
	"context"
	"fmt"
	"os"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
)

func main() {
	defer examples.HandlePanic()

	examples.CheckIfXPrivExists()

	const server = "http://localhost:3003/api/v1"

	client := walletclient.NewWithXPriv(server, examples.ExampleXPriv)
	ctx := context.Background()

	metadata := map[string]any{"some_metadata": "example"}
	createdAccessKey, err := client.CreateAccessKey(ctx, metadata)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Created access key ID: ", createdAccessKey.ID)
	fmt.Println("Metadata: ", createdAccessKey.Metadata)
	fmt.Println("Created at: ", createdAccessKey.CreatedAt)

	fetchedAccessKey, err := client.GetAccessKey(ctx, createdAccessKey.ID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Fetched access key ID: ", fetchedAccessKey.ID)

	revokedAccessKey, err := client.RevokeAccessKey(ctx, createdAccessKey.ID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Revoked access key ID: ", revokedAccessKey.ID)
	fmt.Println("Revoked at: ", revokedAccessKey.RevokedAt)
}
