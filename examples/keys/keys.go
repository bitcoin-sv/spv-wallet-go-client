package main

import (
	"fmt"
	"github.com/BuxOrg/go-buxclient/utils"
)

func main() {
	// Generate a random set of keys
	hdXpriv, hdXpub := utils.GenerateRandomSetOfKeys()

	fmt.Println("Your XPriv: ", hdXpriv)
	fmt.Println("Your XPub: ", hdXpub)

	xpub, err := utils.GetPublicKeyFromHDPrivateKey(hdXpriv.String())
	if err != nil {
		panic(err)
	}

	// hdXpub and xpub should be this same
	if hdXpub.String() != xpub.String() {
		panic("xpub and hdXpub are not the same")
	} else {
		fmt.Println("Keys generated successfully, xpub and hdXpub are the same")
	}
}
