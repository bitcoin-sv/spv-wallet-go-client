package main

import (
	"context"
	"fmt"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
)

func main() {
	adminAPI, err := wallet.NewAdminAPIWithXPub(exampleutil.NewDefaultConfig(), examples.AdminXPriv)
	if err != nil {
		log.Fatalf("Failed to initialize admin API with XPriv: %v", err)
	}

	ctx := context.Background()
	err = adminAPI.DeletePaymail(ctx, examples.Paymail)
	if err != nil {
		log.Fatalf("Failed to delete paymail: %v", err)
	}

	fmt.Printf("Paymail deleted: %s\n", examples.Paymail)
}
