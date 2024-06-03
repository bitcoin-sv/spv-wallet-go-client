package walletclient

import "github.com/bitcoin-sv/spv-wallet/models"

// TransportType the type of transport being used ('http' for usage or 'mock' for testing)
type TransportType string

// SPVWalletUserAgent the spv wallet user agent sent to the spv wallet.
const SPVWalletUserAgent = "SPVWallet: go-client"

const (
	// SPVWalletTransportHTTP uses the http transport for all spv-wallet actions
	SPVWalletTransportHTTP TransportType = "http"

	// SPVWalletTransportMock uses the mock transport for all spv-wallet actions
	SPVWalletTransportMock TransportType = "mock"
)

// Recipients is a struct for recipients
type Recipients struct {
	OpReturn *models.OpReturn `json:"op_return"`
	Satoshis uint64           `json:"satoshis"`
	To       string           `json:"to"`
}

const (
	// FieldMetadata is the field name for metadata
	FieldMetadata = "metadata"

	// FieldQueryParams is the field name for the query params
	FieldQueryParams = "params"

	// FieldXpubKey is the field name for xpub key
	FieldXpubKey = "key"

	// FieldXpubID is the field name for xpub id
	FieldXpubID = "xpub_id"

	// FieldAddress is the field name for paymail address
	FieldAddress = "address"

	// FieldPublicName is the field name for (paymail) public name
	FieldPublicName = "public_name"

	// FieldAvatar is the field name for (paymail) avatar
	FieldAvatar = "avatar"

	// FieldConditions is the field name for conditions
	FieldConditions = "conditions"

	// FieldTo is the field name for "to"
	FieldTo = "to"

	// FieldSatoshis is the field name for "satoshis"
	FieldSatoshis = "satoshis"

	// FieldOpReturn is the field name for "op_return"
	FieldOpReturn = "op_return"

	// FieldConfig is the field name for "config"
	FieldConfig = "config"

	// FieldOutputs is the field name for "outputs"
	FieldOutputs = "outputs"

	// FieldHex is the field name for "hex"
	FieldHex = "hex"

	// FieldReferenceID is the field name for "reference_id"
	FieldReferenceID = "reference_id"

	// FieldID is the id field for most models
	FieldID = "id"

	// FieldLockingScript is the field for locking script
	FieldLockingScript = "locking_script"

	// FieldUserAgent is the field for storing the user agent
	FieldUserAgent = "user_agent"

	// FieldTransactionConfig is the field for the config of a new transaction
	FieldTransactionConfig = "transaction_config"

	// FieldTransactionID is the field for transaction ID
	FieldTransactionID = "tx_id"

	// FieldOutputIndex is the field for "output_index"
	FieldOutputIndex = "output_index"
)
