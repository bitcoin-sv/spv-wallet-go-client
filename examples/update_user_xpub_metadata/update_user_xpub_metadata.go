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
	usersAPI, err := wallet.NewUserAPIWithXPriv(exampleutil.NewDefaultConfig(), examples.UserXPriv)
	if err != nil {
		log.Fatalf("Failed to initialize user API with XPriv: %v", err)
	}

	ctx := context.Background()
	xPub, err := usersAPI.XPub(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch xPub: %v", err)
	}
	exampleutil.PrettyPrint("User xPub info before update", xPub)

	xPub, err = usersAPI.UpdateXPubMetadata(ctx, &commands.UpdateXPubMetadata{
		Metadata: queryparams.Metadata{"new_key": "new_value"},
	})
	if err != nil {
		log.Fatalf("Failed to fetch xPub: %v", err)
	}
	exampleutil.PrettyPrint("User xPub info after update", xPub)
}
