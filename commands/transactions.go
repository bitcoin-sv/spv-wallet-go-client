package commands

import (
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/queryparams"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// RecordTransaction holds the arguments required to record a user transaction.
type RecordTransaction struct {
	Metadata    queryparams.Metadata `json:"metadata"`    // Metadata associated with the transaction.
	Hex         string               `json:"hex"`         // Hexadecimal string representation of the transaction.
	ReferenceID string               `json:"referenceId"` // Reference ID for the transaction.
}

// DraftTransaction holds the arguments required to create user draft transaction.
type DraftTransaction struct {
	Config   response.TransactionConfig `json:"config"`   // Configuration for the transaction.
	Metadata queryparams.Metadata       `json:"metadata"` // Metadata related to the transaction.
}

// UpdateTransactionMetadata holds the arguments required to update the metadata of a user transaction.
// The ID field is ignored in the request body sent to the SPV Wallet API; instead, it is used as part
// of the transaction metadata update endpoint (e.g., /api/v1/transactions/{ID}).
type UpdateTransactionMetadata struct {
	ID       string               `json:"-"`        // Unique identifier of the transaction to be updated.
	Metadata queryparams.Metadata `json:"metadata"` // New metadata to associate with the transaction.
}

// Recipients represents a single recipient in a transaction.
// It includes details about the recipient address, the amount to send,
// and an optional OP_RETURN script for including additional data in the transaction.
type Recipients struct {
	OpReturn *response.OpReturn `json:"op_return"` // Optional OP_RETURN script for attaching data to the transaction.
	Satoshis uint64             `json:"satoshis"`  // Amount to send to the recipient, in satoshis.
	To       string             `json:"to"`        // Paymails address of the recipient.
}

// SendToRecipients holds the arguments required to send a transaction to multiple recipients.
// This includes the list of recipients with their details and optional metadata for the transaction.
type SendToRecipients struct {
	Recipients []*Recipients        `json:"recipients"` // List of recipients for the transaction.
	Metadata   queryparams.Metadata `json:"metadata"`   // Metadata associated with the transaction.
}
