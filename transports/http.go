package transports

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	buxmodels "github.com/BuxOrg/bux-models"
	buxerrors "github.com/BuxOrg/bux-models/bux-errors"
	"github.com/BuxOrg/go-buxclient/utils"
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

// NewPaymail will register a new paymail
func (h *TransportHTTP) NewPaymail(ctx context.Context, rawXpub, paymailAddress, avatar, publicName string, metadata *buxmodels.Metadata) ResponseError {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldAddress:    paymailAddress,
		FieldAvatar:     avatar,
		FieldPublicName: publicName,
		FieldMetadata:   processMetadata(metadata),
		FieldXpubKey:    rawXpub,
	})
	if err != nil {
		return WrapError(err)
	}

	var paymailData interface{}

	return h.doHTTPRequest(
		ctx, http.MethodPost, "/paymail", jsonStr, h.xPriv, true, &paymailData,
	)
}

// GetXPub will get the xpub of the current xpub
func (h *TransportHTTP) GetXPub(ctx context.Context) (*buxmodels.Xpub, ResponseError) {
	var xPub buxmodels.Xpub
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

// UpdateXPubMetadata update the metadata of the logged in xpub
func (h *TransportHTTP) UpdateXPubMetadata(ctx context.Context, metadata *buxmodels.Metadata) (*buxmodels.Xpub, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: processMetadata(metadata),
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var xPub buxmodels.Xpub
	if err := h.doHTTPRequest(
		ctx, http.MethodPatch, "/xpub", jsonStr, h.xPriv, true, &xPub,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("xpub: %v\n", xPub)
	}

	return &xPub, nil
}

// GetAccessKey will get an access key by id
func (h *TransportHTTP) GetAccessKey(ctx context.Context, id string) (*buxmodels.AccessKey, ResponseError) {
	var accessKey buxmodels.AccessKey
	if err := h.doHTTPRequest(
		ctx, http.MethodGet, "/access-key?"+FieldID+"="+id, nil, h.xPriv, true, &accessKey,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("access key: %v\n", accessKey)
	}

	return &accessKey, nil
}

// GetAccessKeys will get all access keys matching the metadata filter
func (h *TransportHTTP) GetAccessKeys(ctx context.Context, metadataConditions *buxmodels.Metadata) ([]*buxmodels.AccessKey, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: processMetadata(metadataConditions),
	})
	if err != nil {
		return nil, WrapError(err)
	}
	var accessKey []*buxmodels.AccessKey
	if err := h.doHTTPRequest(
		ctx, http.MethodPost, "/access-key/search", jsonStr, h.xPriv, true, &accessKey,
	); err != nil {
		return nil, err
	}

	return accessKey, nil
}

// RevokeAccessKey will revoke an access key by id
func (h *TransportHTTP) RevokeAccessKey(ctx context.Context, id string) (*buxmodels.AccessKey, ResponseError) {
	var accessKey buxmodels.AccessKey
	if err := h.doHTTPRequest(
		ctx, http.MethodDelete, "/access-key?"+FieldID+"="+id, nil, h.xPriv, true, &accessKey,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("access key: %v\n", accessKey)
	}

	return &accessKey, nil
}

// CreateAccessKey will create new access key
func (h *TransportHTTP) CreateAccessKey(ctx context.Context, metadata *buxmodels.Metadata) (*buxmodels.AccessKey, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: processMetadata(metadata),
	})
	if err != nil {
		return nil, WrapError(err)
	}
	var accessKey buxmodels.AccessKey
	if err := h.doHTTPRequest(
		ctx, http.MethodPost, "/access-key", jsonStr, h.xPriv, true, &accessKey,
	); err != nil {
		return nil, err
	}

	return &accessKey, nil
}

// GetDestinationByID will get a destination by id
func (h *TransportHTTP) GetDestinationByID(ctx context.Context, id string) (*buxmodels.Destination, ResponseError) {
	var destination buxmodels.Destination
	if err := h.doHTTPRequest(
		ctx, http.MethodGet, "/destination?"+FieldID+"="+id, nil, h.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("destination: %v\n", destination)
	}

	return &destination, nil
}

