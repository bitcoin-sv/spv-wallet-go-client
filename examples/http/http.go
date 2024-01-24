package main

import (
	"fmt"
	"github.com/BuxOrg/go-buxclient"
	"github.com/BuxOrg/go-buxclient/xpriv"
)

func main() {

	// Generate keys
	keys, _ := xpriv.Generate()

	// Create a client
	client, _ := buxclient.New(
		buxclient.WithXPriv(keys.XPriv()),
		buxclient.WithHTTP("localhost:3001"),
		buxclient.WithSignRequest(true),
	)
	fmt.Println(client.IsSignRequest())
}
