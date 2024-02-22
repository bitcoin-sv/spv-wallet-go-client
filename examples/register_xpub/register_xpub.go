package main

import (
	"context"
	"fmt"
	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet/models"
)

func main() {
	// Replace with your admin keys
	keys, _ := xpriv.Generate()

	// Create a client
	walletClient, _ := walletclient.New(
		walletclient.WithXPriv(keys.XPriv()),
		walletclient.WithHTTP("localhost:3003/v1"),
		walletclient.WithSignRequest(true),
	)

	ctx := context.Background()

	_ = walletClient.AdminNewXpub(ctx, keys.XPub().String(), &models.Metadata{"example_field": "example_data"})

	xpubKey, err := walletClient.GetXPub(ctx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(xpubKey)
}
