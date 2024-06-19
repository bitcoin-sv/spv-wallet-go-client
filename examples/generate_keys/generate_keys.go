/*
Package main - generate_keys example
*/
package main

import (
	"fmt"

	"examples"
)

func main() {
	keys := examples.GenerateKeys()
	exampleXPriv := keys.XPriv()
	exampleXPub := keys.XPub().String()

	fmt.Println("exampleXPriv: ", exampleXPriv)
	fmt.Println("exampleXPub: ", exampleXPub)
}
