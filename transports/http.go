package transports

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/BuxOrg/bux"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
)

// TransportHTTP is the struct for HTTP
type TransportHTTP struct {
	accessKey   *bec.PrivateKey
	adminXPriv  *bip32.ExtendedKey
	debug       bool
	httpClient  *http.Client
	server      string
	signRequest bool
	xPriv       *bip32.ExtendedKey
	xPub        *bip32.ExtendedKey
}

// Init will initialize
func (h *TransportHTTP) Init() error {
	return nil
}

// SetDebug turn the debugging on or off
func (h *TransportHTTP) SetDebug(debug bool) {
	h.debug = debug
}

// IsDebug return the debugging status
func (h *TransportHTTP) IsDebug() bool {
	return h.debug
}

// SetSignRequest turn the signing of the http request on or off
func (h *TransportHTTP) SetSignRequest(signRequest bool) {
	h.signRequest = signRequest
}

// IsSignRequest return whether to sign all requests
func (h *TransportHTTP) IsSignRequest() bool {
	return h.signRequest
}

// SetAdminKey set the admin key
func (h *TransportHTTP) SetAdminKey(adminKey *bip32.ExtendedKey) {
	h.adminXPriv = adminKey
}

// RegisterPaymail will register a new paymail
func (h *TransportHTTP) RegisterPaymail(ctx context.Context, rawXpub, paymailAddress string, metadata *bux.Metadata) error {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldAddress:  paymailAddress,
		FieldMetadata: processMetadata(metadata),
		FieldXpubKey:  rawXpub,
	})
	if err != nil {
		return err
	}

	var paymailData interface{}
	if err = h.doHTTPRequest(
		ctx, http.MethodPost, "/paymail", jsonStr, h.xPriv, true, &paymailData,
	); err != nil {
		return err
	}

	return nil
}

// RegisterXpub will register an xPub
func (h *TransportHTTP) RegisterXpub(ctx context.Context, rawXPub string, metadata *bux.Metadata) error {

	// Adding a xpub needs to be signed by an admin key
	if h.adminXPriv == nil {
		return ErrAdminKey
	}

	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: processMetadata(metadata),
		FieldXpubKey:  rawXPub,
	})
	if err != nil {
		return err
	}

	var xPubData bux.Xpub
	if err = h.doHTTPRequest(
		ctx, http.MethodPost, "/xpub", jsonStr, h.adminXPriv, true, &xPubData,
	); err != nil {
		return err
	}

	return nil
}

// GetXPub will get the xpub of the current xpub
func (h *TransportHTTP) GetXPub(ctx context.Context) (*bux.Xpub, error) {
	var xPub bux.Xpub
	if err := h.doHTTPRequest(
		ctx, http.MethodGet, "/xpub", nil, h.xPriv, true, &xPub,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("xpub: %v\n", xPub)
	}

	return &xPub, nil
}

// GetAccessKey will get an access key by id
func (h *TransportHTTP) GetAccessKey(ctx context.Context, id string) (*bux.AccessKey, error) {
	var accessKey bux.AccessKey
	if err := h.doHTTPRequest(
		ctx, http.MethodGet, "/access-key?id="+id, nil, h.xPriv, true, &accessKey,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("access key: %v\n", accessKey)
	}

	return &accessKey, nil
}

// GetAccessKeys will get all access keys matching the metadata filter
func (h *TransportHTTP) GetAccessKeys(ctx context.Context, metadataConditions *bux.Metadata) ([]*bux.AccessKey, error) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: processMetadata(metadataConditions),
	})
	if err != nil {
		return nil, err
	}
	var accessKey []*bux.AccessKey
	if err = h.doHTTPRequest(
		ctx, http.MethodPost, "/access-key/search", jsonStr, h.xPriv, true, &accessKey,
	); err != nil {
		return nil, err
	}

	return accessKey, nil
}

// RevokeAccessKey will revoke an access key by id
func (h *TransportHTTP) RevokeAccessKey(ctx context.Context, id string) (*bux.AccessKey, error) {
	var accessKey bux.AccessKey
	if err := h.doHTTPRequest(
		ctx, http.MethodDelete, "/access-key?id="+id, nil, h.xPriv, true, &accessKey,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("access key: %v\n", accessKey)
	}

	return &accessKey, nil
}

// CreateAccessKey will create new access key
func (h *TransportHTTP) CreateAccessKey(ctx context.Context, metadata *bux.Metadata) (*bux.AccessKey, error) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: processMetadata(metadata),
	})
	if err != nil {
		return nil, err
	}
	var accessKey bux.AccessKey
	if err = h.doHTTPRequest(
		ctx, http.MethodPost, "/access-key", jsonStr, h.xPriv, true, &accessKey,
	); err != nil {
		return nil, err
	}

	return &accessKey, nil
}

// GetDestinationByID will get a destination by id
func (h *TransportHTTP) GetDestinationByID(ctx context.Context, id string) (*bux.Destination, error) {
	var destination bux.Destination
	if err := h.doHTTPRequest(
		ctx, http.MethodGet, "/destination?id="+id, nil, h.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("destination: %v\n", destination)
	}

	return &destination, nil
}

// GetDestinationByAddress will get a destination by address
func (h *TransportHTTP) GetDestinationByAddress(ctx context.Context, address string) (*bux.Destination, error) {
	var destination bux.Destination
	if err := h.doHTTPRequest(
		ctx, http.MethodGet, "/destination?address="+address, nil, h.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("destination: %v\n", destination)
	}

	return &destination, nil
}

