/*
Package main - generate_keys example
*/
package main

import (
	"fmt"
	"os"

	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
)

func main() {
	keys, err := xpriv.Generate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	exampleXPriv := keys.XPriv()
	exampleXPub := keys.XPub().String()

	fmt.Println("exampleXPriv: ", exampleXPriv)
	fmt.Println("exampleXPub: ", exampleXPub)
}
