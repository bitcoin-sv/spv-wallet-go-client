package main

import (
	"context"
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
		buxclient.WithHTTP("localhost:3001"),
		buxclient.WithSignRequest(true),
	)

	_ = buxClient.NewXpub(
		context.Background(), keys.XPub().String(), &buxmodels.Metadata{"example_field": "example_data"},
	)

}
