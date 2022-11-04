package transports

import "github.com/BuxOrg/bux"

// TransportType the type of transport being used (http or graphql)
type TransportType string

// BuxUserAgent the bux user agent sent to the bux server
const BuxUserAgent = "BUX: go-client " + BuxClientVersion

// BuxClientVersion is the version of the client
const BuxClientVersion = "v0.2.4"

const (
	// BuxTransportHTTP uses the http transport for all bux server actions
	BuxTransportHTTP TransportType = "http"

	// BuxTransportGraphQL uses the graphql transport for all bux server actions
	BuxTransportGraphQL TransportType = "graphql"

	// BuxTransportMock uses the mock transport for all bux server actions
	BuxTransportMock TransportType = "mock"
)

// Recipients is a struct for recipients
type Recipients struct {
	OpReturn *bux.OpReturn `json:"op_return"`
	Satoshis uint64        `json:"satoshis"`
	To       string        `json:"to"`
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
)
