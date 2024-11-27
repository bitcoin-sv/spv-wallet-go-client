package main

import (
	"context"
	"fmt"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
	"github.com/bitcoin-sv/spv-wallet/models"
)

func main() {
	usersAPI, err := wallet.NewUserAPIWithXPriv(exampleutil.ExampleConfig, examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	paymail := "john.doe@example.com"
	code := "f22b4214-ab56-45c0-8399-60ed3a4ecf8e"
	err = usersAPI.ConfirmContact(context.Background(), &models.Contact{ID: "b2215c13-5690-469e-868f-e7bc240a0a23"}, code, paymail, 1, 8)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("\n[HTTP POST] Confirm contact - api/v1/contacts/%s/confirmation", paymail))
}
