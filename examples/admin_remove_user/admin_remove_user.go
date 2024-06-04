package main

import (
	"context"
	"examples"
	"fmt"
	"os"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
)

func main() {
	defer examples.HandlePanic()

	const server = "http://localhost:3003/v1"

	if examples.ExampleAdminKey == "" {
		fmt.Println(examples.ErrMessage("adminKey"))
		os.Exit(1)
	}

	adminClient := walletclient.NewWithAdminKey(server, examples.ExampleAdminKey)
	ctx := context.Background()

	adminClient.AdminDeletePaymail(ctx, examples.ExamplePaymail)
}