// GetDestinationByAddress will get a destination by address
func (h *TransportHTTP) GetDestinationByAddress(ctx context.Context, address string) (*buxmodels.Destination, ResponseError) {
	var destination buxmodels.Destination
	if err := h.doHTTPRequest(
		ctx, http.MethodGet, "/destination?"+FieldAddress+"="+address, nil, h.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("destination: %v\n", destination)
	}

	return &destination, nil
}

// GetDestinationByLockingScript will get a destination by locking script
func (h *TransportHTTP) GetDestinationByLockingScript(ctx context.Context, lockingScript string) (*buxmodels.Destination, ResponseError) {
	var destination buxmodels.Destination
	if err := h.doHTTPRequest(
		ctx, http.MethodGet, "/destination?"+FieldLockingScript+"="+lockingScript, nil, h.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("destination: %v\n", destination)
	}

	return &destination, nil
}

// GetDestinations will get all destinations matching the metadata filter
func (h *TransportHTTP) GetDestinations(ctx context.Context, metadataConditions *buxmodels.Metadata) ([]*buxmodels.Destination, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: processMetadata(metadataConditions),
	})
	if err != nil {
		return nil, WrapError(err)
	}
	var destinations []*buxmodels.Destination
	if err := h.doHTTPRequest(
		ctx, http.MethodPost, "/destination/search", jsonStr, h.xPriv, true, &destinations,
	); err != nil {
		return nil, err
	}

	return destinations, nil
}

// NewDestination will create a new destination and return it
func (h *TransportHTTP) NewDestination(ctx context.Context, metadata *buxmodels.Metadata) (*buxmodels.Destination, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: processMetadata(metadata),
	})
	if err != nil {
		return nil, WrapError(err)
	}
	var destination buxmodels.Destination
	if err := h.doHTTPRequest(
		ctx, http.MethodPost, "/destination", jsonStr, h.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("new destination: %v\n", destination)
	}

	return &destination, nil
}

// UpdateDestinationMetadataByID updates the destination metadata by id
func (h *TransportHTTP) UpdateDestinationMetadataByID(ctx context.Context, id string,
	metadata *buxmodels.Metadata,
) (*buxmodels.Destination, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldID:       id,
		FieldMetadata: processMetadata(metadata),
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var destination buxmodels.Destination
	if err := h.doHTTPRequest(
		ctx, http.MethodPatch, "/destination", jsonStr, h.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("destination: %v\n", destination)
	}

	return &destination, nil
}

// UpdateDestinationMetadataByAddress updates the destination metadata by address
func (h *TransportHTTP) UpdateDestinationMetadataByAddress(ctx context.Context, address string,
	metadata *buxmodels.Metadata,
) (*buxmodels.Destination, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldAddress:  address,
		FieldMetadata: processMetadata(metadata),
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var destination buxmodels.Destination
	if err := h.doHTTPRequest(
		ctx, http.MethodPatch, "/destination", jsonStr, h.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("destination: %v\n", destination)
	}

	return &destination, nil
}

// UpdateDestinationMetadataByLockingScript updates the destination metadata by locking script
func (h *TransportHTTP) UpdateDestinationMetadataByLockingScript(ctx context.Context, lockingScript string,
	metadata *buxmodels.Metadata,
) (*buxmodels.Destination, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldLockingScript: lockingScript,
		FieldMetadata:      processMetadata(metadata),
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var destination buxmodels.Destination
	if err := h.doHTTPRequest(
		ctx, http.MethodPatch, "/destination", jsonStr, h.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("destination: %v\n", destination)
	}

	return &destination, nil
}

