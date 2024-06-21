package walletclient

import (
	"encoding/json"
	"errors"
	"github.com/bitcoin-sv/spv-wallet/spverrors"
	"net/http"
)

// ErrAdminKey admin key not set
var ErrAdminKey = errors.New("an admin key must be set to be able to create an xpub")

// ErrNoClientSet is when no client is set
var ErrNoClientSet = errors.New("no transport client set")

// WrapError wraps an error into SPVError
func WrapError(err error) *spverrors.SPVError {
	if err == nil {
		return nil
	}

	return &spverrors.SPVError{
		StatusCode: http.StatusInternalServerError,
		Message:    err.Error(),
		Code:       spverrors.UnknownErrorCode,
	}
}

// WrapResponseError wraps a http response into SPVError
func WrapResponseError(res *http.Response) *spverrors.SPVError {
	if res == nil {
		return nil
	}

	var resError *spverrors.ResponseError

	err := json.NewDecoder(res.Body).Decode(&resError)
	if err != nil {
		return WrapError(err)
	}

	return &spverrors.SPVError{
		StatusCode: res.StatusCode,
		Code:       resError.Code,
		Message:    resError.Message,
	}
}
