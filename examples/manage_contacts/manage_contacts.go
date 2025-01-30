package main

import (
	"context"
	"fmt"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
)

// !!! Adjust the paymail domain to the domain supported by the spv-wallet server
const yourPaymailDomain = "example.com"

// Example configuration – adjust as needed.
// It holds the values required to present the example.
var config = struct {
	setupUsers    bool
	totpDigits    uint
	totpPeriods   uint
	server        string
	paymailDomain string
	alice         user
	bob           user
}{
	// We assume that the users: Alice and Bob are already registered.
	// If they're not, please set this to true to make the example create them.
	setupUsers: false,

	totpDigits:    2,
	totpPeriods:   1200,
	server:        "http://localhost:3003",
	paymailDomain: examplePaymailCorrectlyEdited(yourPaymailDomain),
	alice: user{
		xPriv:   "xprv9s21ZrQH143K2jMwweKF33hFDDvwxEooDtXbZ7mGTJQfmSs8aD77ThuYDsfNrgBAbHr9Yx8FrPaukMLHpxFUyyvBuzAJBMpd4a2xFxr6qts",
		xPub:    "xpub661MyMwAqRbcFDSR3frFQBdymFmSMhXeb7TCMWAt1dweeFCH7kRN1WE257E65MufrqngaLK46ERg5LHHouHiS8DvHKovmo5VhjLs5vgwqdp",
		paymail: "alice" + "@" + yourPaymailDomain,
	},
	bob: user{
		xPriv:   "xprv9s21ZrQH143K3DkTDsWwvUb3pwgKoYGp9hxYe2coqZz3pvE1kQfe1dQLdcN82XSeLmw1nGpMZLnXZktf9hFJTu9NRLBpQnGHwYpo4SmszZY",
		xPub:    "xpub661MyMwAqRbcFhpvKu3xHcXnNyWpCzzfWvt9SR2RPuX2hiZAHwytZRipUtM4qG2PPPF5pZttP3grZM9N9MR5jSek7RRgyggsLJAWFJJUAko",
		paymail: "bob" + "@" + yourPaymailDomain,
	},
}

var clients = struct {
	alice *wallet.UserAPI
	bob   *wallet.UserAPI
	admin *wallet.AdminAPI
}{
	alice: assertNoError(wallet.NewUserAPIWithXPriv(exampleutil.NewDefaultConfig(), config.alice.xPriv)),
	bob:   assertNoError(wallet.NewUserAPIWithXPriv(exampleutil.NewDefaultConfig(), config.bob.xPriv)),
	admin: assertNoError(wallet.NewAdminAPIWithXPriv(exampleutil.NewDefaultConfig(), examples.AdminXPriv)),
}

var ctx = context.Background()

func verificationFlow() (*verificationResults, error) {
	fmt.Println("\n1. Creating initial contacts")

	alicePaymail := config.alice.paymail
	bobPaymail := config.bob.paymail

	_, err := clients.alice.UpsertContact(ctx, commands.UpsertContact{
		ContactPaymail:   bobPaymail,
		FullName:         "Bob Smith",
		RequesterPaymail: alicePaymail,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Bob's contact for Alice: %w", err)
	}

	_, err = clients.bob.UpsertContact(ctx, commands.UpsertContact{
		ContactPaymail:   alicePaymail,
		FullName:         "Alice Smith",
		RequesterPaymail: bobPaymail,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Alice's contact for Bob: %w", err)
	}

	respBob, err := clients.alice.ContactWithPaymail(ctx, bobPaymail)
	if err != nil {
		return nil, fmt.Errorf("failed to get Bob's contact: %w", err)
	}
	bobContact := mapToContactModel(respBob)

	respAlice, err := clients.bob.ContactWithPaymail(ctx, alicePaymail)
	if err != nil {
		return nil, fmt.Errorf("failed to get Alice's contact: %w", err)
	}
	aliceContact := mapToContactModel(respAlice)

	fmt.Println("\n2. Alice initiates verification")
	aliceTotpForBob, err := clients.alice.GenerateTotpForContact(bobContact, config.totpPeriods, config.totpDigits)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Alice's TOTP for Bob: %w", err)
	}
	logSecureMessage("Alice", "Bob", aliceTotpForBob)

	fmt.Println("\n3. Bob validates Alice's TOTP")
	bobValidationErr := clients.bob.ValidateTotpForContact(aliceContact, aliceTotpForBob, respBob.Paymail, config.totpPeriods, config.totpDigits)
	bobValidatedAlicesTotp := bobValidationErr == nil
	fmt.Printf("Validation status: %v\n", bobValidatedAlicesTotp)

	fmt.Println("\n4. Bob initiates verification")
	bobTotpForAlice, err := clients.bob.GenerateTotpForContact(aliceContact, config.totpPeriods, config.totpDigits)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Bob's TOTP for Alice: %w", err)
	}
	logSecureMessage("Bob", "Alice", bobTotpForAlice)

	fmt.Println("\n5. Alice validates Bob's TOTP")
	aliceValidationErr := clients.alice.ValidateTotpForContact(bobContact, bobTotpForAlice, respAlice.Paymail, config.totpPeriods, config.totpDigits)
	aliceValidatedBobsTotp := aliceValidationErr == nil
	fmt.Printf("Validation status: %v\n", aliceValidatedBobsTotp)

	return &verificationResults{
		bobValidatedAlicesTotp: bobValidatedAlicesTotp,
		aliceValidatedBobsTotp: aliceValidatedBobsTotp,
	}, nil
}

func finalizeAndCleanup(results *verificationResults) error {
	isFullyVerified := results.bobValidatedAlicesTotp && results.aliceValidatedBobsTotp
	fmt.Printf("\nBidirectional verification complete: %v\n", isFullyVerified)

	if isFullyVerified {
		fmt.Println("\n6. Admin confirms verified contacts")
		if err := clients.admin.ConfirmContacts(ctx, &commands.ConfirmContacts{
			PaymailA: config.alice.paymail,
			PaymailB: config.bob.paymail,
		}); err != nil {
			_ = fmt.Errorf("failed to confirm contacts: %w", err)
		}
	}

	fmt.Println("\n7. Cleaning up contacts")
	if err := clients.alice.RemoveContact(ctx, config.bob.paymail); err != nil {
		return fmt.Errorf("failed to remove Bob's contact: %w", err)
	}

	if err := clients.bob.RemoveContact(ctx, config.alice.paymail); err != nil {
		return fmt.Errorf("failed to remove Alice's contact: %w", err)
	}

	return nil
}

func main() {
	if config.setupUsers {
		setupUsers()
	} else {
		fmt.Println("We assume that the users: Alice and Bob are already registered.")
		fmt.Println("If they're not, please set config.setupUsers to true to make the example create them.")
	}

	results, err := verificationFlow()
	if err != nil {
		log.Fatalf("Error during verification flow: %v", err)
	}

	if err := finalizeAndCleanup(results); err != nil {
		log.Fatalf("Error during cleanup: %v", err)
	}
}
