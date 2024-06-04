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

	server := "http://localhost:3003/v1"

	if examples.ExampleAdminKey == "" {
		fmt.Println(examples.ErrMessage("adminKey"))
		os.Exit(1)
	}

	adminClient := walletclient.NewWithAdminKey(server, examples.ExampleAdminKey)
	ctx := context.Background()

	metadata := map[string]any{"some_metadata": "example"}

	newXPubRes := adminClient.AdminNewXpub(ctx, examples.ExampleXPub, metadata)
	fmt.Println("AdminNewXpub response: ", newXPubRes)

	createPaymailRes, err := adminClient.AdminCreatePaymail(ctx, examples.ExampleXPub, examples.ExamplePaymail, "Some public name", "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("AdminCreatePaymail response: ", createPaymailRes)
}
