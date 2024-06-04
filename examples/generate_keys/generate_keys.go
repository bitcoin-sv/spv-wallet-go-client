package main

import (
	"examples"
	"fmt"
)

func main() {
	keys := examples.GenerateKeys()
	exampleXPriv := keys.XPriv()
	exampleXPub := keys.XPub().String()

	fmt.Println("exampleXPriv: ", exampleXPriv)
	fmt.Println("exampleXPub: ", exampleXPub)
}
