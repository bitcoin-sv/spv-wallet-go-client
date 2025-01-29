package managecontacts

import (
	"context"
	"fmt"
	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"log"
)

type internalCfg struct {
	totpDigits    uint
	totpPeriods   uint
	server        string
	paymailDomain string
}

type creds struct {
	xPriv string
	xPub  string
}

type verificationResults struct {
	bobValidatedAlicesTotp bool
	aliceValidatedBobsTotp bool
}

type clients struct {
	Alice *wallet.UserAPI
	Bob   *wallet.UserAPI
	Admin *wallet.AdminAPI
}

func newConfig() *internalCfg {
	return &internalCfg{
		totpDigits:  2,
		totpPeriods: 1200,
		server:      "http://localhost:3003",
		// Replace with your own paymail domain!
		paymailDomain: "example.com",
	}
}

func validateConfig(cfg *internalCfg) error {
	if cfg.paymailDomain == "" || cfg.paymailDomain == "example.com" {
		return fmt.Errorf("please replace the paymail domain with your own domain")
	}

	return nil
}

func newCredentials() map[string]creds {
	return map[string]creds{
		"alice": {
			xPriv: "xprv9s21ZrQH143K2jMwweKF33hFDDvwxEooDtXbZ7mGTJQfmSs8aD77ThuYDsfNrgBAbHr9Yx8FrPaukMLHpxFUyyvBuzAJBMpd4a2xFxr6qts",
			xPub:  "xpub661MyMwAqRbcFDSR3frFQBdymFmSMhXeb7TCMWAt1dweeFCH7kRN1WE257E65MufrqngaLK46ERg5LHHouHiS8DvHKovmo5VhjLs5vgwqdp",
		},
		"bob": {
			xPriv: "xprv9s21ZrQH143K3DkTDsWwvUb3pwgKoYGp9hxYe2coqZz3pvE1kQfe1dQLdcN82XSeLmw1nGpMZLnXZktf9hFJTu9NRLBpQnGHwYpo4SmszZY",
			xPub:  "xpub661MyMwAqRbcFhpvKu3xHcXnNyWpCzzfWvt9SR2RPuX2hiZAHwytZRipUtM4qG2PPPF5pZttP3grZM9N9MR5jSek7RRgyggsLJAWFJJUAko",
		},
	}
}

func newClients(creds map[string]creds) (*clients, error) {
	alice, err := wallet.NewUserAPIWithXPriv(exampleutil.NewDefaultConfig(), creds["alice"].xPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Alice's API: %w", err)
	}

	bob, err := wallet.NewUserAPIWithXPriv(exampleutil.NewDefaultConfig(), creds["bob"].xPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Bob's API: %w", err)
	}

	admin, err := wallet.NewAdminAPIWithXPriv(exampleutil.NewDefaultConfig(), examples.AdminXPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize admin API: %w", err)
	}

	return &clients{
		Alice: alice,
		Bob:   bob,
		Admin: admin,
	}, nil
}

