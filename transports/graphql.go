package transports

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
	"github.com/machinebox/graphql"
)

// graphQlService is the interface for GraphQL
type graphQlService interface {
	Run(ctx context.Context, req *graphql.Request, resp interface{}) error
}

// TransportGraphQL is the graphql struct
type TransportGraphQL struct {
	accessKey   *bec.PrivateKey
	adminXPriv  *bip32.ExtendedKey
	client      graphQlService
	debug       bool
	httpClient  *http.Client
	server      string
	signRequest bool
	xPriv       *bip32.ExtendedKey
	xPub        *bip32.ExtendedKey
}

// XPubData is the xpub data
type XPubData struct {
	XPub *bux.Xpub `json:"xpub"`
}

// XPubMetadataData is the xpub data for the metadata mutation
type XPubMetadataData struct {
	XPub *bux.Xpub `json:"xpub_metadata"`
}

// AccessKeyData is the access key data
type AccessKeyData struct {
	AccessKey *bux.AccessKey `json:"access_key"`
}

// AccessKeysData is a slice of access key data
type AccessKeysData struct {
	AccessKeys []*bux.AccessKey `json:"access_keys"`
}

// DestinationData is the destination data
type DestinationData struct {
	Destination *bux.Destination `json:"destination"`
}

// DestinationMetadataData is the destination data for the metadata mutation
type DestinationMetadataData struct {
	Destination *bux.Destination `json:"destination_metadata"`
}

// DestinationsData is a slice of destination data
type DestinationsData struct {
	Destinations []*bux.Destination `json:"destinations"`
}

// DraftTransactionData is a draft transaction
type DraftTransactionData struct {
	NewTransaction *bux.DraftTransaction `json:"new_transaction"`
}

// TransactionData is a transaction
type TransactionData struct {
	Transaction *bux.Transaction `json:"transaction"`
}

// TransactionMetadataData is a transaction for the metadata mutation
type TransactionMetadataData struct {
	Transaction *bux.Transaction `json:"transaction_metadata"`
}

// TransactionsData is a slice of transactions
type TransactionsData struct {
	Transactions []*bux.Transaction `json:"transactions"`
}

// NewTransactionData is a transaction
type NewTransactionData struct {
	Transaction *bux.Transaction `json:"transaction"`
}

// Init will initialize
func (g *TransportGraphQL) Init() error {
	g.client = graphql.NewClient(g.server, graphql.WithHTTPClient(g.httpClient))
	return nil
}

// SetAdminKey set the admin key
func (g *TransportGraphQL) SetAdminKey(adminKey *bip32.ExtendedKey) {
	g.adminXPriv = adminKey
}

// SetDebug turn the debugging on or off
func (g *TransportGraphQL) SetDebug(debug bool) {
	g.debug = debug
}

// IsDebug return the debugging status
func (g *TransportGraphQL) IsDebug() bool {
	return g.debug
}

// SetSignRequest turn the signing of the HTTP request on or off
func (g *TransportGraphQL) SetSignRequest(signRequest bool) {
	g.signRequest = signRequest
}

// IsSignRequest return whether to sign all requests
func (g *TransportGraphQL) IsSignRequest() bool {
	return g.signRequest
}

// NewPaymail will register a new paymail
func (g *TransportGraphQL) NewPaymail(ctx context.Context, rawXpub, paymailAddress, avatar, publicName string, metadata *bux.Metadata) error {
	// TODO: Implement this
	return nil
}

// GetXpub will get an xPub
func (g *TransportGraphQL) GetXpub(ctx context.Context, rawXpub string) (*bux.Xpub, error) {
	// TODO: Implement this
	return nil, nil
}

// GetXPub will get information about the current xPub
func (g *TransportGraphQL) GetXPub(ctx context.Context) (*bux.Xpub, error) {

	reqBody := `
	query {
	  xpub {
		id
		current_balance
		next_internal_num
		next_external_num
		metadata
		created_at
		updated_at
		deleted_at
	  }
	}`

	var respData XPubData
	if err := g.doGraphQLQuery(ctx, reqBody, nil, &respData); err != nil {
		return nil, err
	}

	return respData.XPub, nil
}

// UpdateXPubMetadata update the metadata of the logged in xpub
func (g *TransportGraphQL) UpdateXPubMetadata(ctx context.Context, metadata *bux.Metadata) (*bux.Xpub, error) {

	reqBody := `
    mutation ($metadata: Metadata!) {
  	  xpub_metadata (
  	    metadata: $metadata
  	  ) {
		id
		current_balance
		next_internal_num
		next_external_num
		metadata
		created_at
		updated_at
		deleted_at
	  }
	}`
	variables := map[string]interface{}{
		FieldMetadata: processMetadata(metadata),
	}

	var respData XPubMetadataData
	if err := g.doGraphQLQuery(ctx, reqBody, variables, &respData); err != nil {
		return nil, err
	}

	return respData.XPub, nil
}

