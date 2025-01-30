package main

import (
	"fmt"
	"log"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

type user struct {
	xPriv   string
	xPub    string
	paymail string
}

type verificationResults struct {
	bobValidatedAlicesTotp bool
	aliceValidatedBobsTotp bool
}

func examplePaymailCorrectlyEdited(paymail string) string {
	if paymail == "" || paymail == "example.com" {
		log.Fatal("Invalid configuration - please replace the paymail domain with your own domain")
	}
	return paymail
}

func assertNoError[T any](val T, err error) T {
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
	return val
}

func logSecureMessage(from, to, totp string) {
	fmt.Printf("\n!!! SECURE COMMUNICATION REQUIRED !!!\n%s's TOTP code for %s:\n", from, to)
	fmt.Printf("TOTP code: %s\n", totp)
	fmt.Print("Share using: encrypted message, secure email, phone call or in-person meeting.\n")
}

func mapToContactModel(resp *response.Contact) *models.Contact {
	return &models.Contact{
		ID:       resp.ID,
		FullName: resp.FullName,
		Paymail:  resp.Paymail,
		PubKey:   resp.PubKey,
		Status:   resp.Status,
	}
}

func setupUsers() {
	fmt.Println("0. Setting up users (optional)")

	// Create account for Alice
	assertNoError(clients.admin.CreateXPub(ctx, &commands.CreateUserXpub{
		XPub: config.alice.xPub,
	}))
	assertNoError(clients.admin.CreatePaymail(ctx, &commands.CreatePaymail{
		Key:     config.alice.xPub,
		Address: config.alice.paymail,
	}))

	// Create account for Bob
	assertNoError(clients.admin.CreateXPub(ctx, &commands.CreateUserXpub{
		XPub: config.bob.xPub,
	}))
	assertNoError(clients.admin.CreatePaymail(ctx, &commands.CreatePaymail{
		Key:     config.bob.xPub,
		Address: config.bob.paymail,
	}))
}
