package walletclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ErrAdminKey admin key not set
var ErrAdminKey = errors.New("an admin key must be set to be able to create an xpub")

// ErrNoClientSet is when no client is set
var ErrNoClientSet = errors.New("no transport client set")

// ResError is a struct which contain information about error
type ResError struct {
	StatusCode int
	Message    string
}

// ResponseError is an interface for error
type ResponseError interface {
	Error() string
	GetStatusCode() int
}

// WrapError wraps an error into ResponseError
func WrapError(err error) ResponseError {
	if err == nil {
		return nil
	}

	return &ResError{
		StatusCode: http.StatusInternalServerError,
		Message:    err.Error(),
	}
}

// WrapResponseError wraps a http response into ResponseError
func WrapResponseError(res *http.Response) ResponseError {
	if res == nil {
		return nil
	}

	var errorMsg string

	err := json.NewDecoder(res.Body).Decode(&errorMsg)
	if err != nil {
		// if EOF, then body is empty and we return response status as error message
		if !errors.Is(err, io.EOF) {
			errorMsg = fmt.Sprintf("spv-wallet error message can't be decoded. Reason: %s", err.Error())
		}
		errorMsg = res.Status
	}

	return &ResError{
		StatusCode: res.StatusCode,
		Message:    errorMsg,
	}
}

// Error returns the error message
func (e *ResError) Error() string {
	return e.Message
}

// GetStatusCode returns the status code of error
func (e *ResError) GetStatusCode() int {
	return e.StatusCode
}
