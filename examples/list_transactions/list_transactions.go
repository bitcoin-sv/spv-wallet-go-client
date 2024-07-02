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

	client := walletclient.NewWithXPriv(server, examples.ExampleXPriv)
	ctx := context.Background()

	metadata := map[string]any{
		"note": "user-id-123",
	}

	conditions := filter.TransactionFilter{}
	queryParams := filter.QueryParams{}

	txs, err := client.GetTransactions(ctx, &conditions, metadata, &queryParams)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("GetTransactions response: ", txs)

	conditions = filter.TransactionFilter{BlockHeight: func(i uint64) *uint64 { return &i }(839228)}
	queryParams = filter.QueryParams{PageSize: 100, Page: 1}

	txsFiltered, err := client.GetTransactions(ctx, &conditions, metadata, &queryParams)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Filtered GetTransactions response: ", txsFiltered)
}
