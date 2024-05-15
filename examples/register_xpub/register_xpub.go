package main

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/models"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
)

func main() {
	// Replace with your admin keys
	keys, _ := xpriv.Generate()

	// Create a client
	wc, _ := walletclient.NewWithXPriv(keys.XPriv(), "localhost:3003/v1")
	ctx := context.Background()
	_ = wc.AdminNewXpub(ctx, keys.XPub().String(), &models.Metadata{"example_field": "example_data"})

	xpubKey, err := wc.GetXPub(ctx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(xpubKey)
}
