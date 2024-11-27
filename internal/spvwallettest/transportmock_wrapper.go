package spvwallettest

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/jarcoal/httpmock"
)

type TransportWrapper struct {
	*httpmock.MockTransport

	mu           sync.RWMutex
	lastRequest  *http.Request
	lastResponse *http.Response
	lastError    error
}

// NewTransportWrapper creates a new wrapper around the default httpmock.MockTransport.
func NewTransportWrapper() *TransportWrapper {
	return &TransportWrapper{
		MockTransport: httpmock.NewMockTransport(),
	}
}

// RoundTrip intercepts the request and stores the response and error.
func (tw *TransportWrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := tw.MockTransport.RoundTrip(req)

	tw.mu.Lock()
	defer tw.mu.Unlock()

	tw.lastRequest = req
	tw.lastResponse = resp
	tw.lastError = err

	if err != nil {
		return resp, fmt.Errorf("Round trip error - %w", err)
	}
	return resp, nil
}

// GetResponse retrieves the last response and error.
func (tw *TransportWrapper) GetResponse() (*http.Response, error) {
	tw.mu.RLock()
	defer tw.mu.RUnlock()

	return tw.lastResponse, tw.lastError
}

// GetRequest retrieves the last request.
func (tw *TransportWrapper) GetRequest() *http.Request {
	tw.mu.RLock()
	defer tw.mu.RUnlock()

	return tw.lastRequest
}
