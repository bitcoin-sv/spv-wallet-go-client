package client

import "github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/contacts"

type UpsertContactArgs struct {
	FullName string
	Metadata map[string]any
	Paymail  string
}

func (u UpsertContactArgs) parseUpsertContactRequest() contacts.UpsertContactRequest {
	return contacts.UpsertContactRequest{
		FullName:         u.FullName,
		Metadata:         u.Metadata,
		RequesterPaymail: u.Paymail,
	}
}
