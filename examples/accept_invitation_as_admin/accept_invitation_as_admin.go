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

	id := "cd401f8f-a9ae-4efe-a2a1-de3e6832593d"
	err = adminAPI.AcceptInvitation(context.Background(), id)
	if err != nil {
		log.Fatal(err)
	}
}
