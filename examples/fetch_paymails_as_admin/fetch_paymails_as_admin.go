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
	adminAPI, err := wallet.NewAdminAPIWithXPriv(exampleutil.NewDefaultConfig(), examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	page, err := adminAPI.Paymails(context.Background(), queries.QueryWithPageFilter[filter.AdminPaymailFilter](filter.Page{Size: 1}))
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print("[HTTP GET] Paymails page - api/v1/admin/paymails", page)
}
