/*
Package main - send_op_return example
*/
package main

import (
	"context"
	"fmt"
	"os"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet/models"
)

func main() {
	defer examples.HandlePanic()

	examples.CheckIfXPrivExists()

	const server = "http://localhost:3003/v1"

	client := walletclient.NewWithXPriv(server, examples.ExampleXPriv)
	ctx := context.Background()

	metadata := map[string]any{}

	opReturn := models.OpReturn{StringParts: []string{"hello", "world"}}
	transactionConfig := models.TransactionConfig{Outputs: []*models.TransactionOutput{{OpReturn: &opReturn}}}

	draftTransaction, err := client.DraftTransaction(ctx, &transactionConfig, metadata)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("DraftTransaction response: ", draftTransaction)

	finalized, err := client.FinalizeTransaction(draftTransaction)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	transaction, err := client.RecordTransaction(ctx, finalized, draftTransaction.ID, metadata)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Transaction with OP_RETURN: ", transaction)
}
