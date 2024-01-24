package xpriv

import "fmt"

func ExampleGenerate() {
	keys, _ := Generate()

	fmt.Println("xpriv:", keys.XPriv())
	fmt.Println("xpub:", keys.XPub().String())
}
