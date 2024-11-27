package errutil

import (
	"fmt"
	"net/http"
)

type HTTPErrorFormatter struct {
	Action string
	API    string
	Err    error
}

func (h HTTPErrorFormatter) Format(method string) error {
	return fmt.Errorf("failed to send HTTP %s request to %s via %s: %w", method, h.Action, h.API, h.Err)
}

func (h HTTPErrorFormatter) FormatPutErr() error    { return h.Format(http.MethodPut) }
func (h HTTPErrorFormatter) FormatPatchErr() error  { return h.Format(http.MethodPatch) }
func (h HTTPErrorFormatter) FormatPostErr() error   { return h.Format(http.MethodPost) }
func (h HTTPErrorFormatter) FormatGetErr() error    { return h.Format(http.MethodGet) }
func (h HTTPErrorFormatter) FormatDeleteErr() error { return h.Format(http.MethodDelete) }
