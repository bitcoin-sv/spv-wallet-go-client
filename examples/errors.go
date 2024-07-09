package examples

import (
	"errors"
	"fmt"
	"github.com/bitcoin-sv/spv-wallet/models"
)

func GetFullErrorMessage(err error) {
	var errMsg string

	var spvError models.SPVError
	if errors.As(err, &spvError) {
		errMsg = fmt.Sprintf("Error, Message: %s, Code: %s, HTTP status code: %d", spvError.GetMessage(), spvError.GetCode(), spvError.GetStatusCode())
	} else {
		errMsg = fmt.Sprintf("Error, Message: %s, Code: %s, HTTP status code: %d", err.Error(), models.UnknownErrorCode, 500)
	}
	fmt.Println(errMsg)
}
