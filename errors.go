package walletclient

import (
	"encoding/json"
	"errors"
	"github.com/bitcoin-sv/spv-wallet/models"
	"net/http"
)

// ErrAdminKey admin key not set
var ErrAdminKey = errors.New("an admin key must be set to be able to create an xpub")

// ErrNoClientSet is when no client is set
var ErrNoClientSet = errors.New("no transport client set")

// WrapError wraps an error into SPVError
func WrapError(err error) *models.SPVError {
	if err == nil {
		return nil
	}

	return &models.SPVError{
		StatusCode: http.StatusInternalServerError,
		Message:    err.Error(),
		Code:       models.UnknownErrorCode,
	}
}

// WrapResponseError wraps a http response into SPVError
func WrapResponseError(res *http.Response) *models.SPVError {
	if res == nil {
		return nil
	}

	var resError *models.ResponseError

	err := json.NewDecoder(res.Body).Decode(&resError)
	if err != nil {
		return WrapError(err)
	}

	return &models.SPVError{
		StatusCode: res.StatusCode,
		Code:       resError.Code,
		Message:    resError.Message,
	}
}

func CreateErrorResponse(code string, message string) *models.SPVError {
	return &models.SPVError{
		StatusCode: http.StatusInternalServerError,
		Code:       code,
		Message:    message,
	}
}
