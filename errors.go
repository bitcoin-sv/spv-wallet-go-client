package walletclient

import (
	"encoding/json"
	"github.com/bitcoin-sv/spv-wallet/models"
	"net/http"
)

// ErrAdminKey admin key not set
var ErrAdminKey = models.SPVError{Message: "an admin key must be set to be able to create an xpub", StatusCode: 401, Code: "error-unauthorized-admin-key-not-set"}

// ErrMissingXpriv is when xpriv is missing
var ErrMissingXpriv = models.SPVError{Message: "xpriv is missing", StatusCode: 401, Code: "error-unauthorized-xpriv-missing"}

// ErrInvalidXpriv is when xpriv is invalid
var ErrInvalidXpriv = models.SPVError{Message: "xpriv is invalid", StatusCode: 401, Code: "error-unauthorized-xpriv-invalid"}

// ErrInvalidXpub is when xpub is invalid
var ErrInvalidXpub = models.SPVError{Message: "xpub is invalid", StatusCode: 401, Code: "error-unauthorized-xpub-invalid"}

// ErrInvalidAccessKey is when access key is invalid
var ErrInvalidAccessKey = models.SPVError{Message: "access key is invalid", StatusCode: 401, Code: "error-unauthorized-access-key-invalid"}

// ErrInvalidAdminKey is when admin key is invalid
var ErrInvalidAdminKey = models.SPVError{Message: "admin key is invalid", StatusCode: 401, Code: "error-unauthorized-admin-key-invalid"}

// ErrInvalidServerURL is when server url is invalid
var ErrInvalidServerURL = models.SPVError{Message: "server url is invalid", StatusCode: 401, Code: "error-unauthorized-server-url-invalid"}

// ErrCreateClient is when client creation fails
var ErrCreateClient = models.SPVError{Message: "failed to create client", StatusCode: 500, Code: "error-create-client-failed"}

// ErrMissingKey is when neither xPriv nor adminXPriv is provided
var ErrMissingKey = models.SPVError{Message: "neither xPriv nor adminXPriv is provided", StatusCode: 404, Code: "error-shared-config-key-missing"}

// ErrMissingAccessKey is when access key is missing
var ErrMissingAccessKey = models.SPVError{Message: "access key is missing", StatusCode: 401, Code: "error-unauthorized-access-key-missing"}

// ErrCouldNotFindDraftTransaction is when draft transaction is not found
var ErrCouldNotFindDraftTransaction = models.SPVError{Message: "could not find draft transaction", StatusCode: 404, Code: "error-draft-transaction-not-found"}

// ErrTotpInvalid is when totp is invalid
var ErrTotpInvalid = models.SPVError{Message: "totp is invalid", StatusCode: 400, Code: "error-totp-invalid"}

// ErrContactPubKeyInvalid is when contact's PubKey is invalid
var ErrContactPubKeyInvalid = models.SPVError{Message: "contact's PubKey is invalid", StatusCode: 400, Code: "error-contact-pubkey-invalid"}

// ErrStaleLastEvaluatedKey is when the last evaluated key returned from sync merkleroots is the same as it was in a previous iteration
// indicating sync issue or a potential loop
var ErrStaleLastEvaluatedKey = models.SPVError{Message: "The last evaluated key has not changed between requests, indicating a possible loop or synchronization issue.", StatusCode: 500, Code: "error-stale-last-evaluated-key"}

// ErrStaleLastEvaluatedKey is when the last evaluated key returned from sync merkleroots is the same as it was in a previous iteration
// indicating sync issue or a potential loop
var ErrSyncMerkleRootsTimeout = models.SPVError{Message: "SyncMerkleRoots operation timed out", StatusCode: 500, Code: "error-sync-merkleroots-timeout"}

// WrapError wraps an error into SPVError
func WrapError(err error) error {
	if err == nil {
		return nil
	}

	return models.SPVError{
		StatusCode: http.StatusInternalServerError,
		Message:    err.Error(),
		Code:       models.UnknownErrorCode,
	}
}

// WrapResponseError wraps a http response into SPVError
func WrapResponseError(res *http.Response) error {
	if res == nil {
		return nil
	}

	var resError models.ResponseError

	err := json.NewDecoder(res.Body).Decode(&resError)
	if err != nil {
		return WrapError(err)
	}

	return models.SPVError{
		StatusCode: res.StatusCode,
		Code:       resError.Code,
		Message:    resError.Message,
	}
}

func CreateErrorResponse(code string, message string) error {
	return models.SPVError{
		StatusCode: http.StatusInternalServerError,
		Code:       code,
		Message:    message,
	}
}
