package main

import (
	"context"
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

	_ = buxClient.NewPaymail(context.Background(), keys.XPub().String(), "foo@domain.com", "", "Foo", nil)

}
