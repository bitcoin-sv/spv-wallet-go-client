package commands

// UpdateXPubMetadata contains the parameters needed to update the metadata
// associated with the current user's xpub.
type UpdateXPubMetadata struct {
	Metadata map[string]any `json:"metadata"` // Key-value pairs representing the xpub metadata
}

// GenerateAccessKey contains the parameters needed to generate a new access key
// for the current user, including any associated metadata.
type GenerateAccessKey struct {
	Metadata map[string]any `json:"metadata"` // Key-value pairs representing the access key metadata
}