// GetAccessKey will get an access key by id
func (g *TransportGraphQL) GetAccessKey(ctx context.Context, id string) (*bux.AccessKey, error) {

	reqBody := `
	query ($id: String) {
      access_key (
        id: $id
      ) {
        id
        xpub_id
        key
        metadata
        created_at
        updated_at
        deleted_at
        revoked_at
      }
    }`

	var respData AccessKeyData
	if err := g.doGraphQLQuery(ctx, reqBody, map[string]interface{}{
		FieldID: id,
	}, &respData); err != nil {
		return nil, err
	}

	return respData.AccessKey, nil
}

// GetAccessKeys will get all access keys filtered by the metadata
func (g *TransportGraphQL) GetAccessKeys(ctx context.Context, metadata *bux.Metadata) ([]*bux.AccessKey, error) {

	reqBody := `
	query ($metadata: Metadata) {
      access_keys (
        metadata: $metadata
      ) {
        id
        xpub_id
        key
        metadata
        created_at
        updated_at
        deleted_at
        revoked_at
      }
    }`
	variables := map[string]interface{}{
		FieldMetadata: processMetadata(metadata),
	}

	var respData AccessKeysData
	if err := g.doGraphQLQuery(ctx, reqBody, variables, &respData); err != nil {
		return nil, err
	}

	return respData.AccessKeys, nil
}

// CreateAccessKey will create new access key
func (g *TransportGraphQL) CreateAccessKey(ctx context.Context, metadata *bux.Metadata) (*bux.AccessKey, error) {

	reqBody := `
	  mutation ($metadata: Metadata) {
        access_key (
          metadata: $metadata
        ) {
          id
          xpub_id
          key
          metadata
          created_at
          updated_at
          deleted_at
          revoked_at
        }
      }`
	variables := map[string]interface{}{
		FieldMetadata: processMetadata(metadata),
	}

	var respData AccessKeyData
	if err := g.doGraphQLQuery(ctx, reqBody, variables, &respData); err != nil {
		return nil, err
	}

	return respData.AccessKey, nil
}

// RevokeAccessKey will revoke the given access key
func (g *TransportGraphQL) RevokeAccessKey(ctx context.Context, id string) (*bux.AccessKey, error) {

	reqBody := `
	  mutation ($id: String) {
        access_key_revoke (
          id: $id
        ) {
          id
          xpub_id
          key
          metadata
          created_at
          updated_at
          deleted_at
          revoked_at
        }
      }`
	variables := map[string]interface{}{
		FieldID: id,
	}

	var respData AccessKeyData
	if err := g.doGraphQLQuery(ctx, reqBody, variables, &respData); err != nil {
		return nil, err
	}

	return respData.AccessKey, nil
}

// GetDestinationByID will get a destination by the given id
func (g *TransportGraphQL) GetDestinationByID(ctx context.Context, id string) (*bux.Destination, error) {

	reqBody := `{
	query ($id: String) {
        destination (
          id: $id
        ) {
          id
          xpub_id
          locking_script
          type
          chain
          num
          address
          metadata
          created_at
          updated_at
          deleted_at
        }
      }`
	variables := map[string]interface{}{
		FieldID: id,
	}

	var respData DestinationData
	if err := g.doGraphQLQuery(ctx, reqBody, variables, &respData); err != nil {
		return nil, err
	}

	return respData.Destination, nil
}

// GetDestinationByLockingScript will get a destination by the given locking script
func (g *TransportGraphQL) GetDestinationByLockingScript(ctx context.Context, lockingScript string) (*bux.Destination, error) {

	reqBody := `{
	query ($lockingScript: String) {
        destination (
          locking_script: $lockingScript
        ) {
          id
          xpub_id
          locking_script
          type
          chain
          num
          address
          metadata
          created_at
          updated_at
          deleted_at
        }
      }`
	variables := map[string]interface{}{
		FieldLockingScript: lockingScript,
	}

	var respData DestinationData
	if err := g.doGraphQLQuery(ctx, reqBody, variables, &respData); err != nil {
		return nil, err
	}

	return respData.Destination, nil
}

