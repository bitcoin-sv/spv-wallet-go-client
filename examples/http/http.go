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
	wc, _ := walletclient.NewWalletClientWithXPrivate(keys.XPriv(), "localhost:3001", true)
	fmt.Println(wc.IsSignRequest())
}
