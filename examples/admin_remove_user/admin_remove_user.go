/*
Package main - admin_remove_user example
*/
package main

import (
	"context"
	"fmt"
	"os"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
)

func main() {
	defer examples.HandlePanic()

	examples.CheckIfAdminKeyExists()

	const server = "http://localhost:3003/v1"

	adminClient := walletclient.NewWithAdminKey(server, examples.ExampleAdminKey)
	ctx := context.Background()

	err := adminClient.AdminDeletePaymail(ctx, examples.ExamplePaymail)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