// GetDestinationByAddress will get a destination by the given address
func (g *TransportGraphQL) GetDestinationByAddress(ctx context.Context, address string) (*bux.Destination, error) {

	reqBody := `{
	query ($address: String) {
        destination (
          address: $address
        ) {
          id
          xpub_id
          locking_script
          type
          chain
          num
          address
          metadata
          created_at
          updated_at
          deleted_at
        }
      }`
	variables := map[string]interface{}{
		FieldAddress: address,
	}

	var respData DestinationData
	if err := g.doGraphQLQuery(ctx, reqBody, variables, &respData); err != nil {
		return nil, err
	}

	return respData.Destination, nil
}

// UpdateDestinationMetadataByID updates the destination metadata by id
func (g *TransportGraphQL) UpdateDestinationMetadataByID(ctx context.Context, id string, metadata *bux.Metadata) (*bux.Destination, error) {

	reqBody := `{
      mutation ($id: String, $metadata: Metadata!) {
  	    destination_metadata (
		  id: $id
  	      metadata: $metadata
  	    ) {
          id
          xpub_id
          locking_script
          type
          chain
          num
          address
          metadata
          created_at
          updated_at
          deleted_at
        }
      }`
	variables := map[string]interface{}{
		FieldID:       id,
		FieldMetadata: processMetadata(metadata),
	}

	var respData DestinationMetadataData
	if err := g.doGraphQLQuery(ctx, reqBody, variables, &respData); err != nil {
		return nil, err
	}

	return respData.Destination, nil
}

// UpdateDestinationMetadataByAddress updates the destination metadata by address
func (g *TransportGraphQL) UpdateDestinationMetadataByAddress(ctx context.Context, address string, metadata *bux.Metadata) (*bux.Destination, error) {

	reqBody := `{
      mutation ($address: String, $metadata: Metadata!) {
  	    destination_metadata (
		  address: $address
  	      metadata: $metadata
  	    ) {
          id
          xpub_id
          locking_script
          type
          chain
          num
          address
          metadata
          created_at
          updated_at
          deleted_at
        }
      }`
	variables := map[string]interface{}{
		FieldAddress:  address,
		FieldMetadata: processMetadata(metadata),
	}

	var respData DestinationMetadataData
	if err := g.doGraphQLQuery(ctx, reqBody, variables, &respData); err != nil {
		return nil, err
	}

	return respData.Destination, nil
}

// UpdateDestinationMetadataByLockingScript updates the destination metadata by lockingScript
func (g *TransportGraphQL) UpdateDestinationMetadataByLockingScript(ctx context.Context, lockingScript string, metadata *bux.Metadata) (*bux.Destination, error) {

	reqBody := `{
      mutation ($locking_script: String, $metadata: Metadata!) {
  	    destination_metadata (
		  locking_script: $locking_script
  	      metadata: $metadata
  	    ) {
          id
          xpub_id
          locking_script
          type
          chain
          num
          address
          metadata
          created_at
          updated_at
          deleted_at
        }
      }`
	variables := map[string]interface{}{
		FieldLockingScript: lockingScript,
		FieldMetadata:      processMetadata(metadata),
	}

	var respData DestinationMetadataData
	if err := g.doGraphQLQuery(ctx, reqBody, variables, &respData); err != nil {
		return nil, err
	}

	return respData.Destination, nil
}

// GetDestinations will get all destinations filtered by the medata conditions
func (g *TransportGraphQL) GetDestinations(ctx context.Context, metadataConditions *bux.Metadata) ([]*bux.Destination, error) {

	reqBody := `{
	  query ($metadata: Metadata) {
        destinations (
          metadata: $metadata
        ) {
          id
          xpub_id
          locking_script
          type
          chain
          num
          address
          metadata
          created_at
          updated_at
          deleted_at
        }
      }`
	variables := map[string]interface{}{
		FieldMetadata: processMetadata(metadataConditions),
	}

	var respData DestinationsData
	if err := g.doGraphQLQuery(ctx, reqBody, variables, &respData); err != nil {
		return nil, err
	}

	return respData.Destinations, nil
}

// NewDestination will get a new destination
func (g *TransportGraphQL) NewDestination(ctx context.Context, metadata *bux.Metadata) (*bux.Destination, error) {

	reqBody := `
   	mutation ($metadata: Metadata) {
	  destination(
		metadata: $metadata
	  ) {
		id
		xpub_id
		locking_script
		type
		chain
		num
		address
		metadata
	  }
	}`
	variables := map[string]interface{}{
		FieldMetadata: processMetadata(metadata),
	}

	var respData DestinationData
	if err := g.doGraphQLQuery(ctx, reqBody, variables, &respData); err != nil {
		return nil, err
	}

	return respData.Destination, nil
}

