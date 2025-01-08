package errutil_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/stretchr/testify/require"
)

func TestHTTPErrorFormatter_Format(t *testing.T) {
	// given:
	const (
		API    = "Users API"
		action = "retrieve users page"
	)
	wrappedErr := errors.New(http.StatusText(http.StatusInternalServerError))
	expectedErr := fmt.Errorf("failed to send HTTP %s request to %s via %s: %w", http.MethodPost, action, API, wrappedErr)

	formatter := errutil.HTTPErrorFormatter{
		Action: action,
		API:    API,
		Err:    wrappedErr,
	}

	// when:
	got := formatter.Format(http.MethodPost)

	// then:
	require.Equal(t, got, expectedErr)
}
