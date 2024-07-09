package walletclient

import (
	"encoding/json"
	"github.com/bitcoin-sv/spv-wallet/models"
	"net/http"
)

// ErrAdminKey admin key not set
var ErrAdminKey = models.SPVError{Message: "an admin key must be set to be able to create an xpub", StatusCode: 401, Code: "error-unauthorized-admin-key-not-set"}

// ErrMissingXpriv is when xpriv is missing
var ErrMissingXpriv = models.SPVError{Message: "xpriv missing", StatusCode: 401, Code: "error-unauthorized-xpriv-missing"}

// ErrMissingXprivAndXpub is when xpriv and xpub are missing
var ErrMissingXprivAndXpub = models.SPVError{Message: "xpriv and xpub missing", StatusCode: 404, Code: "error-shared-config-xpriv-and-xpub-missing"}

// ErrMissingAccessKey is when access key is missing
var ErrMissingAccessKey = models.SPVError{Message: "access key missing", StatusCode: 401, Code: "error-unauthorized-access-key-missing"}

// ErrCouldNotFindDraftTransaction is when draft transaction is not found
var ErrCouldNotFindDraftTransaction = models.SPVError{Message: "could not find draft transaction", StatusCode: 404, Code: "error-draft-transaction-not-found"}

// ErrTotpInvalid is when totp is invalid
var ErrTotpInvalid = models.SPVError{Message: "totp invalid", StatusCode: 400, Code: "error-totp-invalid"}

// ErrContactPubKeyInvalid is when contact's PubKey is invalid
var ErrContactPubKeyInvalid = models.SPVError{Message: "contact's PubKey is invalid", StatusCode: 400, Code: "error-contact-pubkey-invalid"}

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
