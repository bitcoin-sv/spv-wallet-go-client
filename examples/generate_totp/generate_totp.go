package main

import (
	"fmt"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
	"github.com/bitcoin-sv/spv-wallet/models"
)

func main() {
	const aliceXPriv = examples.UserXPriv

	// pubKey - PKI can be obtained from the contact's paymail capability
	const bobPKI = "03a48e13dc598dce5fda9b14ea13f32d5dbc4e8d8a34447dda84f9f4c457d57fe7"
	const digits = 4
	const period = 1200

	alice, err := wallet.NewUserAPIWithXPriv(exampleutil.NewDefaultConfig(), aliceXPriv)
	if err != nil {
		log.Fatalf("Failed to initialize user API with XPriv: %v", err)
	}

	bob := &models.Contact{
		PubKey:  bobPKI,
		Paymail: "test@paymail.com",
	}
	code, err := alice.GenerateTotpForContact(bob, period, digits)
	if err != nil {
		log.Fatalf("Failed to generate totp for contact: %v", err)
	}

	fmt.Println("TOTP code from Alice to Bob: ", code)

	err = alice.ValidateTotpForContact(bob, code, bob.Paymail, period, digits)
	if err != nil {
		log.Fatalf("Failed to validate totp for contact: %v", err)
	}

	fmt.Println("TOTP code from Alice to Bob is valid")
}
