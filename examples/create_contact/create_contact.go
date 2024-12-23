package main

import (
	"context"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
)

func main() {
	ctx := context.Background()
	adminAPI, err := wallet.NewAdminAPIWithXPriv(exampleutil.NewDefaultConfig(), examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	paymail := "john.doe@example"
	contact, err := adminAPI.CreateContact(ctx, &commands.CreateContact{
		Paymail:        paymail,
		CreatorPaymail: "admin@example",
		FullName:       "John Doe",
	})
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Printf("Create Paymail - api/v1/admin/contacts/%s", contact, "", 0, paymail)
}
