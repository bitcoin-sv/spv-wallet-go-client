/*
Package main - admin_add_user example
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

	examples.CheckIfAdminKeyExists()

	server := "http://localhost:3003/v1"

	adminClient, err := walletclient.NewWithAdminKey(server, examples.ExampleAdminKey)
	if err != nil {
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}
	ctx := context.Background()

	metadata := map[string]any{"some_metadata": "example"}

	err = adminClient.AdminNewXpub(ctx, examples.ExampleXPub, metadata)
	if err != nil {
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}

	createPaymailRes, err := adminClient.AdminCreatePaymail(ctx, examples.ExampleXPub, examples.ExamplePaymail, "Some public name", "")
	if err != nil {
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}
	fmt.Println("AdminCreatePaymail response: ", createPaymailRes)
}
