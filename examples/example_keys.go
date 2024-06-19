/*
Package examples - key constants to be used in the examples and utility function for generating keys
*/
package examples

import (
	"fmt"
	"os"

	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
)

const (
	// ExampleAdminKey - example admin key
	ExampleAdminKey string = "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"

	// you can generate new keys using `task generate-keys`

	// ExampleXPriv - example private key
	ExampleXPriv string = ""
	// ExampleXPub - example public key
	ExampleXPub string = ""

	// ExamplePaymail - example Paymail address
	ExamplePaymail string = ""
)

// GenerateKeys - function for generating keys (private and public)
func GenerateKeys() xpriv.KeyWithMnemonic {
	keys, err := xpriv.Generate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return keys
}
