/*
Package main - create_transaction example
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

	client := walletclient.NewWithXPriv(server, examples.ExampleXPriv)
	ctx := context.Background()

	recipient := walletclient.Recipients{To: "test-multiple1@pawel.test.4chain.space", Satoshis: 1}
	recipients := []*walletclient.Recipients{&recipient}
	metadata := map[string]any{"some_metadata": "example"}

	newTransaction, err := client.SendToRecipients(ctx, recipients, metadata)
	if err != nil {
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}
	fmt.Println("SendToRecipients response: ", newTransaction)

	tx, err := client.GetTransaction(ctx, newTransaction.ID)
	if err != nil {
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}
	fmt.Println("GetTransaction response: ", tx)
}
