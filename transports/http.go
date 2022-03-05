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

const BuxVersion = "v1"

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
	err = h.doHTTPRequest(ctx, "POST", fmt.Sprintf("/%s/xpubs", BuxVersion), jsonStr, h.adminXPriv, true, &xPubData)
	if err != nil {
		return err
	}

	return nil
}

// GetDestination will get a destination
func (h *TransportHTTP) GetDestination(ctx context.Context, metadata *bux.Metadata) (*bux.Destination, error) {
	jsonData := map[string]interface{}{
		"metadata": processMetadata(metadata),
	}

	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	var destination bux.Destination
	err = h.doHTTPRequest(ctx, "POST", fmt.Sprintf("/%s/destinations", BuxVersion), jsonStr, h.xPriv, true, &destination)
	if err != nil {
		return nil, err
	}
	if h.debug {
		fmt.Printf("Address for new destination: %s\n", destination.Address)
	}

	return &destination, nil
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

func (h *TransportHTTP) createDraftTransaction(ctx context.Context, jsonData map[string]interface{}) (*bux.DraftTransaction, error) {
	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	var draftTransaction *bux.DraftTransaction
	err = h.doHTTPRequest(ctx, "POST", fmt.Sprintf("/%s/transactions/new", BuxVersion), jsonStr, h.xPriv, true, &draftTransaction)
	if err != nil {
		return nil, err
	}
	if h.debug {
		fmt.Printf("Draft transaction: %v\n", draftTransaction)
	}

	return draftTransaction, nil
}

// GetTransaction will get a transaction by ID
func (h *TransportHTTP) GetTransaction(ctx context.Context, txID string) (*bux.Transaction, error) {

	var transaction *bux.Transaction
	err := h.doHTTPRequest(ctx, "GET", fmt.Sprintf("/%s/transaction?id="+txID, BuxVersion), nil, h.xPriv, h.signRequest, &transaction)
	if err != nil {
		return nil, err
	}
	if h.debug {
		fmt.Printf("Transaction: %s\n", transaction.ID)
	}

	return transaction, nil
}

// GetTransactions will get a transactions by
func (h *TransportHTTP) GetTransactions(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata) ([]*bux.Transaction, error) {

	jsonData := map[string]interface{}{
		"conditions": conditions,
		"metadata":   processMetadata(metadata),
	}

	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	var transactions []*bux.Transaction
	err = h.doHTTPRequest(ctx, "POST", fmt.Sprintf("/%s/transactions", BuxVersion), jsonStr, h.xPriv, h.signRequest, &transactions)
	if err != nil {
		return nil, err
	}
	if h.debug {
		fmt.Printf("Transactions: %d\n", len(transactions))
	}

	return transactions, nil
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

	var transaction *bux.Transaction
	err = h.doHTTPRequest(ctx, "POST", fmt.Sprintf("/%s/transactions/record", BuxVersion), jsonStr, h.xPriv, h.signRequest, &transaction)
	if err != nil {
		return nil, err
	}
	if h.debug {
		fmt.Printf("Transaction: %s\n", transaction.ID)
	}

	return transaction, nil
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
