/*
Package main - list_transactions example
*/
package main

import (
	"context"
	"fmt"
	"os"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
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

	metadata := map[string]any{}

	conditions := filter.TransactionFilter{}
	queryParams := filter.QueryParams{}

	txs, err := client.GetTransactions(ctx, &conditions, metadata, &queryParams)
	if err != nil {
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}
	fmt.Println("GetTransactions response: ", txs)

	targetBlockHeight := uint64(839228)
	conditions = filter.TransactionFilter{BlockHeight: &targetBlockHeight}
	queryParams = filter.QueryParams{PageSize: 100, Page: 1}

	txsFiltered, err := client.GetTransactions(ctx, &conditions, metadata, &queryParams)
	if err != nil {
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}
	fmt.Println("Filtered GetTransactions response: ", txsFiltered)
}
