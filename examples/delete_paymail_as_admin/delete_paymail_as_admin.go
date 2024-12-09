package main

import (
	"context"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
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

	addr := exampleutil.RandomPaymail()
	paymail, err := adminAPI.CreatePaymail(ctx, &commands.CreatePaymail{
		Key:      examples.XPub,
		Address:  addr,
		Metadata: querybuilders.Metadata{"key": "value"},
	})
	if err != nil {
		log.Fatal(err)
	}
	exampleutil.Print("[HTTP POST][Step 2] Create Paymail - api/v1/admin/paymails", paymail)

	err = adminAPI.DeletePaymail(ctx, addr)
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print("[HTTP DELETE][Step 3] Delete Paymail: %s - api/v1/admin/paymails", addr)
}
