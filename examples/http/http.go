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
	wc, _ := walletclient.NewWithXPriv(keys.XPriv(), "localhost:3001")
	fmt.Println(wc.IsSignRequest())
}
