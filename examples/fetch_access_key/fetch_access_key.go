package main

import (
	"context"
	"fmt"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
)

func main() {
	usersAPI, err := wallet.NewUserAPIWithXPriv(exampleutil.NewDefaultConfig(), examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	accessKeyID := "35465782-e247-42dd-a2e7-a01ba5b56285"
	accessKey, err := usersAPI.AccessKey(context.Background(), accessKeyID)
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print(fmt.Sprintf("[HTTP GET] Access key - api/v1/users/current/keys/%s", accessKeyID), accessKey)
}
