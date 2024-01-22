package main

import (
	"context"
	"fmt"
	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/go-buxclient"
	"github.com/BuxOrg/go-buxclient/xpriv"
)

func main() {
	// Generate keys
	keys, _ := xpriv.Generate()

	// Create a client
	buxClient, _ := buxclient.New(
		buxclient.WithXPriv(keys.XPriv()),
		buxclient.WithHTTP("localhost:3003/v1"),
		buxclient.WithSignRequest(true),
	)

	ctx := context.Background()

	_ = buxClient.NewXpub(
		ctx, keys.XPub().String(), &buxmodels.Metadata{"example_field": "example_data"},
	)

	xpubKey, err := buxClient.GetXPub(ctx)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(xpubKey)

}
