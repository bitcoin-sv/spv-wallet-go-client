package main

import (
	"fmt"

	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
)

func main() {
	// Generate keys
	keys, err := xpriv.Generate()
	if err != nil {
		panic(err)
	}

	// Generate keys from mnemonic string
	xpriv3, err := xpriv.FromMnemonic(keys.Mnemonic())
	if err != nil {
		panic(err)
	}

	fmt.Println("<-- FromMnemonic method")
	fmt.Println("XPriv: ", xpriv3.XPriv())
	fmt.Println("XPub: ", xpriv3.XPub().String())
	fmt.Println("Mnemonic: ", xpriv3.Mnemonic())

	// Generate keys from string
	xpriv2, err := xpriv.FromString(keys.XPriv())
	if err != nil {
		panic(err)
	}

	fmt.Println("<-- FromString method")
	fmt.Println("XPriv: ", xpriv2.XPriv())
	fmt.Println("XPub: ", xpriv2.XPub().String())
	fmt.Println("Can not get mnemonic from keys generated from string")
}
