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
	wc := walletclient.NewWithXPriv("https://localhost:3001", keys.XPriv())
	wc.AdminCreatePaymail(context.Background(), keys.XPub().String(), "foo@domain.com", "", "Foo")
}
