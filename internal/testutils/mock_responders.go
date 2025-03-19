package testutils

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

// NewJSONFileResponderWithStatusOK returns a new SPVError object with status code 400
func NewJSONFileResponderWithStatusOK(filePath string) httpmock.Responder {
	return httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File(filePath))
}

// NewJSONBodyResponderWithStatusOK returns a new SPVError object with status code 400
func NewJSONBodyResponderWithStatusOK(body any) httpmock.Responder {
	return httpmock.NewJsonResponderOrPanic(http.StatusOK, body)
}

// NewStringResponderStatusOK returns a new responder with status code 200 and a body
func NewStringResponderStatusOK(body string) httpmock.Responder {
	return httpmock.NewStringResponder(http.StatusOK, body)
}

// NewBadRequestSPVErrorResponder returns a new SPVError object with status code 400
func NewBadRequestSPVErrorResponder() httpmock.Responder {
	return httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, NewBadRequestSPVError())
}

// NewResourceNotFoundSPVErrorResponder returns a new SPVError object with status code 404.
func NewResourceNotFoundSPVErrorResponder() httpmock.Responder {
	return httpmock.NewJsonResponderOrPanic(http.StatusNotFound, NewResourceNotFoundSPVError())
}

// NewConflictRequestSPVErrorResponder returns a new SPVError object with status code 409.
func NewConflictRequestSPVErrorResponder() httpmock.Responder {
	return httpmock.NewJsonResponderOrPanic(http.StatusConflict, NewConflictRequestSPVError())
}

// NewInternalServerSPVErrorResponder returns a new SPVError object with status code 500
func NewInternalServerSPVErrorResponder() httpmock.Responder {
	return httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, NewInternalServerSPVError())
}

// NewInternalServerSPVErrorStringResponder returns a new SPVError object with status code 500
func NewInternalServerSPVErrorStringResponder(errMessage string) httpmock.Responder {
	return httpmock.NewStringResponder(http.StatusInternalServerError, errMessage)
}

// NewUnauthorizedAccessSPVErrorResponder returns a new SPVError object with status code 401
func NewUnauthorizedAccessSPVErrorResponder() httpmock.Responder {
	return httpmock.NewJsonResponderOrPanic(http.StatusUnauthorized, NewUnauthorizedAccessSPVError())
}

// RegisterMockResponder registers a mock responder for a given endpoint
func RegisterMockResponder(t *testing.T, client *resty.Client, endpoint string, statusCode int, responseBody interface{}) {
	url := client.BaseURL + endpoint
	httpmock.RegisterResponder("GET", url, httpmock.NewJsonResponderOrPanic(statusCode, responseBody))
	t.Cleanup(func() { httpmock.Reset() })
}

// NewPaginatedJSONResponder creates a responder that simulates paginated responses
func NewPaginatedJSONResponder(t *testing.T, files ...string) httpmock.Responder {
	responses := make([]*http.Response, 0, len(files))
	for _, file := range files {
		responses = append(responses, httpmock.NewStringResponse(http.StatusOK, httpmock.File(file).String())) //nolint: bodyclose
	}
	return httpmock.ResponderFromMultipleResponses(responses)
}
