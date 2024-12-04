package commands

// UpsertContact holds the necessary arguments for adding or updating a user's contact information.
type UpsertContact struct {
	FullName string         `json:"fullName"`         // The full name of the user.
	Metadata map[string]any `json:"metadata"`         // Metadata associated with the contact.
	Paymail  string         `json:"requesterPaymail"` // Paymail address of the user, which is used for secure and simplified payment transfers.
}

// UpdateContact represents the arguments defined for updating a user's contact information.
//
// Note: The `ID` field is not included in the request body sent to the SPV Wallet API.
// Instead, it is used as part of the endpoint path (e.g., /api/v1/admin/contacts/{ID}).
type UpdateContact struct {
	ID       string         `json:"-"`        // Unique identifier of the contact to be updated.
	FullName string         `json:"fullName"` // The full name of the contact.
	Metadata map[string]any `json:"metadata"` // Metadata associated with the contact.
}
