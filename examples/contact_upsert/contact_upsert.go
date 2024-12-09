package main

import (
	"context"
	"fmt"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
)

func main() {
	usersAPI, err := wallet.NewUserAPIWithXPriv(exampleutil.NewDefaultConfig(), examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	paymail := "john.doe@example.com"
	contact, err := usersAPI.UpsertContact(context.Background(), commands.UpsertContact{
		FullName: "John Doe",
		Metadata: map[string]any{
			"key": "value",
		},
		Paymail: paymail,
	})
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print(fmt.Sprintf("[HTTP PUT] Upsert contact - api/v1/contacts/%s", paymail), contact)
}
