package main

import (
	"context"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

func main() {
	usersAPI, err := wallet.NewUserAPIWithXPriv(exampleutil.NewDefaultConfig(), examples.UserXPriv)
	if err != nil {
		log.Fatalf("Failed to initialize user API with XPriv: %v", err)
	}

	page, err := usersAPI.Transactions(context.Background(), queries.QueryWithPageFilter[filter.TransactionFilter](filter.Page{
		Size: 1,
		Sort: "asc",
	}))
	if err != nil {
		log.Fatalf("Failed to fetch transactions: %v", err)
	}
	exampleutil.PrettyPrint("Fetched transactions", page.Content)
}
