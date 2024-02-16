package main

import (
	"fmt"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
)

func main() {
	// Generate keys
	keys, _ := xpriv.Generate()

	// Create a client
	client, _ := walletclient.New(
		walletclient.WithXPriv(keys.XPriv()),
		walletclient.WithHTTP("localhost:3001"),
		walletclient.WithSignRequest(true),
	)
	fmt.Println(client.IsSignRequest())
}
