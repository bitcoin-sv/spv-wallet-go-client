package main

import (
	"context"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
)

func main() {
	// Replace with your admin keys
	keys, _ := xpriv.Generate()

	// Create a client
	walletclient, _ := walletclient.New(
		walletclient.WithXPriv(keys.XPriv()),
		walletclient.WithHTTP("localhost:3001"),
		walletclient.WithSignRequest(true),
	)

	walletclient.AdminCreatePaymail(context.Background(), keys.XPub().String(), "foo@domain.com", "", "Foo")
}
