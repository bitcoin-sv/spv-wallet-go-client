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
	err = adminAPI.DeletePaymail(ctx, "d43ed481ba08aae1db02d880ebefe962f9796168387bb293a95024cb02b953ef")
	if err != nil {
		log.Fatalf("Failed to delete paymail: %v", err)
	}

	fmt.Printf("Paymail deleted: %s\n", examples.Paymail)
}
