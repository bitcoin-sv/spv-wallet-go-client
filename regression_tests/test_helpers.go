package regressiontests

import (
	"context"
	"os"
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
