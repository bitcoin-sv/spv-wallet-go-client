/*
Package main - create_transaction example
*/
package main

import (
	"context"
	"fmt"
	"os"

	"examples"
	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
)

func main() {
	defer examples.HandlePanic()

	examples.CheckIfXPrivExists()

	const server = "http://localhost:3003/v1"

	client := walletclient.NewWithXPriv(server, examples.ExampleXPriv)
	ctx := context.Background()

	recipient := walletclient.Recipients{To: "receiver@example.com", Satoshis: 1}
	recipients := []*walletclient.Recipients{&recipient}
	metadata := map[string]any{"some_metadata": "example"}

	newTransaction, err := client.SendToRecipients(ctx, recipients, metadata)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("SendToRecipients response: ", newTransaction)

	tx, err := client.GetTransaction(ctx, newTransaction.ID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("GetTransaction response: ", tx)

}
