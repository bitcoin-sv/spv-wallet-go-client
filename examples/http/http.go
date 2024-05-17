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
	wc := walletclient.NewWithXPriv("https://localhost:3001", keys.XPriv())
	fmt.Println(wc.IsSignRequest())
}
