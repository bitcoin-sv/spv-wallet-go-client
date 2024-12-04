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
	adminAPI, err := wallet.NewAdminAPIWithXPriv(exampleutil.ExampleConfig, examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	id := "88db6027-e38a-43b7-97a0-45f08d535256"
	contact, err := adminAPI.ContactUpdate(context.Background(), &commands.UpdateContact{
		ID:       id,
		FullName: "John Doe",
		Metadata: map[string]any{
			"phoneNumber": "123456789",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print(fmt.Sprintf("[HTTP PUT] Update contact - api/v1/admin/contacts/%s", id), contact)
}
