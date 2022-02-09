package transports

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/BuxOrg/bux"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
	"github.com/machinebox/graphql"
)

type graphQlService interface {
	Run(ctx context.Context, req *graphql.Request, resp interface{}) error
}

// TransportGraphQL is the graphql struct
type TransportGraphQL struct {
	accessKey   *bec.PrivateKey
	adminXPriv  *bip32.ExtendedKey
	debug       bool
	server      string
	signRequest bool
	xPriv       *bip32.ExtendedKey
	xPub        *bip32.ExtendedKey
	client      graphQlService
}

// DestinationData is the destination data
type DestinationData struct {
	Destination *bux.Destination `json:"destination"`
}

// DraftTransactionData is a draft transaction
type DraftTransactionData struct {
	NewTransaction *bux.DraftTransaction `json:"newTransaction"`
}

// Init will initialize
func (g *TransportGraphQL) Init() error {
	g.client = graphql.NewClient(g.server)
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

// SetSignRequest turn the signing of the http request on or off
func (g *TransportGraphQL) SetSignRequest(signRequest bool) {
	g.signRequest = signRequest
}

// IsSignRequest return whether to sign all requests
func (g *TransportGraphQL) IsSignRequest() bool {
	return g.signRequest
}

// RegisterXpub will register an xPub
func (g *TransportGraphQL) RegisterXpub(ctx context.Context, rawXPub string, metadata *bux.Metadata) error {

	// adding an xpub needs to be signed by an admin key
	if g.adminXPriv == nil {
		return ErrAdminKey
	}

	reqBody := `
   	mutation ($metadata: Map) {
	  xpub(
		xpub: "` + rawXPub + `"
		metadata: $metadata
	  ) {
	    ID
	  }
	}`
	req := graphql.NewRequest(reqBody)
	req.Var("metadata", processMetadata(metadata))
	variables := map[string]interface{}{
		"metadata": processMetadata(metadata),
	}

	bodyString, err := getBodyString(reqBody, variables)
	if err != nil {
		return err
	}
	err = addSignature(&req.Header, g.adminXPriv, bodyString)
	if err != nil {
		return err
	}

	// run it and capture the response
	var xPubData interface{}
	if err = g.client.Run(ctx, req, &xPubData); err != nil {
		return err
	}

	return nil
}

// GetDestination will get a destination
func (g *TransportGraphQL) GetDestination(ctx context.Context, metadata *bux.Metadata) (*bux.Destination, error) {
	reqBody := `
   	mutation ($metadata: Map) {
	  destination(
		metadata: $metadata
	  ) {
		ID
		XpubID
		LockingScript
		Type
		Chain
		Num
		Address
		Metadata
	  }
	}`
	req := graphql.NewRequest(reqBody)
	req.Var("metadata", processMetadata(metadata))
	if g.signRequest {
		variables := map[string]interface{}{
			"metadata": processMetadata(metadata),
		}
		bodyString, err := getBodyString(reqBody, variables)
		if err != nil {
			return nil, err
		}
		err = addSignature(&req.Header, g.xPriv, bodyString)
		if err != nil {
			return nil, err
		}
	} else {
		req.Header.Set("auth_xpub", g.xPub.String())
	}

	// run it and capture the response
	var respData DestinationData
	if err := g.client.Run(ctx, req, &respData); err != nil {
		return nil, err
	}
	destination := respData.Destination
	if g.debug {
		fmt.Printf("Address for new destination: %s\n", destination.Address)
	}

	return destination, nil
}

// DraftTransaction is a draft transaction
func (g *TransportGraphQL) DraftTransaction(ctx context.Context, transactionConfig *bux.TransactionConfig,
	metadata *bux.Metadata) (*bux.DraftTransaction, error) {

	reqBody := `
   	mutation ($transactionConfig: TransactionConfigInput!, $metadata: Map) {
	  newTransaction(
		transactionConfig: $transactionConfig
		metadata: $metadata
	  ) ` + graphqlDraftTransactionFields + `
	}`
	req := graphql.NewRequest(reqBody)
	req.Var("transactionConfig", transactionConfig)
	req.Var("metadata", processMetadata(metadata))
	variables := map[string]interface{}{
		"transactionConfig": transactionConfig,
		"metadata":          processMetadata(metadata),
	}

	return g.draftTransactionCommon(ctx, reqBody, variables, req)
}

// DraftToRecipients is a draft transaction to a slice of recipients
func (g *TransportGraphQL) DraftToRecipients(ctx context.Context, recipients []*Recipients,
	metadata *bux.Metadata) (*bux.DraftTransaction, error) {

	reqBody := `
   	mutation ($outputs: [TransactionOutputInput]!, $metadata: Map) {
	  newTransaction(
		transactionConfig:{
		  Outputs: $outputs
          ChangeNumberOfDestinations:3
          ChangeDestinationsStrategy:"random"
		}
		metadata:$metadata
	  ) ` + graphqlDraftTransactionFields + `
	}`
	req := graphql.NewRequest(reqBody)
	outputs := make([]map[string]interface{}, 0)
	for _, recipient := range recipients {
		outputs = append(outputs, map[string]interface{}{
			"To":       recipient.To,
			"Satoshis": recipient.Satoshis,
			"OpReturn": recipient.OpReturn,
		})
	}
	req.Var("outputs", outputs)
	req.Var("metadata", processMetadata(metadata))
	variables := map[string]interface{}{
		"outputs":  outputs,
		"metadata": processMetadata(metadata),
	}

	return g.draftTransactionCommon(ctx, reqBody, variables, req)
}

func (g *TransportGraphQL) draftTransactionCommon(ctx context.Context, reqBody string,
	variables map[string]interface{}, req *graphql.Request) (*bux.DraftTransaction, error) {

	if g.signRequest {
		bodyString, err := getBodyString(reqBody, variables)
		if err != nil {
			return nil, err
		}
		err = addSignature(&req.Header, g.xPriv, bodyString)
		if err != nil {
			return nil, err
		}
	} else {
		req.Header.Set("auth_xpub", g.xPub.String())
	}

	// run it and capture the response
	var respData DraftTransactionData
	if err := g.client.Run(ctx, req, &respData); err != nil {
		return nil, err
	}
	draftTransaction := respData.NewTransaction
	if g.debug {
		fmt.Printf("Draft transaction: %v\n", draftTransaction)
	}

	return draftTransaction, nil
}

// RecordTransaction will record a transaction
func (g *TransportGraphQL) RecordTransaction(ctx context.Context, hex, referenceID string,
	metadata *bux.Metadata) (string, error) {

	reqBody := `
   	mutation($metadata: Map) {
	  transaction(
		hex:"` + hex + `",
        draftID:"` + referenceID + `"
		metadata: $metadata
	  ) {
		ID
	  }
	}`
	req := graphql.NewRequest(reqBody)
	req.Var("metadata", processMetadata(metadata))

	if g.signRequest {
		variables := map[string]interface{}{
			"metadata": processMetadata(metadata),
		}
		bodyString, err := getBodyString(reqBody, variables)
		if err != nil {
			return "", err
		}
		err = addSignature(&req.Header, g.xPriv, bodyString)
		if err != nil {
			return "", err
		}
	} else {
		req.Header.Set("auth_xpub", g.xPub.String())
	}

	// run it and capture the response
	var transaction bux.Transaction
	if err := g.client.Run(ctx, req, &transaction); err != nil {
		return "", err
	}
	if g.debug {
		fmt.Printf("Transaction: %s\n", transaction.ID)
	}

	return transaction.ID, nil
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

const graphqlDraftTransactionFields = `{
ID
XpubID
Configuration {
  Inputs {
	ID
	Satoshis
	TransactionID
	OutputIndex
	ScriptPubKey
	Destination {
	  ID
	  Address
	  Type
	  Num
	  Chain
	  LockingScript
	}
  }
  Outputs {
	To
	Satoshis
	Scripts {
	  Address
	  Satoshis
	  Script
	}
	PaymailP4 {
	  Alias
	  Domain
	  FromPaymail
	  Note
	  PubKey
	  ReceiveEndpoint
	  ReferenceID
	  ResolutionType
	}
  }
  ChangeDestinations {
	Address
	Chain
	Num
	LockingScript
	DraftID
  }
  ChangeSatoshis
  Fee
}
Status
ExpiresAt
Hex
}`
