package transports

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	jsonData := map[string]interface{}{
		"metadata": processMetadata(metadata),
		"key":      rawXpub,
		"address":  paymailAddress,
	}

	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}

	var paymailData interface{}
	err = h.doHTTPRequest(ctx, "POST", "/paymail", jsonStr, h.xPriv, true, &paymailData)
	if err != nil {
		return err
	}

	return nil
}

// RegisterXpub will register an xPub
func (h *TransportHTTP) RegisterXpub(ctx context.Context, rawXPub string, metadata *bux.Metadata) error {

	// adding an xpub needs to be signed by an admin key
	if h.adminXPriv == nil {
		return ErrAdminKey
	}

	jsonData := map[string]interface{}{
		"metadata": processMetadata(metadata),
		"key":      rawXPub,
	}

	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}

	var xPubData bux.Xpub
	err = h.doHTTPRequest(ctx, "POST", "/xpub", jsonStr, h.adminXPriv, true, &xPubData)
	if err != nil {
		return err
	}

	return nil
}

// GetXPub will get the xpub of the current xpub
func (h *TransportHTTP) GetXPub(ctx context.Context) (*bux.Xpub, error) {
	var xPub bux.Xpub
	err := h.doHTTPRequest(ctx, "GET", "/xpub", nil, h.xPriv, true, &xPub)
	if err != nil {
		return nil, err
	}
	if h.debug {
		fmt.Printf("XPub: %v\n", xPub)
	}

	return &xPub, nil
}

// GetAccessKey will get an access key by id
func (h *TransportHTTP) GetAccessKey(ctx context.Context, id string) (*bux.AccessKey, error) {
	var accessKey bux.AccessKey
	err := h.doHTTPRequest(ctx, "GET", "/access-key?id="+id, nil, h.xPriv, true, &accessKey)
	if err != nil {
		return nil, err
	}
	if h.debug {
		fmt.Printf("Access key: %v\n", accessKey)
	}

	return &accessKey, nil
}

// GetAccessKeys will get all access keys matching the metadata filter
func (h *TransportHTTP) GetAccessKeys(ctx context.Context, metadataConditions *bux.Metadata) ([]*bux.AccessKey, error) {
	jsonData := map[string]interface{}{
		"metadata": processMetadata(metadataConditions),
	}

	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	var accessKey []*bux.AccessKey
	err = h.doHTTPRequest(ctx, "POST", "/access-key/search", jsonStr, h.xPriv, true, &accessKey)
	if err != nil {
		return nil, err
	}

	return accessKey, nil
}

// RevokeAccessKey will revoke an access key by id
func (h *TransportHTTP) RevokeAccessKey(ctx context.Context, id string) (*bux.AccessKey, error) {
	var accessKey bux.AccessKey
	err := h.doHTTPRequest(ctx, "DELETE", "/access-key?id="+id, nil, h.xPriv, true, &accessKey)
	if err != nil {
		return nil, err
	}
	if h.debug {
		fmt.Printf("Access key: %v\n", accessKey)
	}

	return &accessKey, nil
}

// CreateAccessKey will create a new access key
func (h *TransportHTTP) CreateAccessKey(ctx context.Context, metadata *bux.Metadata) (*bux.AccessKey, error) {
	jsonData := map[string]interface{}{
		"metadata": processMetadata(metadata),
	}

	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	var accessKey bux.AccessKey
	err = h.doHTTPRequest(ctx, "POST", "/access-key", jsonStr, h.xPriv, true, &accessKey)
	if err != nil {
		return nil, err
	}

	return &accessKey, nil
}

// GetDestinationByID will get a destination by id
func (h *TransportHTTP) GetDestinationByID(ctx context.Context, id string) (*bux.Destination, error) {
	var destination bux.Destination
	err := h.doHTTPRequest(ctx, "GET", "/destination?id="+id, nil, h.xPriv, true, &destination)
	if err != nil {
		return nil, err
	}
	if h.debug {
		fmt.Printf("Destination: %v\n", destination)
	}

	return &destination, nil
}

// GetDestinationByAddress will get a destination by address
func (h *TransportHTTP) GetDestinationByAddress(ctx context.Context, address string) (*bux.Destination, error) {
	var destination bux.Destination
	err := h.doHTTPRequest(ctx, "GET", "/destination?address="+address, nil, h.xPriv, true, &destination)
	if err != nil {
		return nil, err
	}
	if h.debug {
		fmt.Printf("Destination: %v\n", destination)
	}

	return &destination, nil
}

// GetDestinationByLockingScript will get a destination by locking script
func (h *TransportHTTP) GetDestinationByLockingScript(ctx context.Context, lockingScript string) (*bux.Destination, error) {
	var destination bux.Destination
	err := h.doHTTPRequest(ctx, "GET", "/destination?locking_script="+lockingScript, nil, h.xPriv, true, &destination)
	if err != nil {
		return nil, err
	}
	if h.debug {
		fmt.Printf("Destination: %v\n", destination)
	}

	return &destination, nil
}

