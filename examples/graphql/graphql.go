package main

import (
	"github.com/BuxOrg/go-buxclient"
	"github.com/BuxOrg/go-buxclient/xpriv"
)

func main() {

	// Generate keys
	keys, _ := xpriv.Generate()

	// Create a client
	_, _ = buxclient.New(
		buxclient.WithXPriv(keys.XPriv()),
		buxclient.WithGraphQL("localhost:3001"),
		buxclient.WithSignRequest(true),
	)
}
