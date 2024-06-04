package main

import (
	"fmt"
	"os"

	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
)

func main() {
	// This is an example mnemonic phrase - replace it with your own
	const mnemonicPhrase = "nut same spike popular already mercy kit board rent light illegal local eight filter tube"

	keys, err := xpriv.FromMnemonic(mnemonicPhrase)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("extracted xPriv: ", keys.XPriv())
}
