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
	usersAPI, err := wallet.NewUserAPIWithXPriv(exampleutil.ExampleConfig, examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	paymail := "john.doe@example.com"
	contact, err := usersAPI.ContactWithPaymail(context.Background(), paymail)
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print(fmt.Sprintf("[HTTP GET] Contact by paymail - api/v1/contacts/%s", paymail), contact)
}
