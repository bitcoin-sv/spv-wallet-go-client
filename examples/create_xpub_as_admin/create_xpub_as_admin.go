package main

import (
	"context"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/queryparams"
)

func main() {
	adminAPI, err := wallet.NewAdminAPIWithXPriv(exampleutil.NewDefaultConfig(), examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	xPub, err := adminAPI.CreateXPub(context.Background(), &commands.CreateUserXpub{
		XPub:     "1c318ad8-5ee4-42d3-9cf5-5b0babec9156",
		Metadata: queryparams.Metadata{"key": "value"},
	})
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print("[HTTP POST] Create XPub - api/v1/admin/users", xPub)
}
