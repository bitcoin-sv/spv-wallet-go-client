/*
Package examples - Utility functions for this package
*/
package examples

import (
	"fmt"
	"os"
)

// 'xPriv' | 'xPub' | 'adminKey' | 'Paymail'
type keyType string

/*
ErrMessage - function for displaying errors about missing keys (see CheckIfAdminKeyExists, CheckIfXPrivExists)
*/
func ErrMessage(key keyType) string {
	return fmt.Sprintf("Please provide a valid %s.", key)
}

/*
HandlePanic - function used to handle a recovery after a panic - use with defer
*/
func HandlePanic() {
	r := recover()

	if r != nil {
		fmt.Println("Recovering: ", r)
	}
}

/*
CheckIfXPrivExists - checks if ExampleXPriv is not empty
*/
func CheckIfXPrivExists() {
	if ExampleXPriv == "" {
		fmt.Println(ErrMessage("xPriv"))
		os.Exit(1)
	}
}

/*
CheckIfXPubExists - checks if ExampleXPub is not empty
*/
func CheckIfXPubExists() {
	if ExampleXPub == "" {
		fmt.Println(ErrMessage("xPub"))
		os.Exit(1)
	}
}

/*
CheckIfAdminKeyExists - checks if ExampleAdminKey is not empty
*/
func CheckIfAdminKeyExists() {
	if ExampleAdminKey == "" {
		fmt.Println(ErrMessage("adminKey"))
		os.Exit(1)
	}
}
