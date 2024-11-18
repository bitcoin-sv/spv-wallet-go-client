package commands

// UpsertContact holds the necessary arguments for adding or updating a user's contact information.
type UpsertContact struct {
	FullName string         `json:"fullName"`         // The full name of the user.
	Metadata map[string]any `json:"metadata"`         // Metadata associated with the transaction.
	Paymail  string         `json:"requesterPaymail"` // Paymail address of the user, which is used for secure and simplified payment transfers.
}
