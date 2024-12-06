package commands

import "github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"

// UpdateXPubMetadata contains the parameters needed to update the metadata
// associated with the current user's xpub.
type UpdateXPubMetadata struct {
	Metadata querybuilders.Metadata `json:"metadata"` // Key-value pairs representing the xpub metadata
}

// GenerateAccessKey contains the parameters needed to generate a new access key
// for the current user, including any associated metadata.
type GenerateAccessKey struct {
	Metadata querybuilders.Metadata `json:"metadata"` // Key-value pairs representing the access key metadata
}
