package main

import (
	"context"
	"fmt"
	"os"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"

	"examples"
)

func main() {
	defer examples.HandlePanic()

	const server = "http://localhost:3003/v1"

	if examples.ExampleXPriv == "" {
		fmt.Println(examples.ErrMessage("xPriv"))
		os.Exit(1)
	}

	client := walletclient.NewWithXPriv(server, examples.ExampleXPriv)
	ctx := context.Background()

	xpubInfo, err := client.GetXPub(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Current balance: ", xpubInfo.CurrentBalance)
}
