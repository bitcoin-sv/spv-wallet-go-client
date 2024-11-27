package commands

import "github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"

// CreateUserXpub contains the parameters required to register a user's XPub.
type CreateUserXpub struct {
	Metadata querybuilders.Metadata `json:"metadata"` // Metadata associated with the XPub.
	XPub     string                 `json:"key"`      // The user's XPub key to be recorded.
}
