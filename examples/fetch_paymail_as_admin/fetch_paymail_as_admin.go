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
	ctx := context.Background()
	adminAPI, err := wallet.NewAdminAPIWithXPriv(exampleutil.NewDefaultConfig(), examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	xPub, err := adminAPI.CreateXPub(ctx, &commands.CreateUserXpub{
		Metadata: map[string]any{"xpub_key": "xpub_val"},
		XPub:     examples.XPub,
	})
	if err != nil {
		log.Fatal(err)
	}
	exampleutil.Print("[HTTP POST][Step 1] Create xPub - api/v1/admin/users", xPub)

	paymail, err := adminAPI.CreatePaymail(ctx, &commands.CreatePaymail{
		Key:      examples.XPub,
		Address:  exampleutil.RandomPaymail(),
		Metadata: queryparams.Metadata{"key": "value"},
	})
	if err != nil {
		log.Fatal(err)
	}
	exampleutil.Print("[HTTP POST][Step 2] Create Paymail - api/v1/admin/paymails", paymail)

	paymail, err = adminAPI.Paymail(ctx, paymail.ID)
	if err != nil {
		log.Fatal(err)
	}
	exampleutil.Print("[HTTP GET][Step 3] Fetch Paymail - api/v1/admin/paymails", paymail)
}
