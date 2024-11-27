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
	usersAPI, err := wallet.NewUserAPIWithXPriv(exampleutil.ExampleConfig, examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	xPub, err := usersAPI.UpdateXPubMetadata(context.Background(), &commands.UpdateXPubMetadata{
		Metadata: map[string]any{"key": "value"},
	})
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print("[HTTP PATCH] Current user xPub metadata update - api/v1/users/current", xPub)
}
