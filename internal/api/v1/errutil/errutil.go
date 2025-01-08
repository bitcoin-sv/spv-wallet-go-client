package errutil

import (
	"fmt"
	"net/http"
)

// HTTPErrorFormatter is a utility struct that formats HTTP errors.
type HTTPErrorFormatter struct {
	Action string
	API    string
	Err    error
}

// Format creates a formatted error message for HTTP requests.
func (h HTTPErrorFormatter) Format(method string) error {
	return fmt.Errorf("failed to send HTTP %s request to %s via %s: %w", method, h.Action, h.API, h.Err)
}

// FormatPutErr is a convenience method for formatting HTTP PUT errors.
func (h HTTPErrorFormatter) FormatPutErr() error { return h.Format(http.MethodPut) }

// FormatPatchErr is a convenience method for formatting HTTP PATCH errors.
func (h HTTPErrorFormatter) FormatPatchErr() error { return h.Format(http.MethodPatch) }

// FormatPostErr is a convenience method for formatting HTTP POST errors.
func (h HTTPErrorFormatter) FormatPostErr() error { return h.Format(http.MethodPost) }

// FormatGetErr is a convenience method for formatting HTTP GET errors.
func (h HTTPErrorFormatter) FormatGetErr() error { return h.Format(http.MethodGet) }

// FormatDeleteErr is a convenience method for formatting HTTP DELETE errors.
func (h HTTPErrorFormatter) FormatDeleteErr() error { return h.Format(http.MethodDelete) }

// NewHTTPErrorFormatter creates a new instance of HTTPErrorFormatter.
// It eliminates redundancy and ensures consistency across the codebase.
func NewHTTPErrorFormatter(api string, action string, err error) *HTTPErrorFormatter {
	return &HTTPErrorFormatter{
		API:    api,
		Action: action,
		Err:    err,
	}
}
