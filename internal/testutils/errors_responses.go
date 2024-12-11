package testutils

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// NewBadRequestSPVError creates a new SPVError for bad request
func NewBadRequestSPVError() models.SPVError {
	return models.SPVError{
		Message:    http.StatusText(http.StatusBadRequest),
		StatusCode: http.StatusBadRequest,
		Code:       models.UnknownErrorCode,
	}
}

// NewUnauthorizedAccessSPVError creates a new SPVError for unauthorized access
func NewUnauthorizedAccessSPVError() models.SPVError {
	return models.SPVError{
		Message:    "unauthorized",
		StatusCode: http.StatusUnauthorized,
		Code:       "error-unauthorized",
	}
}

// NewInternalServerSPVError creates a new SPVError for internal server error
func NewInternalServerSPVError() models.SPVError {
	return models.SPVError{
		Message:    http.StatusText(http.StatusInternalServerError),
		StatusCode: http.StatusInternalServerError,
		Code:       models.UnknownErrorCode,
	}
}

// NewUnrecognizedAPIResponseError creates a new SPVError for unrecognized API response
func NewUnrecognizedAPIResponseError() models.SPVError {
	return models.SPVError{
		Message:    errors.ErrUnrecognizedAPIResponse.Error(),
		StatusCode: http.StatusInternalServerError,
		Code:       "internal-server-error",
	}
}

// NewInvalidRequestError creates a new SPVError for invalid request
func NewInvalidRequestError() models.SPVError {
	return models.SPVError{
		Message:    "Invalid request",
		StatusCode: http.StatusBadRequest,
		Code:       "invalid-request",
	}
}
