package configstest

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/models"
)

func NewBadRequestSPVError() models.SPVError {
	return models.SPVError{
		Message:    http.StatusText(http.StatusBadRequest),
		StatusCode: http.StatusBadRequest,
		Code:       "invalid-data-format",
	}
}
