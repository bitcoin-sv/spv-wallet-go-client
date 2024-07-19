package regressiontests

import (
	"testing"
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
	userOne, err := createUser(userName, sharedConfigInstanceOne.PaymailDomains[0], rtConfig.ClientOneURL, adminXPriv)
	if err != nil {
		t.Errorf(errorWhileCreatingUser, err)
	}

	userName = "instanceTwoUser1"
	userTwo, err := createUser(userName, sharedConfigInstanceTwo.PaymailDomains[0], rtConfig.ClientTwoURL, adminXPriv)
	if err != nil {
		t.Errorf(errorWhileCreatingUser, err)
	}

	t.Run("TestInitialBalancesAndTransactionsBeforeAndAfterFundTransfers", func(t *testing.T) {})
}
