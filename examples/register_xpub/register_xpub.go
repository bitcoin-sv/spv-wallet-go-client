package main

import (
	"context"
	"fmt"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
	"github.com/bitcoin-sv/spv-wallet/models"
)

func main() {
	// Generate keys
	keys, _ := xpriv.Generate()

	// Create a client
	buxClient, _ := walletclient.New(
		walletclient.WithXPriv(keys.XPriv()),
		walletclient.WithHTTP("localhost:3003/v1"),
		walletclient.WithSignRequest(true),
	)

	ctx := context.Background()

	_ = buxClient.NewXpub(
		ctx, keys.XPub().String(), &models.Metadata{"example_field": "example_data"},
	)

	xpubKey, err := buxClient.GetXPub(ctx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(xpubKey)
}