// GetDestinations will get all destinations matching the metadata filter
func (h *TransportHTTP) GetDestinations(ctx context.Context, metadataConditions *bux.Metadata) ([]*bux.Destination, error) {
	jsonData := map[string]interface{}{
		"metadata": processMetadata(metadataConditions),
	}

	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	var destinations []*bux.Destination
	err = h.doHTTPRequest(ctx, "POST", "/destination/search", jsonStr, h.xPriv, true, &destinations)
	if err != nil {
		return nil, err
	}

	return destinations, nil
}

// NewDestination will create a new destination and return it
func (h *TransportHTTP) NewDestination(ctx context.Context, metadata *bux.Metadata) (*bux.Destination, error) {
	jsonData := map[string]interface{}{
		"metadata": processMetadata(metadata),
	}

	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	var destination bux.Destination
	err = h.doHTTPRequest(ctx, "POST", "/destination", jsonStr, h.xPriv, true, &destination)
	if err != nil {
		return nil, err
	}
	if h.debug {
		fmt.Printf("New destination: %v\n", destination)
	}

	return &destination, nil
}

// GetTransaction will get a transaction by ID
func (h *TransportHTTP) GetTransaction(ctx context.Context, txID string) (*bux.Transaction, error) {

	var transaction bux.Transaction
	err := h.doHTTPRequest(ctx, "GET", "/transaction?id="+txID, nil, h.xPriv, h.signRequest, &transaction)
	if err != nil {
		return nil, err
	}
	if h.debug {
		fmt.Printf("Transaction: %v\n", transaction)
	}

	return &transaction, nil
}

// GetTransactions will get a transactions by
func (h *TransportHTTP) GetTransactions(ctx context.Context, conditions map[string]interface{},
	metadataConditions *bux.Metadata) ([]*bux.Transaction, error) {

	jsonData := map[string]interface{}{
		"conditions": conditions,
		"metadata":   processMetadata(metadataConditions),
	}

	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	var transactions []*bux.Transaction
	err = h.doHTTPRequest(ctx, "POST", "/transaction/search", jsonStr, h.xPriv, h.signRequest, &transactions)
	if err != nil {
		return nil, err
	}
	if h.debug {
		fmt.Printf("Transactions: %d\n", len(transactions))
	}

	return transactions, nil
}

// DraftToRecipients is a draft transaction to a slice of recipients
func (h *TransportHTTP) DraftToRecipients(ctx context.Context, recipients []*Recipients,
	metadata *bux.Metadata) (*bux.DraftTransaction, error) {

	outputs := make([]map[string]interface{}, 0)
	for _, recipient := range recipients {
		outputs = append(outputs, map[string]interface{}{
			"to":        recipient.To,
			"satoshis":  recipient.Satoshis,
			"op_return": recipient.OpReturn,
		})
	}
	jsonData := map[string]interface{}{
		"config": map[string]interface{}{
			"outputs": outputs,
		},
		"metadata": processMetadata(metadata),
	}

	return h.createDraftTransaction(ctx, jsonData)
}

// DraftTransaction is a draft transaction
func (h *TransportHTTP) DraftTransaction(ctx context.Context, transactionConfig *bux.TransactionConfig,
	metadata *bux.Metadata) (*bux.DraftTransaction, error) {

	jsonData := map[string]interface{}{
		"config":   transactionConfig,
		"metadata": processMetadata(metadata),
	}

	return h.createDraftTransaction(ctx, jsonData)
}

func (h *TransportHTTP) createDraftTransaction(ctx context.Context, jsonData map[string]interface{}) (*bux.DraftTransaction, error) {
	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	var draftTransaction bux.DraftTransaction
	err = h.doHTTPRequest(ctx, "POST", "/transaction", jsonStr, h.xPriv, true, &draftTransaction)
	if err != nil {
		return nil, err
	}
	if h.debug {
		fmt.Printf("Draft transaction: %v\n", draftTransaction)
	}

	return &draftTransaction, nil
}

// RecordTransaction will record a transaction
func (h *TransportHTTP) RecordTransaction(ctx context.Context, hex, referenceID string,
	metadata *bux.Metadata) (*bux.Transaction, error) {

	jsonData := map[string]interface{}{
		"hex":          hex,
		"reference_id": referenceID,
		"metadata":     processMetadata(metadata),
	}

	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	var transaction bux.Transaction
	err = h.doHTTPRequest(ctx, "POST", "/transaction/record", jsonStr, h.xPriv, h.signRequest, &transaction)
	if err != nil {
		return nil, err
	}
	if h.debug {
		fmt.Printf("Transaction: %s\n", transaction.ID)
	}

	return &transaction, nil
}

func (h *TransportHTTP) doHTTPRequest(ctx context.Context, method string, path string, jsonStr []byte, xPriv *bip32.ExtendedKey, sign bool, responseJSON interface{}) error {

	url := h.server + path
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if sign {
		err = addSignature(&req.Header, xPriv, string(jsonStr))
		if err != nil {
			return err
		}
	} else {
		var xPub string
		xPub, err = bitcoin.GetExtendedPublicKey(xPriv)
		if err != nil {
			return err
		}
		req.Header.Set("auth_xpub", xPub)
	}

	resp, err := h.httpClient.Do(req) //nolint:bodyclose // done in defer function
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New("server error: " + strconv.Itoa(resp.StatusCode) + " - " + resp.Status)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	err = json.NewDecoder(resp.Body).Decode(&responseJSON)
	if err != nil {
		return err
	}

	return nil
}
