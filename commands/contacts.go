package commands

// UpsertContact holds the necessary arguments for adding or updating a user's contact information.
type UpsertContact struct {
	ContactPaymail   string         `json:"-"`                // Paymail address of the new/updating contact.
	FullName         string         `json:"fullName"`         // The full name of the user.
	Metadata         map[string]any `json:"metadata"`         // Metadata associated with the contact.
	RequesterPaymail string         `json:"requesterPaymail"` // RequesterPaymail address of the user, which is used for secure and simplified payment transfers.
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

// ConfirmContacts represents the body defined for confirming contact's between two users.
type ConfirmContacts struct {
	PaymailA string `json:"paymailA"` // The paymail address of the first user in the contact relationship
	PaymailB string `json:"paymailB"` // The paymail address of the second user in the contact relationship
}

// CreateContact holds the necessary arguments for creating a contact in the SPV Wallet system.
// It includes the paymail of the creator (the user initiating the contact addition),
// the full name of the contact, and any associated metadata.
// Note: The `Paymail` field is not included in the request body sent to the SPV Wallet API.
// Instead, it is used as part of the endpoint path (e.g., /api/v1/admin/contacts/{paymail}).
type CreateContact struct {
	// CreatorPaymail is the paymail address of the user who is adding the contact.
	// It identifies the owner or creator of the contact.
	CreatorPaymail string `json:"creatorPaymail"`

	Paymail string `json:"-"` // Paymail identifier of the contact to be created.

	// FullName is the full name of the contact to be added.
	// This is the name that will be associated with the contact in the system.
	FullName string `json:"fullName"`

	// Metadata is additional information that can be associated with the contact.
	// It is a key-value pair where the value can be of any type.
	Metadata map[string]any `json:"metadata"`
}
