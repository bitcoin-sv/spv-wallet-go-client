//go:build regression
// +build regression

package regressiontests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
)

func TestRegressionWorkflow(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	spvWalletPG, spvWalletSL := initServers(t)

	t.Log("Step 1: Setup success: created SPV client instances with test users")
	t.Logf("SPV clients for env: %s, user: %s, admin: %s, leader: %s", spvWalletPG.cfg.envURL, spvWalletPG.user.alias, spvWalletPG.admin.alias, spvWalletPG.leader.alias)
	t.Logf("SPV clients for env: %s, user: %s, admin: %s, leader: %s", spvWalletSL.cfg.envURL, spvWalletSL.user.alias, spvWalletSL.admin.alias, spvWalletSL.leader.alias)

	t.Run("Step 2: The leader clients attempt to fetch the shared configuration response from their SPV Wallet API instance.", func(t *testing.T) {
		tests := []struct {
			name                string
			server              *spvWalletServer
			expectedPaymailsLen int
		}{
			{
				name:                fmt.Sprintf("%s should set paymail domain after fetching shared config", spvWalletPG.leader.alias),
				server:              spvWalletPG,
				expectedPaymailsLen: 1,
			},
			{
				name:                fmt.Sprintf("%s should set paymail domain after fetching shared config", spvWalletSL.leader.alias),
				server:              spvWalletSL,
				expectedPaymailsLen: 1,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// given:
				leader := tc.server.leader

				// when:
				got, err := leader.client.SharedConfig(ctx)

				// then:
				assert.NoError(err, "Shared config wasn't successful retrieved by %s. Expect to get nil error", leader.paymail)

				if assert.NotNil(got.PaymailDomains, "Shared config should contain non-nil paymail domains slice") {
					actualLen := len(got.PaymailDomains)

					assert.Equal(tc.expectedPaymailsLen, actualLen, "Retrieved shared config  should have %s paymail domains. Got: %d paymail domains", tc.expectedPaymailsLen, actualLen)
					assert.NotEmpty(got.PaymailDomains[0], "Retrieved shared config should not be an empty string")

					tc.server.setPaymailDomains(got.PaymailDomains[0])

					logSuccessOp(t, err, "%s retrieved the shared configuration successfully. Leader set the paymail domains to admin, leader, user clients.", leader.paymail)
				}
			})
		}
	})

	t.Run("Step 3: The SPV Wallet admin clients attempt to add a user's xPub records within the same environment by making a request to their SPV Wallet API instance.", func(t *testing.T) {
		tests := []struct {
			name   string
			server *spvWalletServer
		}{
			{
				name:   fmt.Sprintf("%s should add xPub record %s for %s", spvWalletPG.admin.paymail, spvWalletPG.user.xPub, spvWalletPG.user.paymail),
				server: spvWalletPG,
			},
			{
				name:   fmt.Sprintf("%s should add xPub record %s for %s", spvWalletPG.admin.paymail, spvWalletSL.user.xPub, spvWalletSL.user.paymail),
				server: spvWalletSL,
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// given:
				admin := tc.server.admin
				user := tc.server.user

				// when:
				xPub, err := admin.client.CreateXPub(ctx, &commands.CreateUserXpub{XPub: user.xPub})

				// then:
				assert.NoError(err, "xPub record %s wasn't created successfully for %s by %s. Expect to get nil error", user.xPub, user.paymail, admin.paymail)
				assert.NotNil(xPub, "Expected to get non-nil xPub response after sending creation request by %s.", admin.paymail)

				logSuccessOp(t, err, "xPub record %s was created successfully for %s by %s", user.xPub, user.paymail, admin.paymail)
			})
		}
	})

	t.Run("Step 4: The SPV Wallet admin clients attempt to add a user's paymail record within the same environment by making a request to their SPV Wallet API instance.", func(t *testing.T) {
		tests := []struct {
			name   string
			server *spvWalletServer
		}{
			{
				name:   fmt.Sprintf("%s should add paymail record %s for the user %s", spvWalletPG.admin.paymail, spvWalletPG.user.paymail, spvWalletPG.user.alias),
				server: spvWalletPG,
			},
			{
				name:   fmt.Sprintf("%s should add paymail record %s for the user %s", spvWalletPG.admin.paymail, spvWalletSL.user.paymail, spvWalletSL.user.alias),
				server: spvWalletSL,
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// given:
				admin := tc.server.admin
				user := tc.server.user

				// when:
				paymail, err := admin.client.CreatePaymail(ctx, &commands.CreatePaymail{
					Key:        user.xPub,
					Address:    user.paymail,
					PublicName: "Regression tests",
				})

				// then:
				assert.NoError(err, "Paymail record %s wasn't created successfully for %s by %s. Expect to get nil error", user.paymail, user.alias, admin.paymail)
				assert.NotNil(paymail, "Expected to get non-nil paymail addresss response after sending creation request by %s.", admin.paymail)
				user.paymailID = paymail.ID

				logSuccessOp(t, err, "Paymail record %s was created successfully for %s by %s.", user.paymail, user.alias, admin.paymail)
			})
		}
	})

	t.Run("Step 5: The leader clients from one environment attempt to make internal transfers to users within their environment using the appropriate SPV Wallet API instance.", func(t *testing.T) {
		tests := []struct {
			name   string
			server *spvWalletServer
			funds  uint64
		}{
			{
				server: spvWalletPG,
				funds:  2,
				name:   fmt.Sprintf("%s should transfer 2 satoshis to the user %s", spvWalletPG.leader.paymail, spvWalletPG.user.paymail),
			},
			{
				server: spvWalletSL,
				funds:  3,
				name:   fmt.Sprintf("%s should transfer 3 satoshis to the user %s", spvWalletSL.leader.paymail, spvWalletSL.user.paymail),
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// given:
				user := tc.server.user
				leader := tc.server.leader

				userBalance, err := user.balance(ctx)
				assert.NoError(err, "Expected to get nil error after fetching balance by %s", user.paymail)

				leaderBalance, err := leader.balance(ctx)
				assert.NoError(err, "Expected to get nil error after fetching balance by %s", leader.paymail)

				// when:
				transaction, err := leader.transferFunds(ctx, user.paymail, tc.funds)

				// then:
				assert.NoError(err, "Transfer funds %d wasn't successful from %s to %s. Expect to get nil error", tc.funds, leader.paymail, user.paymail)

				if assert.NotNil(transaction, "Expected to get non-nil transaction response after transfer funds %d from %s to %s", tc.funds, leader.paymail, user.paymail) {
					// Verify sender's balance after the transaction
					senderBalanceCorrect := assertBalanceAfterTransaction(ctx, t, leader, leaderBalance-tc.funds-transaction.Fee)

					// Verify recipient's balance after the transaction
					recipientBalanceCorrect := assertBalanceAfterTransaction(ctx, t, user, userBalance+tc.funds)

					// Verify that the transaction appears in the recipient's transaction list
					page, err := user.client.Transactions(ctx)
					assert.NoError(err, "Failed to retrieve transactions for recipient %s. Expected to get nil error", user.paymail)

					recipientTransactionsCorrect := assert.True(transactionsSlice(page.Content).Has(transaction.ID), "Transaction %s made by %s was not found in %s's transaction list.", transaction.ID, leader.paymail, user.paymail)

					if senderBalanceCorrect && recipientBalanceCorrect && recipientTransactionsCorrect {
						logSuccessOp(t, nil, "Transfer funds %d was successful from leader %s to user %s", tc.funds, leader.paymail, user.paymail)
					}
				}
			})
		}
	})

	t.Run("Step 6: The user from one env attempts to transfer funds to the user from external env using the appropriate SPV Wallet API instance", func(t *testing.T) {
		tests := []struct {
			name      string
			sender    *user
			recipient *user
			funds     uint64
		}{
			{
				name:      fmt.Sprintf("%s should transfer 2 satoshis to %s", spvWalletSL.user.paymail, spvWalletPG.user.paymail),
				sender:    spvWalletSL.user,
				recipient: spvWalletPG.user,
				funds:     2,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// given:
				sender := tc.sender
				recipient := tc.recipient

				recipientBalance, err := recipient.balance(ctx)
				assert.NoError(err, "Expected to get nil error after fetching balance by %s", recipient.paymail)

				senderBalance, err := sender.balance(ctx)
				assert.NoError(err, "Expected to get nil error after fetching balance by %s", sender.paymail)

				// when:
				transaction, err := sender.transferFunds(ctx, recipient.paymail, tc.funds)

				// then:
				assert.NoError(err, "Transfer funds %d wasn't successful from sender %s to recipient %s. Expect to get nil error after making transaction, got error: %v", tc.funds, sender.paymail, recipient.paymail)

				if assert.NotNil(transaction, "Expected to get non-nil transaction response after transfer funds %d from sender %s to recipient %s", tc.funds, sender.paymail, recipient.paymail) {
					// Verify sender's balance after the transaction
					senderBalanceCorrect := assertBalanceAfterTransaction(ctx, t, sender, senderBalance-tc.funds-transaction.Fee)

					// Verify recipient's balance after the transaction
					recipientBalanceCorrect := assertBalanceAfterTransaction(ctx, t, tc.recipient, recipientBalance+tc.funds)

					if senderBalanceCorrect && recipientBalanceCorrect {
						logSuccessOp(t, nil, "Transfer funds %d was successful from sender %s to recipient %s", tc.funds, sender.paymail, recipient.paymail)
					}
				}
			})
		}
	})

	// The flow is designed to test the contacts management feature in the SPV Wallet API.
	// The verification flow is a series of steps that involve creating contacts, initiating verification, and confirming contacts.

	// given: Bob, Alice, Tom, Jerry are initialized as users for the contacts flow verification.
	bob, alice, tom, jerry := prepareUsersForContactsFlowVerification(t, spvWalletPG, spvWalletSL)

	t.Run("Step 7: Bob adds Alice as contact", func(t *testing.T) {
		// when: Bob attempts to add Alice as a contact
		contact, err := bob.addContact(ctx, alice.paymail, alice.alias)

		// then: The operation should succeed, and the contact should not be nil
		require.NoError(t, err, "Bob failed to add Alice as contact")
		require.NotNil(t, contact, "Bob failed to add Alice as contact")

		logSuccessOp(t, nil, "Contact successfully created between %s and %s", bob.alias, alice.alias)
	})

	t.Run("Step 8: Alice adds Bob as contact", func(t *testing.T) {
		// when: Alice attempts to add Bob as a contact
		contact, err := alice.addContact(ctx, bob.paymail, bob.alias)

		// then: The operation should succeed, and the contact should not be nil
		require.NoError(t, err, "Alice failed to add Bob as contact")
		require.NotNil(t, contact, "Alice failed to add Bob as contact")

		logSuccessOp(t, nil, "Contact successfully created between %s and %s", alice.alias, bob.alias)
	})

	t.Run("Step 9: Bob should confirm contact with Alice", func(t *testing.T) {
		// when: Alice generates a TOTP for Bob
		totp, err := alice.generateTotp(ctx, bob.paymail)
		require.NoError(t, err, "Failed to generate TOTP for Bob")
		require.NotEmpty(t, totp, "TOTP should not be empty")

		// when: Bob confirms Alice as a contact using the generated TOTP
		err = bob.confirmContact(ctx, alice.paymail, totp)

		// then: The contact status should be confirmed
		require.NoError(t, err, "Bob failed to confirm Alice")

		contact, err := bob.getContact(ctx, alice.paymail)
		require.NoError(t, err, "Failed to fetch contact")
		require.Equal(t, response.ContactConfirmed, contact.Status, "Alice's contact status should be confirmed")

		logSuccessOp(t, nil, "Contact successfully confirmed between %s and %s", bob.alias, alice.alias)
	})

	t.Run("Step 10: Bob should unconfirm contact with Alice", func(t *testing.T) {
		// when: Bob unconfirms Alice as a contact
		err := bob.unconfirmContact(ctx, alice.paymail)
		require.NoError(t, err, "Bob failed to unconfirm Alice")

		// then: The contact status should be unconfirmed
		contact, err := bob.getContact(ctx, alice.paymail)
		require.NoError(t, err, "Failed to fetch contact")
		require.Equal(t, response.ContactNotConfirmed, contact.Status, "Alice's contact status should be unconfirmed")

		logSuccessOp(t, nil, "Contact successfully unconfirmed between %s and %s", bob.alias, alice.alias)
	})

	t.Run("Step 11: Bob should remove Alice from contacts", func(t *testing.T) {
		// when: Bob removes Alice from his contacts
		err := bob.removeContact(ctx, alice.paymail)

		// then: The operation should succeed
		require.NoError(t, err, "Bob failed to remove Alice")

		logSuccessOp(t, nil, "Contacts successfully removed between %s and %s", bob.alias, alice.alias)
	})

	t.Run("Step 12: Tom adds Jerry as contact", func(t *testing.T) {
		// when: Tom attempts to add Jerry as a contact
		contact, err := tom.addContact(ctx, jerry.paymail, jerry.alias)

		// then: The operation should succeed, and the contact should not be nil
		require.NoError(t, err, "Tom failed to add Jerry as contact")
		require.NotNil(t, contact, "Tom failed to add Jerry as contact")

		logSuccessOp(t, nil, "Contacts successfully created between %s and %s", tom.alias, jerry.alias)
	})

	t.Run("Step 13: Jerry adds Tom as contact", func(t *testing.T) {
		// when: Jerry attempts to add Tom as a contact
		contact, err := jerry.addContact(ctx, tom.paymail, tom.alias)

		// then: The operation should succeed, and the contact should not be nil
		require.NoError(t, err, "Jerry failed to add Tom as contact")
		require.NotNil(t, contact, "Jerry failed to add Tom as contact")

		logSuccessOp(t, nil, "Contacts successfully created between %s and %s", jerry.alias, tom.alias)
	})

	t.Run("Step 14: Tom should confirm contact with Jerry", func(t *testing.T) {
		// when: Jerry generates a TOTP for Tom
		totp, err := jerry.generateTotp(ctx, tom.paymail)
		require.NoError(t, err, "Failed to generate TOTP for Jerry")
		require.NotEmpty(t, totp, "TOTP should not be empty")

		// when: Tom confirms Jerry as a contact using the generated TOTP
		err = tom.confirmContact(ctx, jerry.paymail, totp)

		// then: The contact status should be confirmed
		require.NoError(t, err, "Tom failed to confirm Jerry")

		contact, err := tom.getContact(ctx, jerry.paymail)
		require.NoError(t, err, "Failed to fetch contact")
		require.Equal(t, response.ContactConfirmed, contact.Status, "Jerry's contact status should be confirmed")

		logSuccessOp(t, nil, "Contact successfully confirmed between %s and %s", tom.alias, jerry.alias)
	})

	t.Run("Step 15: Tom should unconfirm contact with Jerry", func(t *testing.T) {
		// when: Tom unconfirms Jerry as a contact
		err := tom.unconfirmContact(ctx, jerry.paymail)
		require.NoError(t, err, "Tom failed to unconfirm Jerry")

		// then: The contact status should be unconfirmed
		contact, err := tom.getContact(ctx, jerry.paymail)
		require.NoError(t, err, "Failed to fetch contact")
		require.Equal(t, response.ContactNotConfirmed, contact.Status, "Jerry's contact status should be unconfirmed")

		logSuccessOp(t, nil, "Contact successfully unconfirmed between %s and %s", tom.alias, jerry.alias)
	})

	t.Run("Step 16: Tom should remove Jerry from contacts", func(t *testing.T) {
		// when: Tom removes Jerry from his contacts
		err := tom.removeContact(ctx, jerry.paymail)

		// then: The operation should succeed
		require.NoError(t, err, "Tom failed to remove Jerry")

		logSuccessOp(t, nil, "Contact successfully removed between %s and %s", tom.alias, jerry.alias)
	})

	t.Run("Step 17: The SPV Wallet users manage their access keys.", func(t *testing.T) {
		tests := []struct {
			name   string
			server *spvWalletServer
		}{
			{
				name:   fmt.Sprintf("%s should generate and retrieve an access key", spvWalletPG.user.paymail),
				server: spvWalletPG,
			},
			{
				name:   fmt.Sprintf("%s should generate and retrieve an access key", spvWalletSL.user.paymail),
				server: spvWalletSL,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// given:
				user := tc.server.user

				// when: Generate a new access key
				accessKey, err := user.generateAccessKey(ctx)
				require.NoError(t, err, "Failed to generate access key for %s", user.paymail)
				require.NotNil(t, accessKey, "Expected non-nil access key response")
				logSuccessOp(t, err, "User %s successfully generated an access key", user.paymail)

				// when: Fetch access key by ID
				retrievedKey, err := user.getAccessKeyByID(ctx, accessKey.ID)
				require.NoError(t, err, "Failed to retrieve access key %s for user %s", accessKey.ID, user.paymail)
				require.Equal(t, accessKey.ID, retrievedKey.ID, "Fetched access key ID should match generated key ID")
				logSuccessOp(t, err, "User %s successfully retrieved access key %s", user.paymail, accessKey.ID)

				keys, err := user.getAccessKeys(ctx)
				require.NoError(t, err, "Failed to retrieve access keys for user %s", user.paymail)
				require.NotNil(t, keys, "Expected non-nil access keys response")
				assert.True(len(keys) > 0, "Expected at least one access key to be present")
				logSuccessOp(t, err, "User %s successfully retrieved access keys", user.paymail)

				// when: Revoke access key
				err = user.client.RevokeAccessKey(ctx, accessKey.ID)
				require.NoError(t, err, "Failed to revoke access key %s for user %s", accessKey.ID, user.paymail)
				logSuccessOp(t, err, "User %s successfully revoked access key %s", user.paymail, accessKey.ID)
			})
		}
	})

	t.Run("Step 18: The SPV Wallet admin clients fetch all access keys.", func(t *testing.T) {
		tests := []struct {
			name   string
			server *spvWalletServer
		}{
			{
				name:   fmt.Sprintf("%s should fetch all access keys", spvWalletPG.admin.paymail),
				server: spvWalletPG,
			},
			{
				name:   fmt.Sprintf("%s should fetch all access keys", spvWalletSL.admin.paymail),
				server: spvWalletSL,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// given:
				admin := tc.server.admin
				user := tc.server.user
				accessKey, err := user.generateAccessKey(ctx)
				require.NoError(t, err, "Failed to generate access key for %s", user.paymail)
				require.NotNil(t, accessKey, "Expected non-nil access key response")

				// when:
				keys, err := admin.getAccessKeysAdmin(ctx)

				// then:
				assert.NoError(err, "Fetching access keys failed for %s", admin.paymail)
				assert.NotNil(keys, "Expected non-nil access keys response")
				assert.True(len(keys) > 0, "Expected at least one access key to be present")
				logSuccessOp(t, err, "Admin %s successfully retrieved access keys", admin.paymail)
			})
		}
	})

	// given: The SPV Wallet admin clients are prepared for contacts flow verification.
	adminSL, adminPG := prepareAdminForContactsFlowVerification(t, spvWalletSL, spvWalletPG)
	t.Run("Step 19: The SPV Wallet admin clients attempt to add user contacts.", func(t *testing.T) {
		tests := []struct {
			title string
			admin *admin
			user  *user
		}{
			{
				title: fmt.Sprintf("%s should add Bob as contact", adminPG.paymail),
				admin: adminPG,
				user:  bob,
			},
			{
				title: fmt.Sprintf("%s should add Tom as contact", adminSL.paymail),
				admin: adminSL,
				user:  tom,
			},
		}

		for _, tc := range tests {
			t.Run(tc.title, func(t *testing.T) {
				// Create Contact
				// when: Admin creates a contact for the user
				contact, err := tc.admin.createContact(ctx, tc.user.paymail, tc.user.alias)
				// then: The operation should succeed, and the contact should not be nil
				require.NoError(t, err, "Admin %s failed to create contact for %s", tc.admin.paymail, tc.user.paymail)
				require.NotNil(t, contact, "Expected non-nil contact response")

				logSuccessOp(t, err, "Contact %s successfully created by %s", contact.Paymail, tc.admin.paymail)

				// Update Contact
				// given: The contact ID
				contactID := contact.ID
				updatedName := fmt.Sprintf("%s Updated", tc.user.alias)
				// when: Admin updates the contact name
				updatedContact, err := tc.admin.updateContact(ctx, contactID, updatedName)
				// then: The operation should succeed, and the updated contact name should match
				require.NoError(t, err, "Failed to update contact for %s", tc.user.paymail)
				require.Equal(t, updatedContact.FullName, updatedName, "Updated contact name should match")

				logSuccessOp(t, err, "Contact %s successfully updated by %s", contactID, tc.admin.paymail)

				// Remove Contact
				// when: Admin removes the contact
				err = tc.admin.deleteContact(ctx, contactID)
				// then: The operation should succeed
				require.NoError(t, err, "Admin %s failed to remove contact %s", tc.admin.paymail, contactID)

				logSuccessOp(t, err, "Contact %s successfully removed by %s", contactID, tc.admin.paymail)
			})
		}
	})

	t.Run("Step 20: The SPV Wallet admin clients attempt to confirm contacts.", func(t *testing.T) {
		tests := []struct {
			title string
			admin *admin
			userA *user
			userB *user
		}{
			{
				title: "Admin should confirm contact between Alice and Bob",
				admin: adminPG,
				userA: bob,
				userB: alice,
			},
			{
				title: "Admin should confirm contact between Tom and Jerry",
				admin: adminSL,
				userA: tom,
				userB: jerry,
			},
		}

		for _, tc := range tests {
			t.Run(tc.title, func(t *testing.T) {
				// given: Create contacts
				contactA, errA := tc.userA.addContact(ctx, tc.userB.paymail, tc.userB.alias)
				contactB, errB := tc.userB.addContact(ctx, tc.userA.paymail, tc.userA.alias)
				require.NoError(t, errA, "Failed to create contact A")
				require.NoError(t, errB, "Failed to create contact B")

				// when: Admin confirms the contact between the users
				err := tc.admin.confirmContact(ctx, contactA.Paymail, contactB.Paymail)

				// then: The operation should succeed
				require.NoError(t, err, "Admin failed to confirm contact between %s and %s", contactA.Paymail, contactB.Paymail)

				logSuccessOp(t, err, "Contact between %s and %s confirmed successfully", contactA.Paymail, contactB.Paymail)
			})
		}
	})

	t.Run("Step 21: The admin clients attempt to remove created actor paymails using the appropriate SPV Wallet API instance.", func(t *testing.T) {
		tests := []struct {
			name   string
			server *spvWalletServer
		}{
			{
				name:   fmt.Sprintf("%s should delete %s paymail record", spvWalletPG.admin.paymail, spvWalletPG.user.paymail),
				server: spvWalletPG,
			},
			{
				name:   fmt.Sprintf("%s should delete %s paymail record", spvWalletSL.admin.paymail, spvWalletSL.user.paymail),
				server: spvWalletSL,
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// given:
				admin := tc.server.admin
				paymail := tc.server.user.paymail
				paymailID := tc.server.user.paymailID

				// when:
				err := admin.client.DeletePaymail(ctx, paymailID)

				// then:
				assert.NoError(err, "Delete paymail %s wasn't successful by %s. Expect to get nil error, got error: %v", paymail, admin.paymail, err)
				logSuccessOp(t, err, "Delete paymail %s was successful by %s", paymail, admin.paymail)
			})
		}
	})
}

