//go:build regression
// +build regression

package regressiontests

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	fundsPerTest                  = 2
	adminXPriv                    = "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"
	adminXPub                     = "xpub661MyMwAqRbcFgfmdkPgE2m5UjHXu9dj124DbaGLSjaqVESTWfCD4VuNmEbVPkbYLCkykwVZvmA8Pbf8884TQr1FgdG2nPoHR8aB36YdDQh"
	errorWhileGettingTransaction  = "error while getting transaction: %s"
	errorWhileCreatingUser        = "error while creating user: %s"
	errorWhileGettingBalance      = "error while getting transaction: %s"
	errorWhileSendingFunds        = "error while sending funds: %s"
	errorWhileGettingSharedConfig = "error while getting shared config: %s"
	errorWhileGettingEnvVariables = "error while getting env variables: %s"
)

var (
	clientOneURL         string
	clientTwoURL         string
	clientOneLeaderXPriv string
	clientTwoLeaderXPriv string
)

func TestRegression(t *testing.T) {
	ctx := context.Background()
	rtConfig, err := getEnvVariables()
	require.NoError(t, err, fmt.Sprintf(errorWhileGettingEnvVariables, err))

	sharedConfigInstanceOne, err := getSharedConfig(adminXPub, rtConfig.ClientOneURL)
	require.NoError(t, err, fmt.Sprintf(errorWhileGettingSharedConfig, err))

	sharedConfigInstanceTwo, err := getSharedConfig(adminXPub, rtConfig.ClientTwoURL)
	require.NoError(t, err, fmt.Sprintf(errorWhileGettingSharedConfig, err))

	userName := "instanceOneUser1"
	userOne, err := createUser(ctx, userName, sharedConfigInstanceOne.PaymailDomains[0], rtConfig.ClientOneURL, adminXPriv)
	require.NoError(t, err, fmt.Sprintf(errorWhileCreatingUser, err))

	defer deleteUser(ctx, userOne.Paymail, rtConfig.ClientOneURL, adminXPriv)

	userName = "instanceTwoUser1"
	userTwo, err := createUser(ctx, userName, sharedConfigInstanceTwo.PaymailDomains[0], rtConfig.ClientTwoURL, adminXPriv)
	require.NoError(t, err, fmt.Sprintf(errorWhileCreatingUser, err))

	defer deleteUser(ctx, userTwo.Paymail, rtConfig.ClientTwoURL, adminXPriv)

	t.Run("TestInitialBalancesAndTransactionsBeforeAndAfterFundTransfers", func(t *testing.T) {
		// given
		balance, err := getBalance(ctx, rtConfig.ClientOneURL, userOne.XPriv)
		require.NoError(t, err, fmt.Sprintf(errorWhileGettingBalance, err))

		transactions, err := getTransactions(ctx, rtConfig.ClientOneURL, userOne.XPriv)
		require.NoError(t, err, fmt.Sprintf(errorWhileGettingTransaction, err))

		assert.Equal(t, 0, balance)
		assert.Equal(t, 0, len(transactions))

		balance, err = getBalance(ctx, rtConfig.ClientTwoURL, userTwo.XPriv)
		require.NoError(t, err, fmt.Sprintf(errorWhileGettingBalance, err))

		transactions, err = getTransactions(ctx, rtConfig.ClientTwoURL, userTwo.XPriv)
		require.NoError(t, err, fmt.Sprintf(errorWhileGettingTransaction, err))

		assert.Equal(t, 0, balance)
		assert.Equal(t, 0, len(transactions))

		// when
		transactionOne, err := sendFunds(ctx, rtConfig.ClientOneURL, rtConfig.ClientOneLeaderXPriv, userTwo.Paymail, 2)
		require.NoError(t, err, fmt.Sprintf(errorWhileSendingFunds, err))
		assert.GreaterOrEqual(t, int64(-1), transactionOne.OutputValue)

		transactionTwo, err := sendFunds(ctx, rtConfig.ClientTwoURL, rtConfig.ClientTwoLeaderXPriv, userOne.Paymail, 2)
		require.NoError(t, err, fmt.Sprintf(errorWhileSendingFunds, err))
		assert.GreaterOrEqual(t, int64(-1), transactionTwo.OutputValue)

		// then
		balance, err = getBalance(ctx, rtConfig.ClientOneURL, userOne.XPriv)
		require.NoError(t, err, fmt.Sprintf(errorWhileGettingBalance, err))

		transactions, err = getTransactions(ctx, rtConfig.ClientOneURL, userOne.XPriv)
		require.NoError(t, err, fmt.Sprintf(errorWhileGettingTransaction, err))

		assert.GreaterOrEqual(t, balance, 1)
		assert.GreaterOrEqual(t, len(transactions), 1)

		balance, err = getBalance(ctx, rtConfig.ClientTwoURL, userTwo.XPriv)
		require.NoError(t, err, fmt.Sprintf(errorWhileGettingBalance, err))

		transactions, err = getTransactions(ctx, rtConfig.ClientTwoURL, userTwo.XPriv)
		require.NoError(t, err, fmt.Sprintf(errorWhileGettingTransaction, err))
		assert.GreaterOrEqual(t, balance, 1)
		assert.GreaterOrEqual(t, len(transactions), 1)
	})
}
