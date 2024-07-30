/*
Package main - update_xpub_metadata example
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

	examples.CheckIfXPrivExists()

	const server = "http://localhost:3003/api/v1"

	client := walletclient.NewWithXPriv(server, examples.ExampleXPriv)
	ctx := context.Background()

	xpubInfo, err := client.GetXPub(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("XPub metadata: ", xpubInfo.Metadata)
	fmt.Println("XPub (updated_at): ", xpubInfo.UpdatedAt)

	metadata := map[string]any{"some_metadata_2": "example2"}
	updatedXpubInfo, err := client.UpdateXPubMetadata(ctx, metadata)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Updated XPub metadata: ", updatedXpubInfo.Metadata)
	fmt.Println("Updated XPub (updated_at): ", updatedXpubInfo.UpdatedAt)
}
