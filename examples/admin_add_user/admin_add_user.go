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
	adminAPI, err := wallet.NewAdminAPIWithXPub(exampleutil.NewDefaultConfig(), examples.AdminXPriv)
	if err != nil {
		log.Fatalf("Failed to initialize admin API with XPriv: %v", err)
	}

	ctx := context.Background()
	xPub, err := adminAPI.CreateXPub(ctx, &commands.CreateUserXpub{
		XPub:     examples.UserXPub,
		Metadata: queryparams.Metadata{"key": "value"},
	})
	if err != nil {
		log.Fatalf("Failed to create xPub: %v", err)
	}
	exampleutil.PrettyPrint("Created XPub", xPub)

	paymail, err := adminAPI.CreatePaymail(ctx, &commands.CreatePaymail{
		Metadata: queryparams.Metadata{"key": "value"},
		Key:      examples.UserXPub,
		Address:  examples.Paymail,
	})
	if err != nil {
		log.Fatalf("Failed to create paymail: %v", err)
	}
	exampleutil.PrettyPrint("Created paymail", paymail)
}
