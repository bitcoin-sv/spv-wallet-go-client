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

	examples.CheckIfXPubExists()

	const server = "http://localhost:3003/v1"

	client := walletclient.NewWithXPub(server, examples.ExampleXPub)
	ctx := context.Background()

	status, err := client.AdminGetStatus(ctx)
	if err != nil {
		fmt.Println("Response status: ", err.GetStatusCode())
		fmt.Println("Content: ", err.Error())

		os.Exit(1)
	}

	fmt.Println("Status: ", status)
}