// GetTransaction will get a transaction by ID
func (h *TransportHTTP) GetTransaction(ctx context.Context, txID string) (*buxmodels.Transaction, ResponseError) {
	var transaction buxmodels.Transaction
	if err := h.doHTTPRequest(
		ctx, http.MethodGet, "/transaction?"+FieldID+"="+txID, nil, h.xPriv, h.signRequest, &transaction,
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
	metadataConditions *buxmodels.Metadata, queryParams *QueryParams,
) ([]*buxmodels.Transaction, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldConditions:  conditions,
		FieldMetadata:    processMetadata(metadataConditions),
		FieldQueryParams: queryParams,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var transactions []*buxmodels.Transaction
	if err := h.doHTTPRequest(
		ctx, http.MethodPost, "/transaction/search", jsonStr, h.xPriv, h.signRequest, &transactions,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("transactions: %d\n", len(transactions))
	}

	return transactions, nil
}

// GetTransactionsCount get number of user transactions
func (h *TransportHTTP) GetTransactionsCount(ctx context.Context, conditions map[string]interface{},
	metadata *buxmodels.Metadata,
) (int64, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldConditions: conditions,
		FieldMetadata:   processMetadata(metadata),
	})
	if err != nil {
		return 0, WrapError(err)
	}

	var count int64
	if err := h.doHTTPRequest(
		ctx, http.MethodPost, "/transaction/count", jsonStr, h.xPriv, h.signRequest, &count,
	); err != nil {
		return 0, err
	}
	if h.debug {
		log.Printf("Transactions count: %v\n", count)
	}

	return count, nil
}

// DraftToRecipients is a draft transaction to a slice of recipients
func (h *TransportHTTP) DraftToRecipients(ctx context.Context, recipients []*Recipients,
	metadata *buxmodels.Metadata,
) (*buxmodels.DraftTransaction, ResponseError) {
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
func (h *TransportHTTP) DraftTransaction(ctx context.Context, transactionConfig *buxmodels.TransactionConfig,
	metadata *buxmodels.Metadata,
) (*buxmodels.DraftTransaction, ResponseError) {
	return h.createDraftTransaction(ctx, map[string]interface{}{
		FieldConfig:   transactionConfig,
		FieldMetadata: processMetadata(metadata),
	})
}

// createDraftTransaction will create a draft transaction
func (h *TransportHTTP) createDraftTransaction(ctx context.Context,
	jsonData map[string]interface{},
) (*buxmodels.DraftTransaction, ResponseError) {
	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		return nil, WrapError(err)
	}

	var draftTransaction *buxmodels.DraftTransaction
	if err := h.doHTTPRequest(
		ctx, http.MethodPost, "/transaction", jsonStr, h.xPriv, true, &draftTransaction,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("draft transaction: %v\n", draftTransaction)
	}
	if draftTransaction == nil {
		return nil, WrapError(buxerrors.ErrDraftNotFound)
	}

	return draftTransaction, nil
}

// RecordTransaction will record a transaction
func (h *TransportHTTP) RecordTransaction(ctx context.Context, hex, referenceID string,
	metadata *buxmodels.Metadata,
) (*buxmodels.Transaction, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldHex:         hex,
		FieldReferenceID: referenceID,
		FieldMetadata:    processMetadata(metadata),
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var transaction buxmodels.Transaction
	if err := h.doHTTPRequest(
		ctx, http.MethodPost, "/transaction/record", jsonStr, h.xPriv, h.signRequest, &transaction,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("transaction: %s\n", transaction.ID)
	}

	return &transaction, nil
}

// UpdateTransactionMetadata update the metadata of a transaction
func (h *TransportHTTP) UpdateTransactionMetadata(ctx context.Context, txID string,
	metadata *buxmodels.Metadata,
) (*buxmodels.Transaction, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldID:       txID,
		FieldMetadata: processMetadata(metadata),
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var transaction buxmodels.Transaction
	if err := h.doHTTPRequest(
		ctx, http.MethodPatch, "/transaction", jsonStr, h.xPriv, h.signRequest, &transaction,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("Transaction: %v\n", transaction)
	}

	return &transaction, nil
}