func getPaymail(name, domain string) string {
	return fmt.Sprintf("%s@%s", name, domain)
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

func setupUsers(ctx context.Context, clients *clients, cfg *internalCfg, creds map[string]creds) error {
	fmt.Println("0. Setting up users (optional - uncomment if users are not registered)")

	_, err := clients.Admin.CreateXPub(ctx, &commands.CreateUserXpub{
		XPub: creds["alice"].xPub,
	})
	if err != nil {
		return fmt.Errorf("failed to create Alice's xPub: %w", err)
	}

	_, err = clients.Admin.CreatePaymail(ctx, &commands.CreatePaymail{
		Key:     creds["alice"].xPub,
		Address: getPaymail("alice", cfg.paymailDomain),
	})
	if err != nil {
		return fmt.Errorf("failed to create Alice's paymail: %w", err)
	}

	_, err = clients.Admin.CreateXPub(ctx, &commands.CreateUserXpub{
		XPub: creds["bob"].xPub,
	})
	if err != nil {
		return fmt.Errorf("failed to create Bob's xPub: %w", err)
	}

	_, err = clients.Admin.CreatePaymail(ctx, &commands.CreatePaymail{
		Key:     creds["bob"].xPub,
		Address: getPaymail("bob", cfg.paymailDomain),
	})
	if err != nil {
		return fmt.Errorf("failed to create Bob's paymail: %w", err)
	}

	return nil
}

func verificationFlow(ctx context.Context, clients *clients, cfg *internalCfg) (*verificationResults, error) {
	fmt.Println("1. Creating initial contacts")

	alicePaymail := getPaymail("alice", cfg.paymailDomain)
	bobPaymail := getPaymail("bob", cfg.paymailDomain)

	_, err := clients.Alice.UpsertContact(ctx, commands.UpsertContact{
		ContactPaymail:   bobPaymail,
		FullName:         "Bob Smith",
		RequesterPaymail: alicePaymail,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Bob's contact for Alice: %w", err)
	}

	_, err = clients.Bob.UpsertContact(ctx, commands.UpsertContact{
		ContactPaymail:   alicePaymail,
		FullName:         "Alice Smith",
		RequesterPaymail: bobPaymail,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Alice's contact for Bob: %w", err)
	}

	respBob, err := clients.Alice.ContactWithPaymail(ctx, bobPaymail)
	if err != nil {
		return nil, fmt.Errorf("failed to get Bob's contact: %w", err)
	}
	bobContact := mapToContactModel(respBob)

	respAlice, err := clients.Bob.ContactWithPaymail(ctx, alicePaymail)
	if err != nil {
		return nil, fmt.Errorf("failed to get Alice's contact: %w", err)
	}
	aliceContact := mapToContactModel(respAlice)

	fmt.Println("\n2. Alice initiates verification")
	aliceTotpForBob, err := clients.Alice.GenerateTotpForContact(bobContact, cfg.totpPeriods, cfg.totpDigits)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Alice's TOTP for Bob: %w", err)
	}
	logSecureMessage("Alice", "Bob", aliceTotpForBob)

	fmt.Println("3. Bob validates Alice's TOTP")
	bobValidationErr := clients.Bob.ValidateTotpForContact(aliceContact, aliceTotpForBob, respBob.Paymail, cfg.totpPeriods, cfg.totpDigits)
	bobValidatedAlicesTotp := bobValidationErr == nil
	fmt.Printf("Validation status: %v\n", bobValidatedAlicesTotp)

	fmt.Println("\n4. Bob initiates verification")
	bobTotpForAlice, err := clients.Bob.GenerateTotpForContact(aliceContact, cfg.totpPeriods, cfg.totpDigits)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Bob's TOTP for Alice: %w", err)
	}
	logSecureMessage("Bob", "Alice", bobTotpForAlice)

	fmt.Println("5. Alice validates Bob's TOTP")
	aliceValidationErr := clients.Alice.ValidateTotpForContact(bobContact, bobTotpForAlice, respAlice.Paymail, cfg.totpPeriods, cfg.totpDigits)
	aliceValidatedBobsTotp := aliceValidationErr == nil
	fmt.Printf("Validation status: %v\n", aliceValidatedBobsTotp)

	return &verificationResults{
		bobValidatedAlicesTotp: bobValidatedAlicesTotp,
		aliceValidatedBobsTotp: aliceValidatedBobsTotp,
	}, nil
}

func finalizeAndCleanup(ctx context.Context, clients *clients, cfg *internalCfg, results *verificationResults) error {
	isFullyVerified := results.bobValidatedAlicesTotp && results.aliceValidatedBobsTotp
	fmt.Printf("\nBidirectional verification complete: %v\n", isFullyVerified)

	if isFullyVerified {
		fmt.Println("\n6. Admin confirms verified contacts")
		if err := clients.Admin.ConfirmContacts(ctx, &commands.ConfirmContacts{
			PaymailA: getPaymail("alice", cfg.paymailDomain),
			PaymailB: getPaymail("bob", cfg.paymailDomain),
		}); err != nil {
			_ = fmt.Errorf("failed to confirm contacts: %w", err)
		}
	}

	fmt.Println("\n7. Cleaning up contacts")
	if err := clients.Alice.RemoveContact(ctx, getPaymail("bob", cfg.paymailDomain)); err != nil {
		return fmt.Errorf("failed to remove Bob's contact: %w", err)
	}

	if err := clients.Bob.RemoveContact(ctx, getPaymail("alice", cfg.paymailDomain)); err != nil {
		return fmt.Errorf("failed to remove Alice's contact: %w", err)
	}

	return nil
}

func main() {
	internalCfg := newConfig()
	if err := validateConfig(internalCfg); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	creds := newCredentials()
	clients, err := newClients(creds)
	if err != nil {
		log.Fatalf("Failed to initialize clients: %v", err)
	}

	ctx := context.Background()

	fmt.Println("We assume that the users: Alice and Bob are already registered.")
	fmt.Println("If they're not, please uncomment the setupUsers() call below.")

	// Uncomment to setup users
	// if err := setupUsers(ctx, clients, internalCfg, creds); err != nil {
	//     log.Fatalf("Failed to setup users: %v", err)
	// }

	results, err := verificationFlow(ctx, clients, internalCfg)
	if err != nil {
		log.Fatalf("Error during verification flow: %v", err)
	}

	if err := finalizeAndCleanup(ctx, clients, internalCfg, results); err != nil {
		log.Fatalf("Error during cleanup: %v", err)
	}
}
