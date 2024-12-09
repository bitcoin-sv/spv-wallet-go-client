package main

import (
	"context"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
)

func main() {
	adminAPI, err := wallet.NewAdminAPIWithXPriv(exampleutil.NewDefaultConfig(), examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	id := "88db6027-e38a-43b7-97a0-45f08d535256"
	err = adminAPI.DeleteContact(context.Background(), id)
	if err != nil {
		log.Fatal(err)
	}
}
