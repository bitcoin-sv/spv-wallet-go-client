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

	id := "bc84cb22-c9aa-415f-93e1-2ed0af27a6ef"
	err = adminAPI.RejectInvitation(context.Background(), id)
	if err != nil {
		log.Fatal(err)
	}
}
