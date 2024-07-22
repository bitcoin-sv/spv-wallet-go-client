//go:build regression
// +build regression

package regressiontests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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
	if err != nil {
		t.Errorf(errorWhileGettingEnvVariables, err)
	}

	sharedConfigInstanceOne, err := getSharedConfig(adminXPub, rtConfig.ClientOneURL)
	if err != nil {
		t.Errorf(errorWhileGettingSharedConfig, err)
	}
	sharedConfigInstanceTwo, err := getSharedConfig(adminXPub, rtConfig.ClientTwoURL)
	if err != nil {
		t.Errorf(errorWhileGettingSharedConfig, err)
	}

	userName := "instanceOneUser1"
	userOne, err := createUser(ctx, userName, sharedConfigInstanceOne.PaymailDomains[0], rtConfig.ClientOneURL, adminXPriv)
	if err != nil {
		t.Errorf(errorWhileCreatingUser, err)
	}
	defer deleteUser(ctx, userOne.Paymail, rtConfig.ClientOneURL, adminXPriv)

	userName = "instanceTwoUser1"
	userTwo, err := createUser(ctx, userName, sharedConfigInstanceTwo.PaymailDomains[0], rtConfig.ClientTwoURL, adminXPriv)
	if err != nil {
		t.Errorf(errorWhileCreatingUser, err)
	}
	defer deleteUser(ctx, userTwo.Paymail, rtConfig.ClientTwoURL, adminXPriv)

	t.Run("TestInitialBalancesAndTransactionsBeforeAndAfterFundTransfers", func(t *testing.T) {
		// Given
		balance, err := getBalance(ctx, rtConfig.ClientOneURL, userOne.XPriv)
		if err != nil {
			t.Errorf(errorWhileGettingBalance, err)
		}
		transactions, err := getTransactions(ctx, rtConfig.ClientOneURL, userOne.XPriv)
		if err != nil {
			t.Errorf(errorWhileGettingTransaction, err)
		}
		assert.Equal(t, 0, balance)
		assert.Equal(t, 0, len(transactions))

		balance, err = getBalance(ctx, rtConfig.ClientTwoURL, userTwo.XPriv)
		if err != nil {
			t.Errorf(errorWhileGettingBalance, err)
		}
		transactions, err = getTransactions(ctx, rtConfig.ClientTwoURL, userTwo.XPriv)
		if err != nil {
			t.Errorf(errorWhileGettingTransaction, err)
		}
		assert.Equal(t, 0, balance)
		assert.Equal(t, 0, len(transactions))

		// When
		transactionOne, err := sendFunds(ctx, rtConfig.ClientOneURL, rtConfig.ClientOneLeaderXPriv, userTwo.Paymail, 2)
		if err != nil {
			t.Errorf(errorWhileSendingFunds, err)
		}
		assert.GreaterOrEqual(t, int64(-1), transactionOne.OutputValue)

		transactionTwo, err := sendFunds(ctx, rtConfig.ClientTwoURL, rtConfig.ClientTwoLeaderXPriv, userOne.Paymail, 2)
		if err != nil {
			t.Errorf(errorWhileSendingFunds, err)
		}
		assert.GreaterOrEqual(t, int64(-1), transactionTwo.OutputValue)

		// Then
		balance, err = getBalance(ctx, rtConfig.ClientOneURL, userOne.XPriv)
		if err != nil {
			t.Errorf(errorWhileGettingBalance, err)
		}
		transactions, err = getTransactions(ctx, rtConfig.ClientOneURL, userOne.XPriv)
		if err != nil {
			t.Errorf(errorWhileGettingTransaction, err)
		}
		assert.GreaterOrEqual(t, balance, 1)
		assert.GreaterOrEqual(t, len(transactions), 1)

		balance, err = getBalance(ctx, rtConfig.ClientTwoURL, userTwo.XPriv)
		if err != nil {
			t.Errorf(errorWhileGettingBalance, err)
		}
		transactions, err = getTransactions(ctx, rtConfig.ClientTwoURL, userTwo.XPriv)
		if err != nil {
			t.Errorf(errorWhileGettingTransaction, err)
		}
		assert.GreaterOrEqual(t, balance, 1)
		assert.GreaterOrEqual(t, len(transactions), 1)
	})
}