// SetSignatureFromAccessKey will set the signature on the header for the request from an access key
func SetSignatureFromAccessKey(header *http.Header, privateKeyHex, bodyString string) ResponseError {
	// Create the signature
	authData, err := createSignatureAccessKey(privateKeyHex, bodyString)
	if err != nil {
		return WrapError(err)
	}

	// Set the auth header
	header.Set(buxmodels.AuthAccessKey, authData.AccessKey)

	return setSignatureHeaders(header, authData)
}

// UnreserveUtxos will unreserve utxos from draft transaction
func (h *TransportHTTP) UnreserveUtxos(ctx context.Context, referenceID string) ResponseError {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldReferenceID: referenceID,
	})
	if err != nil {
		return WrapError(err)
	}

	return h.doHTTPRequest(ctx, http.MethodPatch, "/utxos/unreserve", jsonStr, h.xPriv, h.signRequest, nil)
}

// createSignatureAccessKey will create a signature for the given access key & body contents
func createSignatureAccessKey(privateKeyHex, bodyString string) (payload *buxmodels.AuthPayload, err error) {
	// No key?
	if privateKeyHex == "" {
		err = buxerrors.ErrMissingAccessKey
		return
	}

	var privateKey *bec.PrivateKey
	if privateKey, err = bitcoin.PrivateKeyFromString(
		privateKeyHex,
	); err != nil {
		return
	}
	publicKey := privateKey.PubKey()

	// Get the xPub
	payload = new(buxmodels.AuthPayload)
	payload.AccessKey = hex.EncodeToString(publicKey.SerialiseCompressed())

	// auth_nonce is a random unique string to seed the signing message
	// this can be checked server side to make sure the request is not being replayed
	payload.AuthNonce, err = utils.RandomHex(32)
	if err != nil {
		return nil, err
	}

	return createSignatureCommon(payload, bodyString, privateKey)
}

// doHTTPRequest will create and submit the HTTP request
func (h *TransportHTTP) doHTTPRequest(ctx context.Context, method string, path string,
	rawJSON []byte, xPriv *bip32.ExtendedKey, sign bool, responseJSON interface{},
) ResponseError {
	req, err := http.NewRequestWithContext(ctx, method, h.server+path, bytes.NewBuffer(rawJSON))
	if err != nil {
		return WrapError(err)
	}
	req.Header.Set("Content-Type", "application/json")

	if xPriv != nil {
		err := h.authenticateWithXpriv(sign, req, xPriv, rawJSON)
		if err != nil {
			return err
		}
	} else {
		err := h.authenticateWithAccessKey(req, rawJSON)
		if err != nil {
			return err
		}
	}

	var resp *http.Response
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	if resp, err = h.httpClient.Do(req); err != nil {
		return WrapError(err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return WrapResponseError(resp)
	}

	if responseJSON == nil {
		return nil
	}

	err = json.NewDecoder(resp.Body).Decode(&responseJSON)
	if err != nil {
		return WrapError(err)
	}
	return nil
}

func (h *TransportHTTP) authenticateWithXpriv(sign bool, req *http.Request, xPriv *bip32.ExtendedKey, rawJSON []byte) ResponseError {
	if sign {
		if err := addSignature(&req.Header, xPriv, string(rawJSON)); err != nil {
			return err
		}
	} else {
		var xPub string
		xPub, err := bitcoin.GetExtendedPublicKey(xPriv)
		if err != nil {
			return WrapError(err)
		}
		req.Header.Set(buxmodels.AuthHeader, xPub)
		req.Header.Set("", xPub)
	}
	return nil
}

func (h *TransportHTTP) authenticateWithAccessKey(req *http.Request, rawJSON []byte) ResponseError {
	return SetSignatureFromAccessKey(&req.Header, hex.EncodeToString(h.accessKey.Serialise()), string(rawJSON))
}
