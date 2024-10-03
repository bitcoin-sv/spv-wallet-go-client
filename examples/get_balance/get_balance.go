/*
Package main - get_balance example
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

	const server = "http://localhost:3003/v1"

	client, err := walletclient.NewWithXPriv(server, examples.ExampleXPriv)
	if err != nil {
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}
	ctx := context.Background()

	xpubInfo, err := client.GetXPub(ctx)
	if err != nil {
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}
	fmt.Println("Current balance: ", xpubInfo.CurrentBalance)
}
