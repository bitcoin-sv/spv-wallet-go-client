//go:build regression
// +build regression

package regressiontests

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	fundsPerTest = 2

	adminXPriv = "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"
	adminXPub  = "xpub661MyMwAqRbcFgfmdkPgE2m5UjHXu9dj124DbaGLSjaqVESTWfCD4VuNmEbVPkbYLCkykwVZvmA8Pbf8884TQr1FgdG2nPoHR8aB36YdDQh"

	errGettingEnvVariables = "failed to get environment variables: %s"
	errGettingSharedConfig = "failed to get shared config: %s"
	errCreatingUser        = "failed to create user: %s"
	errDeletingUserPaymail = "failed to delete user's paymail: %s"
	errSendingFunds        = "failed to send funds: %s"
	errGettingBalance      = "failed to get balance: %s"
	errGettingTransactions = "failed to get transactions: %s"
)

func TestRegression(t *testing.T) {
	ctx := context.Background()
	rtConfig, err := getEnvVariables()
	require.NoError(t, err, fmt.Sprintf(errGettingEnvVariables, err))

	var paymailDomainInstanceOne, paymailDomainInstanceTwo string
	var userOne, userTwo *regressionTestUser

	t.Run("Initialize Shared Configurations", func(t *testing.T) {
		t.Run("Should get sharedConfig for instance one", func(t *testing.T) {
			paymailDomainInstanceOne, err = getPaymailDomain(adminXPub, rtConfig.ClientOneURL)
			require.NoError(t, err, fmt.Sprintf(errGettingSharedConfig, err))
		})

		t.Run("Should get shared config for instance two", func(t *testing.T) {
			paymailDomainInstanceTwo, err = getPaymailDomain(adminXPub, rtConfig.ClientTwoURL)
			require.NoError(t, err, fmt.Sprintf(errGettingSharedConfig, err))
		})
	})

	t.Run("Create Users", func(t *testing.T) {
		t.Run("Should create user for instance one", func(t *testing.T) {
			userName := "instanceOneUser1"
			userOne, err = createUser(ctx, userName, paymailDomainInstanceOne, rtConfig.ClientOneURL, adminXPriv)
			require.NoError(t, err, fmt.Sprintf(errCreatingUser, err))
		})

		t.Run("Should create user for instance two", func(t *testing.T) {
			userName := "instanceTwoUser1"
			userTwo, err = createUser(ctx, userName, paymailDomainInstanceTwo, rtConfig.ClientTwoURL, adminXPriv)
			require.NoError(t, err, fmt.Sprintf(errCreatingUser, err))
		})
	})

	defer func() {
		t.Run("Cleanup: Remove Paymails", func(t *testing.T) {
			t.Run("Should remove user's paymail on first instance", func(t *testing.T) {
				if userOne != nil {
					err := removeRegisteredPaymail(ctx, userOne.Paymail, rtConfig.ClientOneURL, adminXPriv)
					require.NoError(t, err, fmt.Sprintf(errDeletingUserPaymail, err))
				}
			})

			t.Run("Should remove user's paymail on second instance", func(t *testing.T) {
				if userTwo != nil {
					err := removeRegisteredPaymail(ctx, userTwo.Paymail, rtConfig.ClientTwoURL, adminXPriv)
					require.NoError(t, err, fmt.Sprintf(errDeletingUserPaymail, err))
				}
			})
		})
	}()

	t.Run("Perform Transactions", func(t *testing.T) {
		t.Run("Send money to instance 1", func(t *testing.T) {
			transaction, err := sendFunds(ctx, rtConfig.ClientTwoURL, rtConfig.ClientTwoLeaderXPriv, userOne.Paymail, fundsPerTest)
			require.NoError(t, err, fmt.Sprintf(errSendingFunds, err))
			require.GreaterOrEqual(t, int64(-1), transaction.OutputValue)

			balance, err := getBalance(ctx, rtConfig.ClientOneURL, userOne.XPriv)
			require.NoError(t, err, fmt.Sprintf(errGettingBalance, err))
			require.GreaterOrEqual(t, balance, 1)

			transactions, err := getTransactions(ctx, rtConfig.ClientOneURL, userOne.XPriv)
			require.NoError(t, err, fmt.Sprintf(errGettingTransactions, err))
			require.GreaterOrEqual(t, len(transactions), 1)
		})

		t.Run("Send money to instance 2", func(t *testing.T) {
			transaction, err := sendFunds(ctx, rtConfig.ClientOneURL, rtConfig.ClientOneLeaderXPriv, userTwo.Paymail, fundsPerTest)
			require.NoError(t, err, fmt.Sprintf(errSendingFunds, err))
			require.GreaterOrEqual(t, int64(-1), transaction.OutputValue)

			balance, err := getBalance(ctx, rtConfig.ClientTwoURL, userTwo.XPriv)
			require.NoError(t, err, fmt.Sprintf(errGettingBalance, err))
			require.GreaterOrEqual(t, balance, 1)

			transactions, err := getTransactions(ctx, rtConfig.ClientTwoURL, userTwo.XPriv)
			require.NoError(t, err, fmt.Sprintf(errGettingTransactions, err))
			require.GreaterOrEqual(t, len(transactions), 1)
		})

		t.Run("Send money from instance 1 to instance 2", func(t *testing.T) {
			transaction, err := sendFunds(ctx, rtConfig.ClientOneURL, userOne.XPriv, userTwo.Paymail, fundsPerTest)
			require.NoError(t, err, fmt.Sprintf(errSendingFunds, err))
			require.GreaterOrEqual(t, int64(-1), transaction.OutputValue)

			balance, err := getBalance(ctx, rtConfig.ClientTwoURL, userTwo.XPriv)
			require.NoError(t, err, fmt.Sprintf(errGettingBalance, err))
			require.GreaterOrEqual(t, balance, 1)

			transactions, err := getTransactions(ctx, rtConfig.ClientTwoURL, userTwo.XPriv)
			require.NoError(t, err, fmt.Sprintf(errGettingTransactions, err))
			require.GreaterOrEqual(t, len(transactions), 1)

			balance, err = getBalance(ctx, rtConfig.ClientOneURL, userOne.XPriv)
			require.NoError(t, err, fmt.Sprintf(errGettingBalance, err))
			require.GreaterOrEqual(t, balance, 1)

			transactions, err = getTransactions(ctx, rtConfig.ClientOneURL, userOne.XPriv)
			require.NoError(t, err, fmt.Sprintf(errGettingTransactions, err))
			require.GreaterOrEqual(t, len(transactions), 1)
		})
	})
}