// GetTransaction get a transaction by ID
func (g *TransportGraphQL) GetTransaction(ctx context.Context, txID string) (*bux.Transaction, error) {

	reqBody := `
   	query {
	  transaction(
		id:"` + txID + `",
	  ) {
        id
        hex
        block_hash
        block_height
        fee
        number_of_inputs
        number_of_outputs
        output_value
        total_value
        direction
        metadata
        created_at
        updated_at
        deleted_at
	  }
	}`
	var respData TransactionData
	if err := g.doGraphQLQuery(ctx, reqBody, nil, &respData); err != nil {
		return nil, err
	}

	return respData.Transaction, nil
}

// GetTransactions get a transactions, filtered by the given metadata
func (g *TransportGraphQL) GetTransactions(ctx context.Context, conditions map[string]interface{},
	metadataConditions *bux.Metadata) ([]*bux.Transaction, error) {

	querySignature := ""
	queryArguments := ""

	// is there a better way to do this ?
	if conditions != nil {
		querySignature += "( $conditions Map "
		queryArguments += " conditions: $conditions\n"
	}
	if metadataConditions != nil {
		if conditions == nil {
			querySignature += "( "
		} else {
			querySignature += ", "
		}
		querySignature += "$metadata Map"
		queryArguments += " metadata: $metadata\n"
	} else {
		querySignature += " )"
	}

	reqBody := `
   	query ` + querySignature + `{
	  transactions ` + queryArguments + ` {
        id
        hex
        block_hash
        block_height
        fee
        number_of_inputs
        number_of_outputs
        output_value
        total_value
        direction
        metadata
        created_at
        updated_at
        deleted_at
	  }
	}`
	variables := make(map[string]interface{})
	if conditions != nil {
		variables[FieldConditions] = conditions
	}
	if metadataConditions != nil {
		variables[FieldMetadata] = metadataConditions
	}

	var respData TransactionsData
	if err := g.doGraphQLQuery(ctx, reqBody, variables, &respData); err != nil {
		return nil, err
	}

	return respData.Transactions, nil
}

// DraftToRecipients is a draft transaction to a slice of recipients
func (g *TransportGraphQL) DraftToRecipients(ctx context.Context, recipients []*Recipients,
	metadata *bux.Metadata) (*bux.DraftTransaction, error) {

	reqBody := `
   	mutation ($outputs: [TransactionOutputInput]!, $metadata: Metadata) {
	  new_transaction(
		transaction_config:{
		  outputs: $outputs
          change_number_of_destinations:3
          change_destinations_strategy:"random"
		}
		metadata:$metadata
	  ) ` + graphqlDraftTransactionFields + `
	}`
	req := graphql.NewRequest(reqBody)
	outputs := make([]map[string]interface{}, 0)
	for _, recipient := range recipients {
		outputs = append(outputs, map[string]interface{}{
			FieldTo:       recipient.To,
			FieldSatoshis: recipient.Satoshis,
			FieldOpReturn: recipient.OpReturn,
		})
	}
	req.Var(FieldOutputs, outputs)
	req.Var(FieldMetadata, processMetadata(metadata))
	variables := map[string]interface{}{
		FieldOutputs:  outputs,
		FieldMetadata: processMetadata(metadata),
	}

	return g.draftTransactionCommon(ctx, reqBody, variables, req)
}

// DraftTransaction is a draft transaction
func (g *TransportGraphQL) DraftTransaction(ctx context.Context, transactionConfig *bux.TransactionConfig,
	metadata *bux.Metadata) (*bux.DraftTransaction, error) {

	reqBody := `
   	mutation ($transactionConfig: TransactionConfigInput!, $metadata: Metadata) {
	  new_transaction(
		transaction_config: $transactionConfig
		metadata: $metadata
	  ) ` + graphqlDraftTransactionFields + `
	}`
	req := graphql.NewRequest(reqBody)
	req.Var("transactionConfig", transactionConfig)
	req.Var(FieldMetadata, processMetadata(metadata))
	variables := map[string]interface{}{
		FieldTransactionConfig: transactionConfig,
		FieldMetadata:          processMetadata(metadata),
	}

	return g.draftTransactionCommon(ctx, reqBody, variables, req)
}

