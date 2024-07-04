/*
Package main - handle_exceptions example
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet/models"
)

func main() {
	defer examples.HandlePanic()

	examples.CheckIfXPubExists()

	const server = "http://localhost:3003/v1"

	client := walletclient.NewWithXPub(server, examples.ExampleXPub)
	ctx := context.Background()

	status, err := client.AdminGetStatus(ctx)
	if err != nil {
		var extendedErr models.ExtendedError
		if errors.As(err, &extendedErr) {
			fmt.Printf("Extended error: [%d] '%s': %s\n", extendedErr.GetStatusCode(), extendedErr.GetCode(), extendedErr.GetMessage())
		} else {
			fmt.Println("Error: ", err.Error())
		}

		os.Exit(1)
	}

	fmt.Println("Status: ", status)
}
