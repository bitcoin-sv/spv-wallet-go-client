/*
Package main - xpub_from_xpriv example
*/
package main

import (
	"fmt"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"os"

	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
)

func main() {
	// This is an example xPriv key - replace it with your own
	const xPriv = "xprv9s21ZrQH143K4VneY3UWCF1o5Kk2tmgGrGtMtsrThCTsHsszEZ6H1iP37ZTwuUBvMwudG68SRkcfTjeu8h3rkayfyqkjKAStFBkuNsBnAkS"

	keys, err := xpriv.FromString(xPriv)
	if err != nil {
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}

	fmt.Println("extracted xPub: ", keys.XPub().String())
}
