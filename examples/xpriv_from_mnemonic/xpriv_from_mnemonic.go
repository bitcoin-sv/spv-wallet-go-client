/*
Package main - xpriv_from_mnemonic example
*/
package main

import (
	"fmt"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"os"

	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
)

func main() {
	// This is an example mnemonic phrase - replace it with your own
	const mnemonicPhrase = "nut same spike popular already mercy kit board rent light illegal local eight filter tube"

	keys, err := xpriv.FromMnemonic(mnemonicPhrase)
	if err != nil {
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}

	fmt.Println("extracted xPriv: ", keys.XPriv())
}
