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
	wc := walletclient.NewWithXPriv(keys.XPriv(), "localhost:3001")
	wc.AdminCreatePaymail(context.Background(), keys.XPub().String(), "foo@domain.com", "", "Foo")
}
