package regressiontests

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// assertBalanceAfterTransaction checks if a user's balance matches the expected value after a transaction.
// This function is intended for use in tests and relies on the `assert` library for assertions.
// It calls t.Helper() to ensure that errors are reported at the caller's location.
func assertBalanceAfterTransaction(ctx context.Context, t *testing.T, user *user, expectedBalance uint64) bool {
	t.Helper()

	actualBalance, err := user.balance(ctx)
	if !assert.NoError(t, err, "Failed to retrieve balance for %s", user.paymail) {
		return false
	}

	return assert.Equal(t, expectedBalance, actualBalance, "Balance mismatch for %s.", user.paymail)
}

// lookupEnvOrDefault retrieves the value of the specified environment variable.
// If the variable is not set, it returns the provided default value.
// This function is intended for use in tests. It calls t.Helper()
// to ensure that errors are reported at the caller's location.
func lookupEnvOrDefault(t *testing.T, env string, defaultValue string) string {
	t.Helper()

	v, ok := os.LookupEnv(env)
	if !ok {
		t.Logf("Environment variable %s not set, using default: %s", env, defaultValue)
		return defaultValue
	}
	return v
}

// logSuccessOp logs a success message if there is no error.
// This function is intended for use in tests. It calls t.Helper()
// to ensure that errors are reported at the caller's location.
func logSuccessOp(t *testing.T, err error, format string, args ...any) {
	t.Helper()

	if err != nil {
		return
	}
	t.Logf(format, args...)
}

// prepareUsersForContactsFlowVerification prepares users for the contacts flow verification.
// It initializes the users and returns them as user objects.
// This function is intended for use in tests. It calls t.Helper()
// to ensure that errors are reported at the caller's location.
func prepareUsersForContactsFlowVerification(t *testing.T, spvWalletPG, spvWalletSL *spvWalletServer) (bob *user, alice *user, tom *user, jerry *user) {
	t.Helper()

	// PG = Postgres
	bob = spvWalletPG.leader
	alice = spvWalletPG.user
	// SL = SQLite
	tom = spvWalletSL.leader
	jerry = spvWalletSL.user

	alice = &user{
		alias:     "Alice",
		paymail:   strings.ToLower(alice.paymail),
		xPub:      alice.xPub,
		xPriv:     alice.xPriv,
		paymailID: alice.paymailID,
		client:    alice.client,
	}
	bob = &user{
		alias:     "Bob",
		paymail:   strings.ToLower(bob.paymail),
		xPub:      bob.xPub,
		xPriv:     bob.xPriv,
		paymailID: bob.paymailID,
		client:    bob.client,
	}
	tom = &user{
		alias:     "Tom",
		paymail:   strings.ToLower(tom.paymail),
		xPub:      tom.xPub,
		xPriv:     tom.xPriv,
		paymailID: tom.paymailID,
		client:    tom.client,
	}
	jerry = &user{
		alias:     "Jerry",
		paymail:   strings.ToLower(jerry.paymail),
		xPub:      jerry.xPub,
		xPriv:     jerry.xPriv,
		paymailID: jerry.paymailID,
		client:    jerry.client,
	}
	return
}

// prepareAdminForContactsFlowVerification prepares the admin for the contacts flow verification.
// It initializes the admin and returns it as an admin object.
// This function is intended for use in tests. It calls t.Helper()
// to ensure that errors are reported at the caller's location.
func prepareAdminForContactsFlowVerification(t *testing.T, spvWalletPG, spvWalletSL *spvWalletServer) (adminSL *admin, adminPG *admin) {
	t.Helper()

	// PG = Postgres
	adminPG = spvWalletPG.admin
	// SL = SQLite
	adminSL = spvWalletSL.admin

	adminPG = &admin{
		alias:   "AdminPG",
		paymail: strings.ToLower(adminPG.paymail),
		xPriv:   adminPG.xPriv,
		client:  adminPG.client,
	}
	adminSL = &admin{
		alias:   "AdminSL",
		paymail: strings.ToLower(adminSL.paymail),
		xPriv:   adminSL.xPriv,
		client:  adminSL.client,
	}
	return
}
