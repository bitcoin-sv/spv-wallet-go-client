/*
Package examples - Utility functions for this package
*/
package examples

import (
	"fmt"
	"os"
)

func printMissingKeyError(key string) {
	fmt.Printf("Please provide a valid %s. ", key)
}

// HandlePanic - function used to handle a recovery after a panic - use with defer
func HandlePanic() {
	r := recover()

	if r != nil {
		fmt.Println("Recovering: ", r)
	}
}

// CheckIfXPrivExists - checks if ExampleXPriv is not empty
func CheckIfXPrivExists() {
	if ExampleXPriv == "" {
		printMissingKeyError("xPriv")
		os.Exit(1)
	}
}

// CheckIfXPubExists - checks if ExampleXPub is not
func CheckIfXPubExists() {
	if ExampleXPub == "" {
		printMissingKeyError("xPub")
		os.Exit(1)
	}
}

// CheckIfAdminKeyExists - checks if ExampleAdminKey is not empty
func CheckIfAdminKeyExists() {
	if ExampleAdminKey == "" {
		printMissingKeyError("adminKey")
		os.Exit(1)
	}
}
