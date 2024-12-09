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
	usersAPI, err := wallet.NewUserAPIWithXPriv(exampleutil.NewDefaultConfig(), examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	accessKey, err := usersAPI.GenerateAccessKey(ctx, &commands.GenerateAccessKey{
		Metadata: map[string]any{"key": "value"},
	})
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print("[HTTP POST] Generate access key - api/v1/users/current/keys", accessKey)
}