func (g *TransportGraphQL) draftTransactionCommon(ctx context.Context, reqBody string,
	variables map[string]interface{}, req *graphql.Request) (*bux.DraftTransaction, error) {

	err := g.signGraphQLRequest(req, reqBody, variables, g.xPriv, g.xPub)
	if err != nil {
		return nil, err
	}

	// run it and capture the response
	var respData DraftTransactionData
	if err = g.client.Run(ctx, req, &respData); err != nil {
		return nil, err
	}
	draftTransaction := respData.NewTransaction
	if g.debug {
		log.Printf("Draft transaction: %v\n", draftTransaction)
	}

	return draftTransaction, nil
}

// RecordTransaction will record a transaction
func (g *TransportGraphQL) RecordTransaction(ctx context.Context, hex, referenceID string,
	metadata *bux.Metadata) (*bux.Transaction, error) {

	reqBody := `
   	mutation($metadata: Metadata) {
	  transaction(
		hex:"` + hex + `",
        draft_id:"` + referenceID + `"
		metadata: $metadata
	  ) {
		id
		hex
		block_hash
		block_height
		fee
		number_of_inputs
		number_of_outputs
		output_value
		total_value
		direction
		metadata
		created_at
		updated_at
		deleted_at
	  }
	}`
	req := graphql.NewRequest(reqBody)
	req.Var(FieldMetadata, processMetadata(metadata))

	variables := map[string]interface{}{
		FieldMetadata: processMetadata(metadata),
	}
	err := g.signGraphQLRequest(req, reqBody, variables, g.xPriv, g.xPub)
	if err != nil {
		return nil, err
	}

	// run it and capture the response
	var respData NewTransactionData
	if err = g.client.Run(ctx, req, &respData); err != nil {
		return nil, err
	}
	transaction := respData.Transaction
	if g.debug {
		log.Printf("transaction: %s\n", transaction.ID)
	}

	return transaction, nil
}

// UpdateTransactionMetadata update the metadata of a transaction
func (g *TransportGraphQL) UpdateTransactionMetadata(ctx context.Context, txID string, metadata *bux.Metadata) (*bux.Transaction, error) {

	reqBody := `
    mutation ($id: String!, $metadata: Metadata!) {
  	  destination_metadata (
        id: $id
  	    metadata: $metadata
 	  ) {
        id
        hex
        block_hash
        block_height
        fee
        number_of_inputs
        number_of_outputs
        output_value
        total_value
        direction
        metadata
        created_at
        updated_at
        deleted_at
	  }
	}`
	variables := map[string]interface{}{
		FieldID:       txID,
		FieldMetadata: processMetadata(metadata),
	}

	var respData TransactionMetadataData
	if err := g.doGraphQLQuery(ctx, reqBody, variables, &respData); err != nil {
		return nil, err
	}

	return respData.Transaction, nil
}

func (g *TransportGraphQL) doGraphQLQuery(ctx context.Context, reqBody string, variables map[string]interface{},
	respData interface{}) error {

	req := graphql.NewRequest(reqBody)
	for key, value := range variables {
		req.Var(key, value)
	}

	err := g.signGraphQLRequest(req, reqBody, variables, g.xPriv, g.xPub)
	if err != nil {
		return err
	}

	// run it and capture the response
	if err = g.client.Run(ctx, req, &respData); err != nil {
		return err
	}
	if g.debug {
		log.Printf("model: %v\n", respData)
	}

	return nil
}

func getBodyString(reqBody string, variables map[string]interface{}) (string, error) {
	requestBodyObj := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}{
		Query:     reqBody,
		Variables: variables,
	}

	body, err := json.Marshal(requestBodyObj)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (g *TransportGraphQL) signGraphQLRequest(req *graphql.Request, reqBody string, variables map[string]interface{},
	xPriv *bip32.ExtendedKey, xPub *bip32.ExtendedKey) error {

	if g.signRequest {
		bodyString, err := getBodyString(reqBody, variables)
		if err != nil {
			return err
		}
		err = addSignature(&req.Header, xPriv, bodyString)
		if err != nil {
			return err
		}
	} else {
		req.Header.Set(bux.AuthHeader, xPub.String())
	}
	return nil
}

const graphqlDraftTransactionFields = `{
id
xpub_id
configuration {
  inputs {
	id
	satoshis
	transaction_id
	output_index
	script_pub_key
	destination {
	  id
	  address
	  type
	  num
	  chain
	  locking_script
	}
  }
  outputs {
	to
	satoshis
	scripts {
	  address
	  satoshis
	  script
	}
	paymail_p4 {
	  alias
	  domain
	  from_paymail
	  note
	  pub_key
	  receive_endpoint
      reference_id
	  resolution_type
	}
  }
  change_destinations {
	address
	chain
	num
	locking_script
	draft_id
  }
  change_satoshis
  fee
}
status
expires_at
hex
}`
