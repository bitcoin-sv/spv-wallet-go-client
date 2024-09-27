/*
Package main - admin_remove_user example
*/
package main

import (
	"context"
	"os"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
)

func main() {
	defer examples.HandlePanic()

	examples.CheckIfAdminKeyExists()

	const server = "http://localhost:3003/v1"

	adminClient, err := walletclient.NewWithAdminKey(server, examples.ExampleAdminKey)
	if err != nil {
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}
	ctx := context.Background()

	err = adminClient.AdminDeletePaymail(ctx, examples.ExamplePaymail)
	if err != nil {
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}
}
