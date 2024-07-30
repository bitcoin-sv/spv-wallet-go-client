/*
Package main - get_shared_config example
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

	sharedConfig, err := client.GetSharedConfig(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Shared config (PaymailDomains): ", sharedConfig.PaymailDomains)
	fmt.Println("Shared config (ExperimentalFeatures): ", sharedConfig.ExperimentalFeatures)
}
