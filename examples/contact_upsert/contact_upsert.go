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

	contactPaymail := "john.doe@example.com"
	contact, err := usersAPI.UpsertContact(context.Background(), commands.UpsertContact{
		ContactPaymail: contactPaymail,
		FullName:       "John Doe",
		Metadata: map[string]any{
			"key": "value",
		},
		// Optional field representing the paymail address of the user who is creating the contact.
		// It is required in case if user has multiple paymail addresses associated with single xPub.
		//RequesterPaymail: "",
	})
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print(fmt.Sprintf("[HTTP PUT] Upsert contact - api/v1/contacts/%s", contactPaymail), contact)
}
