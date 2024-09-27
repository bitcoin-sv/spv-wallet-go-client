/*
Package main - generate_totp example
*/
package main

import (
	"fmt"
	"os"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet/models"
)

func main() {
	defer examples.HandlePanic()

	const server = "http://localhost:3003/v1"
	const aliceXPriv = "xprv9s21ZrQH143K4JFXqGhBzdrthyNFNuHPaMUwvuo8xvpHwWXprNK7T4JPj1w53S1gojQncyj8JhSh8qouYPZpbocsq934cH5G1t1DRBfgbod"
	const bobPKI = "03a48e13dc598dce5fda9b14ea13f32d5dbc4e8d8a34447dda84f9f4c457d57fe7"
	const digits = 4
	const period = 1200 // 20 minutes

	client, err := walletclient.NewWithXPriv(server, aliceXPriv)
	if err != nil {
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}

	mockContact := &models.Contact{
		PubKey:  bobPKI,
		Paymail: "test@paymail.com",
	}

	totpCode, err := client.GenerateTotpForContact(mockContact, period, digits)
	if err != nil {
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}
	fmt.Println("TOTP code from Alice to Bob: ", totpCode)

	valid, err := client.ValidateTotpForContact(mockContact, totpCode, mockContact.Paymail, period, digits)
	if err != nil {
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}
	fmt.Println("Is TOTP code valid: ", valid)
}
