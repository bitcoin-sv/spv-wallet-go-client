package main

import (
	"fmt"
	"github.com/BuxOrg/go-buxclient/xpriv"
)

func main() {
	//Generate keys
	keys, err := xpriv.Generate()
	if err != nil {
		panic(err)
	}

	fmt.Println("<-- Generate method")
	fmt.Println("XPriv: ", keys.String())
	fmt.Println("XPub: ", keys.XPub().String())
	fmt.Println("Mnemonic: ", keys.Mnemonic())

	//Generate keys from mnemonic string
	xpriv3, err := xpriv.FromMnemonic(keys.Mnemonic())
	if err != nil {
		panic(err)
	}

	fmt.Println("<-- FromMnemonic method")
	fmt.Println("XPriv: ", xpriv3.String())
	fmt.Println("XPub: ", xpriv3.XPub().String())
	fmt.Println("Mnemonic: ", xpriv3.Mnemonic())

	//Generate keys from string
	xpriv2, err := xpriv.FromString(keys.String())
	if err != nil {
		panic(err)
	}

	fmt.Println("<-- FromString method")
	fmt.Println("XPriv: ", xpriv2.String())
	fmt.Println("XPub: ", xpriv2.XPub().String())
	fmt.Println("Can not get mnemonic from keys generated from string")
}