// GetDestinationByLockingScript will get a destination by locking script
func (h *TransportHTTP) GetDestinationByLockingScript(ctx context.Context, lockingScript string) (*bux.Destination, error) {
	var destination bux.Destination
	if err := h.doHTTPRequest(
		ctx, http.MethodGet, "/destination?locking_script="+lockingScript, nil, h.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("destination: %v\n", destination)
	}

	return &destination, nil
}

// GetDestinations will get all destinations matching the metadata filter
func (h *TransportHTTP) GetDestinations(ctx context.Context, metadataConditions *bux.Metadata) ([]*bux.Destination, error) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: processMetadata(metadataConditions),
	})
	if err != nil {
		return nil, err
	}
	var destinations []*bux.Destination
	if err = h.doHTTPRequest(
		ctx, http.MethodPost, "/destination/search", jsonStr, h.xPriv, true, &destinations,
	); err != nil {
		return nil, err
	}

	return destinations, nil
}

// NewDestination will create a new destination and return it
func (h *TransportHTTP) NewDestination(ctx context.Context, metadata *bux.Metadata) (*bux.Destination, error) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: processMetadata(metadata),
	})
	if err != nil {
		return nil, err
	}
	var destination bux.Destination
	if err = h.doHTTPRequest(
		ctx, http.MethodPost, "/destination", jsonStr, h.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("new destination: %v\n", destination)
	}

	return &destination, nil
}

// GetTransaction will get a transaction by ID
func (h *TransportHTTP) GetTransaction(ctx context.Context, txID string) (*bux.Transaction, error) {
	var transaction bux.Transaction
	if err := h.doHTTPRequest(
		ctx, http.MethodGet, "/transaction?id="+txID, nil, h.xPriv, h.signRequest, &transaction,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("Transaction: %v\n", transaction)
	}

	return &transaction, nil
}

// GetTransactions will get a transactions by conditions
func (h *TransportHTTP) GetTransactions(ctx context.Context, conditions map[string]interface{},
	metadataConditions *bux.Metadata) ([]*bux.Transaction, error) {

	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldConditions: conditions,
		FieldMetadata:   processMetadata(metadataConditions),
	})
	if err != nil {
		return nil, err
	}

	var transactions []*bux.Transaction
	if err = h.doHTTPRequest(
		ctx, http.MethodPost, "/transaction/search", jsonStr, h.xPriv, h.signRequest, &transactions,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("transactions: %d\n", len(transactions))
	}

	return transactions, nil
}

// DraftToRecipients is a draft transaction to a slice of recipients
func (h *TransportHTTP) DraftToRecipients(ctx context.Context, recipients []*Recipients,
	metadata *bux.Metadata) (*bux.DraftTransaction, error) {

	outputs := make([]map[string]interface{}, 0)
	for _, recipient := range recipients {
		outputs = append(outputs, map[string]interface{}{
			FieldTo:       recipient.To,
			FieldSatoshis: recipient.Satoshis,
			FieldOpReturn: recipient.OpReturn,
		})
	}

	return h.createDraftTransaction(ctx, map[string]interface{}{
		FieldConfig: map[string]interface{}{
			FieldOutputs: outputs,
		},
		FieldMetadata: processMetadata(metadata),
	})
}

// DraftTransaction is a draft transaction
func (h *TransportHTTP) DraftTransaction(ctx context.Context, transactionConfig *bux.TransactionConfig,
	metadata *bux.Metadata) (*bux.DraftTransaction, error) {

	return h.createDraftTransaction(ctx, map[string]interface{}{
		FieldConfig:   transactionConfig,
		FieldMetadata: processMetadata(metadata),
	})
}

// createDraftTransaction will create a draft transaction
func (h *TransportHTTP) createDraftTransaction(ctx context.Context,
	jsonData map[string]interface{}) (*bux.DraftTransaction, error) {

	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	var draftTransaction bux.DraftTransaction
	if err = h.doHTTPRequest(
		ctx, http.MethodPost, "/transaction", jsonStr, h.xPriv, true, &draftTransaction,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("draft transaction: %v\n", draftTransaction)
	}

	return &draftTransaction, nil
}

// RecordTransaction will record a transaction
func (h *TransportHTTP) RecordTransaction(ctx context.Context, hex, referenceID string,
	metadata *bux.Metadata) (*bux.Transaction, error) {

	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldHex:         hex,
		FieldReferenceID: referenceID,
		FieldMetadata:    processMetadata(metadata),
	})
	if err != nil {
		return nil, err
	}

	var transaction bux.Transaction
	if err = h.doHTTPRequest(
		ctx, http.MethodPost, "/transaction/record", jsonStr, h.xPriv, h.signRequest, &transaction,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("transaction: %s\n", transaction.ID)
	}

	return &transaction, nil
}

// doHTTPRequest will create and submit the HTTP request
func (h *TransportHTTP) doHTTPRequest(ctx context.Context, method string, path string,
	rawJSON []byte, xPriv *bip32.ExtendedKey, sign bool, responseJSON interface{}) error {

	req, err := http.NewRequestWithContext(ctx, method, h.server+path, bytes.NewBuffer(rawJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if sign {
		if err = addSignature(&req.Header, xPriv, string(rawJSON)); err != nil {
			return err
		}
	} else {
		var xPub string
		if xPub, err = bitcoin.GetExtendedPublicKey(xPriv); err != nil {
			return err
		}
		req.Header.Set(bux.AuthHeader, xPub)
	}

	var resp *http.Response
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp, err = h.httpClient.Do(req); err != nil {
		return err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return errors.New("server error: " + strconv.Itoa(resp.StatusCode) + " - " + resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(&responseJSON)
}
