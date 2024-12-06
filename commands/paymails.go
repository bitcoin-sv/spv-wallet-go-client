package commands

import "github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"

// CreatePaymail defines the parameters required to create a new paymail address,
// including associated metadata such as the public name and avatar.
type CreatePaymail struct {
	Metadata   querybuilders.Metadata `json:"metadata"`    // Metadata associated with the paymail as key-value pairs.
	Key        string                 `json:"key"`         // The xpub key linked to the paymail.
	Address    string                 `json:"address"`     // The paymail address to be created.
	PublicName string                 `json:"public_name"` // The public display name associated with the paymail.
	Avatar     string                 `json:"avatar"`      // The URL of the paymail's avatar image.
}
