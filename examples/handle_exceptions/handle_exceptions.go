/*
Package main - handle_exceptions example
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

	fmt.Println("Handle exceptions example")

	examples.CheckIfXPubExists()

	fmt.Println("XPub exists")

	const server = "http://localhost:3003/v1"

	client, err := walletclient.NewWithXPub(server, examples.ExampleAdminKey)
	if err != nil {
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}
	ctx := context.Background()

	fmt.Println("Client created")

	status, err := client.AdminGetStatus(ctx)
	if err != nil {
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}

	fmt.Println("Status: ", status)
}
