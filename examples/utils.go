package examples

import (
	"fmt"
	"os"

	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
)

// 'xPriv' | 'xPub' | 'adminKey' | 'Paymail'
type keyType string

func ErrMessage(key keyType) string {
	return fmt.Sprintf("Please provide a valid %s.", key)
}

func GenerateKeys() xpriv.KeyWithMnemonic {
	keys, err := xpriv.Generate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return keys
}

func HandlePanic() {
	r := recover()

	if r != nil {
		fmt.Println("Recovering: ", r)
	}
}