func initServers(t *testing.T) (*spvWalletServer, *spvWalletServer) {
	const adminXPriv = "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"
	const (
		clientOneURL         = "CLIENT_ONE_URL"
		clientOneLeaderXPriv = "CLIENT_ONE_LEADER_XPRIV"
		clientTwoURL         = "CLIENT_TWO_URL"
		clientTwoLeaderXPriv = "CLIENT_TWO_LEADER_XPRIV"
	)
	const (
		alias1 = "UserSLRegressionTest"
		alias2 = "UserPGRegressionTest"
	)

	spvWalletSL, err := initSPVWalletServer(alias1, &spvWalletServerConfig{
		envURL:     lookupEnvOrDefault(t, clientOneURL, ""),
		envXPriv:   lookupEnvOrDefault(t, clientOneLeaderXPriv, ""),
		adminXPriv: adminXPriv,
	})
	require.NoError(t, err, "Step 1: Setup failed could not initialize the clients for env: %s", spvWalletSL.cfg.envURL)

	spvWalletPG, err := initSPVWalletServer(alias2, &spvWalletServerConfig{
		envURL:     lookupEnvOrDefault(t, clientTwoURL, ""),
		envXPriv:   lookupEnvOrDefault(t, clientTwoLeaderXPriv, ""),
		adminXPriv: adminXPriv,
	})

	require.NoError(t, err, "Step 1: Setup failed could not initialize the clients for env: %s", spvWalletPG.cfg.envURL)

	return spvWalletPG, spvWalletSL
}
