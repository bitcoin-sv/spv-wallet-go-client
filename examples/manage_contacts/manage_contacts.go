package main

import (
	"context"
	"fmt"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	connectionConfig "github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// !!! Adjust the server url
const server = "http://localhost:3003"

// !!! Adjust the paymail domain to the domain supported by the spv-wallet server
const yourPaymailDomain = "example.com"

// We assume that the users: Alice and Bob are already registered.
// If they're not, please set this to true to make the example create them.
const setupUsers = false

// Example configuration – adjust as needed.
// It holds the values required to present the example.
var config = struct {
	totpDigits    uint
	totpPeriods   uint
	paymailDomain string
	alice         user
	bob           user
}{
	totpDigits:    2,
	totpPeriods:   1200,
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

var conConfig = connectionConfig.New(connectionConfig.WithAddr(server))

var clients = struct {
	alice *wallet.UserAPI
	bob   *wallet.UserAPI
	admin *wallet.AdminAPI
}{
	alice: assertNoError(wallet.NewUserAPIWithXPriv(conConfig, config.alice.xPriv)),
	bob:   assertNoError(wallet.NewUserAPIWithXPriv(conConfig, config.bob.xPriv)),
	admin: assertNoError(wallet.NewAdminAPIWithXPriv(conConfig, examples.AdminXPriv)),
}

var ctx = context.Background()

func createInitialContacts() error {
	// Initial Contacts Creation: Step 1  Create Alice's contact with Bob
	fmt.Println("\n Creating Alice's contact with Bob")
	_, err := clients.alice.UpsertContact(ctx, commands.UpsertContact{
		ContactPaymail:   config.bob.paymail,
		FullName:         "Bob Smith",
		RequesterPaymail: config.alice.paymail,
	})
	if err != nil {
		return fmt.Errorf("failed to create Alice's contact with Bob: %w", err)
	}

	// Initial Contacts Creation: Step 2  Create Bob's contact with Alice
	fmt.Println("\n Creating Bob's contact with Alice")
	_, err = clients.bob.UpsertContact(ctx, commands.UpsertContact{
		ContactPaymail:   config.alice.paymail,
		FullName:         "Alice Smith",
		RequesterPaymail: config.bob.paymail,
	})
	if err != nil {
		return fmt.Errorf("failed to create Bob's contact with Alice: %w", err)
	}

	return nil
}

func initiateContactVerification() error {
	// Verification: Step 1  Alice generate TOTP
	fmt.Println("\n Alice generate TOTP for Bob")
	aliceTotpForBob, err := generateTOTP(clients.alice, config.bob)
	if err != nil {
		return fmt.Errorf("generating TOTP by Alice failed: %w", err)
	}

	logSecureMessage("Alice", "Bob", aliceTotpForBob)

	// Verification: Step 2  Bob validates Alice's TOTP
	fmt.Println("\n Bob validates Alice's TOTP")
	err = validateTOTP(clients.bob, config.bob.paymail, config.alice.paymail, aliceTotpForBob)
	if err != nil {
		return fmt.Errorf("validating TOTP by Bob failed: %w", err)
	}

	// As err here is nil, Alice's contact is validated
	aliceValidated := true

	// Verification: Step 3  Bob generate TOTP
	fmt.Println("\n Bob generate TOTP for Alice")
	bobTotpForAlice, err := generateTOTP(clients.bob, config.alice)
	if err != nil {
		return fmt.Errorf("generating TOTP by Bob failed: %w", err)
	}

	logSecureMessage("Bob", "Alice", bobTotpForAlice)

	// Verification: Step 4  Alice validates Bob's TOTP
	fmt.Println("\n Alice validates Bob's TOTP")
	err = validateTOTP(clients.alice, config.alice.paymail, config.bob.paymail, bobTotpForAlice)
	if err != nil {
		return fmt.Errorf("validating TOTP by Alice failed: %w", err)
	}

	// As err here is nil, Bob's contact is validated
	bobValidated := true

	if bobValidated && aliceValidated {
		fmt.Println("Both TOTP verifications succeeded.")
	}
	return nil
}

func cleanup() error {
	// Cleanup: Step 1  Admin deletes Alice's contact with Bob
	fmt.Println("\n Deleting contact and unconfirm other side")
	aliceToBobContact, err := adminGetContact(config.alice.xPub, config.bob.paymail)
	if err != nil {
		return fmt.Errorf("failed to get Alice's contact with Bob: %w", err)
	}

	fmt.Println("\n Admin deletes Alice's contact with Bob")
	err = clients.admin.DeleteContact(context.Background(), aliceToBobContact.ID)
	if err != nil {
		return fmt.Errorf("failed to delete Alice's contact with Bob: %w", err)
	}

	// Cleanup: Step 2  Admin unconfirms Bob's contact with Alice
	bobToAliceContact, err := adminGetContact(config.bob.xPub, config.alice.paymail)
	if err != nil {
		return fmt.Errorf("failed to get Bob's contact with Alice: %w", err)
	}

	fmt.Println("\n Admin unconfirms Bob's contact with Alice")
	err = clients.admin.UnconfirmContact(context.Background(), bobToAliceContact.ID)
	if err != nil {
		return fmt.Errorf("failed to unconfirm Bob's contact with Alice: %w", err)
	}

	// Cleanup: Step 3  Bob removes contact with Alice
	fmt.Println("\n Bob removes contact with Alice")
	if err := clients.bob.RemoveContact(ctx, config.alice.paymail); err != nil {
		return fmt.Errorf("failed to remove contact with Alice: %w", err)
	}

	return nil
}

func adminGetContact(creatorOfContactXPub, counterpartyPaymail string) (*response.Contact, error) {
	creatorOfContactXPubID := exampleutil.CreateXpubID(creatorOfContactXPub)
	response, err := clients.admin.Contacts(
		context.Background(),
		queries.QueryOption[filter.AdminContactFilter](
			queries.QueryWithFilter(filter.AdminContactFilter{
				ContactFilter: filter.ContactFilter{
					Paymail: &counterpartyPaymail,
				},
				XPubID: &creatorOfContactXPubID,
			}),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get contact %w", err)
	}
	contact := response.Content[0]

	return contact, nil
}

func generateTOTP(initiatorClient *wallet.UserAPI, counterparty user) (totp string, err error) {
	// Contact initiator gets the contact info
	response, err := initiatorClient.ContactWithPaymail(ctx, counterparty.paymail)
	if err != nil {
		return "", fmt.Errorf("failed to get contact from initiator's perspective: %w", err)
	}
	contact := mapToContactModel(response)

	// Contact initiator generates TOTP
	totp, err = initiatorClient.GenerateTotpForContact(contact, config.totpPeriods, config.totpDigits)
	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP: %w", err)
	}

	return totp, nil
}

func validateTOTP(validatorClient *wallet.UserAPI, validatorPaymail string, generatingUserPaymail string, totp string) error {
	// Contact counterparty gets the contact info
	respCounterparty, err := validatorClient.ContactWithPaymail(ctx, generatingUserPaymail)
	if err != nil {
		return fmt.Errorf("failed to get contact from counterparty's perspective: %w", err)
	}

	contact := mapToContactModel(respCounterparty)

	// Contact counterparty validates TOTP
	validationErr := validatorClient.ValidateTotpForContact(
		contact,
		totp,
		validatorPaymail,
		config.totpPeriods,
		config.totpDigits,
	)
	if validationErr != nil {
		fmt.Printf("[WARN] TOTP validation failed: %v\n", validationErr)
		return validationErr
	}

	return nil
}

func main() {
	if setupUsers {
		// Initiate users
		fmt.Println("\n======== Initiating users ========")
		initiateUsers()
	} else {
		fmt.Println("We assume that the users: Alice and Bob are already registered.")
		fmt.Println("If they're not, please set config.SetupUsers to true.")
	}

	// Create initial contacts
	fmt.Println("\n======== Creating initial contacts ========")
	err := createInitialContacts()
	if err != nil {
		log.Fatalf("Error during initial contacts creation: %v", err)
	}

	// Initiate contact verification
	fmt.Println("\n======== Initiating contact verification ========")
	if err := initiateContactVerification(); err != nil {
		log.Fatalf("Error during verification flow: %v", err)
	}

	// Confirm contacts using the admin API
	fmt.Println("\n======== Confirming contacts ========")
	err = clients.admin.ConfirmContacts(ctx, &commands.ConfirmContacts{
		PaymailA: config.alice.paymail,
		PaymailB: config.bob.paymail,
	})
	if err != nil {
		log.Fatalf("Error during confirmation of contacts: %v", err)
	}

	// Clean up
	fmt.Println("\n======== Cleaning up ========")
	if err := cleanup(); err != nil {
		log.Fatalf("Error during cleanup: %v", err)
	}
}
