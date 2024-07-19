package walletclient

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	bec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet-go-client/utils"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoinschema/go-bitcoin/v2"
)

// TODO: Issues with "github.com/bitcoinschema/go-bitcoin/v2"

// SetSignRequest turn the signing of the http request on or off
func (wc *WalletClient) SetSignRequest(signRequest bool) {
	wc.signRequest = signRequest
}

// IsSignRequest return whether to sign all requests
func (wc *WalletClient) IsSignRequest() bool {
	return wc.signRequest
}

// SetAdminKey set the admin key
func (wc *WalletClient) SetAdminKey(adminKey *bip32.ExtendedKey) {
	wc.adminXPriv = adminKey
}

// GetXPub will get the xpub of the current xpub
func (wc *WalletClient) GetXPub(ctx context.Context) (*models.Xpub, error) {
	var xPub models.Xpub
	if err := wc.doHTTPRequest(
		ctx, http.MethodGet, "/xpub", nil, wc.xPriv, true, &xPub,
	); err != nil {
		return nil, err
	}

	return &xPub, nil
}

// UpdateXPubMetadata update the metadata of the logged in xpub
func (wc *WalletClient) UpdateXPubMetadata(ctx context.Context, metadata map[string]any) (*models.Xpub, error) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: metadata,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var xPub models.Xpub
	if err := wc.doHTTPRequest(
		ctx, http.MethodPatch, "/xpub", jsonStr, wc.xPriv, true, &xPub,
	); err != nil {
		return nil, err
	}

	return &xPub, nil
}

// GetAccessKey will get an access key by id
func (wc *WalletClient) GetAccessKey(ctx context.Context, id string) (*models.AccessKey, error) {
	var accessKey models.AccessKey
	if err := wc.doHTTPRequest(
		ctx, http.MethodGet, "/access-key?"+FieldID+"="+id, nil, wc.xPriv, true, &accessKey,
	); err != nil {
		return nil, err
	}

	return &accessKey, nil
}

// GetAccessKeys will get all access keys matching the metadata filter
func (wc *WalletClient) GetAccessKeys(
	ctx context.Context,
	conditions *filter.AccessKeyFilter,
	metadata map[string]any,
	queryParams *filter.QueryParams,
) ([]*models.AccessKey, error) {
	return Search[filter.AccessKeyFilter, []*models.AccessKey](
		ctx, http.MethodPost,
		"/access-key/search",
		wc.xPriv,
		conditions,
		metadata,
		queryParams,
		wc.doHTTPRequest,
	)
}

// GetAccessKeysCount will get the count of access keys
func (wc *WalletClient) GetAccessKeysCount(
	ctx context.Context,
	conditions *filter.AccessKeyFilter,
	metadata map[string]any,
) (int64, error) {
	return Count[filter.AccessKeyFilter](
		ctx, http.MethodPost,
		"/access-key/count",
		wc.xPriv,
		conditions,
		metadata,
		wc.doHTTPRequest,
	)
}

// RevokeAccessKey will revoke an access key by id
func (wc *WalletClient) RevokeAccessKey(ctx context.Context, id string) (*models.AccessKey, error) {
	var accessKey models.AccessKey
	if err := wc.doHTTPRequest(
		ctx, http.MethodDelete, "/access-key?"+FieldID+"="+id, nil, wc.xPriv, true, &accessKey,
	); err != nil {
		return nil, err
	}

	return &accessKey, nil
}

// CreateAccessKey will create new access key
func (wc *WalletClient) CreateAccessKey(ctx context.Context, metadata map[string]any) (*models.AccessKey, error) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: metadata,
	})
	if err != nil {
		return nil, WrapError(err)
	}
	var accessKey models.AccessKey
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, "/access-key", jsonStr, wc.xPriv, true, &accessKey,
	); err != nil {
		return nil, err
	}

	return &accessKey, nil
}

// GetDestinationByID will get a destination by id
func (wc *WalletClient) GetDestinationByID(ctx context.Context, id string) (*models.Destination, error) {
	var destination models.Destination
	if err := wc.doHTTPRequest(
		ctx, http.MethodGet, fmt.Sprintf("/destination?%s=%s", FieldID, id), nil, wc.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}

	return &destination, nil
}

// GetDestinationByAddress will get a destination by address
func (wc *WalletClient) GetDestinationByAddress(ctx context.Context, address string) (*models.Destination, error) {
	var destination models.Destination
	if err := wc.doHTTPRequest(
		ctx, http.MethodGet, "/destination?"+FieldAddress+"="+address, nil, wc.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}

	return &destination, nil
}

// GetDestinationByLockingScript will get a destination by locking script
func (wc *WalletClient) GetDestinationByLockingScript(ctx context.Context, lockingScript string) (*models.Destination, error) {
	var destination models.Destination
	if err := wc.doHTTPRequest(
		ctx, http.MethodGet, "/destination?"+FieldLockingScript+"="+lockingScript, nil, wc.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}

	return &destination, nil
}

// GetDestinations will get all destinations matching the metadata filter
func (wc *WalletClient) GetDestinations(ctx context.Context, conditions *filter.DestinationFilter, metadata map[string]any, queryParams *filter.QueryParams) ([]*models.Destination, error) {
	return Search[filter.DestinationFilter, []*models.Destination](
		ctx, http.MethodPost,
		"/destination/search",
		wc.xPriv,
		conditions,
		metadata,
		queryParams,
		wc.doHTTPRequest,
	)
}

// GetDestinationsCount will get the count of destinations matching the metadata filter
func (wc *WalletClient) GetDestinationsCount(ctx context.Context, conditions *filter.DestinationFilter, metadata map[string]any) (int64, error) {
	return Count(
		ctx,
		http.MethodPost,
		"/destination/count",
		wc.xPriv,
		conditions,
		metadata,
		wc.doHTTPRequest,
	)
}

// NewDestination will create a new destination and return it
func (wc *WalletClient) NewDestination(ctx context.Context, metadata map[string]any) (*models.Destination, error) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: metadata,
	})
	if err != nil {
		return nil, WrapError(err)
	}
	var destination models.Destination
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, "/destination", jsonStr, wc.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}

	return &destination, nil
}

// UpdateDestinationMetadataByID updates the destination metadata by id
func (wc *WalletClient) UpdateDestinationMetadataByID(ctx context.Context, id string, metadata map[string]any) (*models.Destination, error) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldID:       id,
		FieldMetadata: metadata,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var destination models.Destination
	if err := wc.doHTTPRequest(
		ctx, http.MethodPatch, "/destination", jsonStr, wc.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}

	return &destination, nil
}

// UpdateDestinationMetadataByAddress updates the destination metadata by address
func (wc *WalletClient) UpdateDestinationMetadataByAddress(ctx context.Context, address string, metadata map[string]any) (*models.Destination, error) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldAddress:  address,
		FieldMetadata: metadata,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var destination models.Destination
	if err := wc.doHTTPRequest(
		ctx, http.MethodPatch, "/destination", jsonStr, wc.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}

	return &destination, nil
}

// UpdateDestinationMetadataByLockingScript updates the destination metadata by locking script
func (wc *WalletClient) UpdateDestinationMetadataByLockingScript(ctx context.Context, lockingScript string, metadata map[string]any) (*models.Destination, error) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldLockingScript: lockingScript,
		FieldMetadata:      metadata,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var destination models.Destination
	if err := wc.doHTTPRequest(
		ctx, http.MethodPatch, "/destination", jsonStr, wc.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}

	return &destination, nil
}

// GetTransaction will get a transaction by ID
func (wc *WalletClient) GetTransaction(ctx context.Context, txID string) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := wc.doHTTPRequest(ctx, http.MethodGet, "/transaction?"+FieldID+"="+txID, nil, wc.xPriv, wc.signRequest, &transaction); err != nil {
		return nil, err
	}

	return &transaction, nil
}

// GetTransactions will get transactions by conditions
func (wc *WalletClient) GetTransactions(
	ctx context.Context,
	conditions *filter.TransactionFilter,
	metadata map[string]any,
	queryParams *filter.QueryParams,
) ([]*models.Transaction, error) {
	return Search[filter.TransactionFilter, []*models.Transaction](
		ctx, http.MethodPost,
		"/transaction/search",
		wc.xPriv,
		conditions,
		metadata,
		queryParams,
		wc.doHTTPRequest,
	)
}

// GetTransactionsCount get number of user transactions
func (wc *WalletClient) GetTransactionsCount(
	ctx context.Context,
	conditions *filter.TransactionFilter,
	metadata map[string]any,
) (int64, error) {
	return Count[filter.TransactionFilter](
		ctx, http.MethodPost,
		"/transaction/count",
		wc.xPriv,
		conditions,
		metadata,
		wc.doHTTPRequest,
	)
}

// DraftToRecipients is a draft transaction to a slice of recipients
func (wc *WalletClient) DraftToRecipients(ctx context.Context, recipients []*Recipients, metadata map[string]any) (*models.DraftTransaction, error) {
	outputs := make([]map[string]interface{}, 0)
	for _, recipient := range recipients {
		outputs = append(outputs, map[string]interface{}{
			FieldTo:       recipient.To,
			FieldSatoshis: recipient.Satoshis,
			FieldOpReturn: recipient.OpReturn,
		})
	}

	return wc.createDraftTransaction(ctx, map[string]interface{}{
		FieldConfig: map[string]interface{}{
			FieldOutputs: outputs,
		},
		FieldMetadata: metadata,
	})
}

// DraftTransaction is a draft transaction
func (wc *WalletClient) DraftTransaction(ctx context.Context, transactionConfig *models.TransactionConfig, metadata map[string]any) (*models.DraftTransaction, error) {
	return wc.createDraftTransaction(ctx, map[string]interface{}{
		FieldConfig:   transactionConfig,
		FieldMetadata: metadata,
	})
}

// createDraftTransaction will create a draft transaction
func (wc *WalletClient) createDraftTransaction(ctx context.Context,
	jsonData map[string]interface{},
) (*models.DraftTransaction, error) {
	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		return nil, WrapError(err)
	}

	var draftTransaction *models.DraftTransaction
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, "/transaction", jsonStr, wc.xPriv, true, &draftTransaction,
	); err != nil {
		return nil, err
	}
	if draftTransaction == nil {
		return nil, ErrCouldNotFindDraftTransaction
	}

	return draftTransaction, nil
}

// RecordTransaction will record a transaction
func (wc *WalletClient) RecordTransaction(ctx context.Context, hex, referenceID string, metadata map[string]any) (*models.Transaction, error) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldHex:         hex,
		FieldReferenceID: referenceID,
		FieldMetadata:    metadata,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var transaction models.Transaction
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, "/transaction/record", jsonStr, wc.xPriv, wc.signRequest, &transaction,
	); err != nil {
		return nil, err
	}

	return &transaction, nil
}

// UpdateTransactionMetadata update the metadata of a transaction
func (wc *WalletClient) UpdateTransactionMetadata(ctx context.Context, txID string, metadata map[string]any) (*models.Transaction, error) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldID:       txID,
		FieldMetadata: metadata,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var transaction models.Transaction
	if err := wc.doHTTPRequest(
		ctx, http.MethodPatch, "/transaction", jsonStr, wc.xPriv, wc.signRequest, &transaction,
	); err != nil {
		return nil, err
	}

	return &transaction, nil
}

// SetSignatureFromAccessKey will set the signature on the header for the request from an access key
func SetSignatureFromAccessKey(header *http.Header, privateKeyHex, bodyString string) error {
	// Create the signature
	authData, err := createSignatureAccessKey(privateKeyHex, bodyString)
	if err != nil {
		return WrapError(err)
	}

	// Set the auth header
	header.Set(models.AuthAccessKey, authData.AccessKey)

	setSignatureHeaders(header, authData)

	return nil
}

// GetUtxo will get a utxo by transaction ID
func (wc *WalletClient) GetUtxo(ctx context.Context, txID string, outputIndex uint32) (*models.Utxo, error) {
	outputIndexStr := strconv.FormatUint(uint64(outputIndex), 10)

	url := fmt.Sprintf("/utxo?%s=%s&%s=%s", FieldTransactionID, txID, FieldOutputIndex, outputIndexStr)

	var utxo models.Utxo
	if err := wc.doHTTPRequest(
		ctx, http.MethodGet, url, nil, wc.xPriv, true, &utxo,
	); err != nil {
		return nil, err
	}

	return &utxo, nil
}

// GetUtxos will get a list of utxos filtered by conditions and metadata
func (wc *WalletClient) GetUtxos(ctx context.Context, conditions *filter.UtxoFilter, metadata map[string]any, queryParams *filter.QueryParams) ([]*models.Utxo, error) {
	return Search[filter.UtxoFilter, []*models.Utxo](
		ctx, http.MethodPost,
		"/utxo/search",
		wc.xPriv,
		conditions,
		metadata,
		queryParams,
		wc.doHTTPRequest,
	)
}

// GetUtxosCount will get the count of utxos filtered by conditions and metadata
func (wc *WalletClient) GetUtxosCount(ctx context.Context, conditions *filter.UtxoFilter, metadata map[string]any) (int64, error) {
	return Count[filter.UtxoFilter](
		ctx, http.MethodPost,
		"/utxo/count",
		wc.xPriv,
		conditions,
		metadata,
		wc.doHTTPRequest,
	)
}

// createSignatureAccessKey will create a signature for the given access key & body contents
func createSignatureAccessKey(privateKeyHex, bodyString string) (payload *models.AuthPayload, err error) {
	// No key?
	if privateKeyHex == "" {
		err = CreateErrorResponse("error-unauthorized-missing-access-key", "missing access key")
		return
	}

	var privateKey *bec.PrivateKey
	if privateKey, err = bitcoin.PrivateKeyFromString(
		privateKeyHex,
	); err != nil {
		return
	}
	publicKey := privateKey.PubKey()

	// Get the AccessKey
	payload = new(models.AuthPayload)
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
func (wc *WalletClient) doHTTPRequest(ctx context.Context, method string, path string,
	rawJSON []byte, xPriv *bip32.ExtendedKey, sign bool, responseJSON interface{},
) error {
	req, err := http.NewRequestWithContext(ctx, method, wc.server+path, bytes.NewBuffer(rawJSON))
	if err != nil {
		return WrapError(err)
	}
	req.Header.Set("Content-Type", "application/json")

	if xPriv != nil {
		err := wc.authenticateWithXpriv(sign, req, xPriv, rawJSON)
		if err != nil {
			return err
		}
	} else {
		err := wc.authenticateWithAccessKey(req, rawJSON)
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
	if resp, err = wc.httpClient.Do(req); err != nil {
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

func (wc *WalletClient) authenticateWithXpriv(sign bool, req *http.Request, xPriv *bip32.ExtendedKey, rawJSON []byte) error {
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
		req.Header.Set(models.AuthHeader, xPub)
		req.Header.Set("", xPub)
	}
	return nil
}

func (wc *WalletClient) authenticateWithAccessKey(req *http.Request, rawJSON []byte) error {
	if wc.accessKey == nil {
		return WrapError(errors.New("access key is missing"))
	}
	return SetSignatureFromAccessKey(&req.Header, hex.EncodeToString(wc.accessKey.Serialise()), string(rawJSON))
}

// AcceptContact will accept the contact associated with the paymail
func (wc *WalletClient) AcceptContact(ctx context.Context, paymail string) error {
	if err := wc.doHTTPRequest(
		ctx, http.MethodPatch, "/contact/accepted/"+paymail, nil, wc.xPriv, wc.signRequest, nil,
	); err != nil {
		return err
	}

	return nil
}

// RejectContact will reject the contact associated with the paymail
func (wc *WalletClient) RejectContact(ctx context.Context, paymail string) error {
	if err := wc.doHTTPRequest(
		ctx, http.MethodPatch, "/contact/rejected/"+paymail, nil, wc.xPriv, wc.signRequest, nil,
	); err != nil {
		return err
	}

	return nil
}

// ConfirmContact will confirm the contact associated with the paymail
func (wc *WalletClient) ConfirmContact(ctx context.Context, contact *models.Contact, passcode, requesterPaymail string, period, digits uint) error {
	isTotpValid, err := wc.ValidateTotpForContact(contact, passcode, requesterPaymail, period, digits)
	if err != nil {
		return WrapError(fmt.Errorf("totp validation failed: %w", err))
	}

	if !isTotpValid {
		return WrapError(errors.New("totp is invalid"))
	}

	if err := wc.doHTTPRequest(
		ctx, http.MethodPatch, "/contact/confirmed/"+contact.Paymail, nil, wc.xPriv, wc.signRequest, nil,
	); err != nil {
		return err
	}

	return nil
}

// GetContacts will get contacts by conditions
func (wc *WalletClient) GetContacts(ctx context.Context, conditions *filter.ContactFilter, metadata map[string]any, queryParams *filter.QueryParams) (*models.SearchContactsResponse, error) {
	return Search[filter.ContactFilter, *models.SearchContactsResponse](
		ctx, http.MethodPost,
		"/contact/search",
		wc.xPriv,
		conditions,
		metadata,
		queryParams,
		wc.doHTTPRequest,
	)
}

// UpsertContact add or update contact. When adding a new contact, the system utilizes Paymail's PIKE capability to dispatch an invitation request, asking the counterparty to include the current user in their contacts.
func (wc *WalletClient) UpsertContact(ctx context.Context, paymail, fullName, requesterPaymail string, metadata map[string]any) (*models.Contact, error) {
	return wc.UpsertContactForPaymail(ctx, paymail, fullName, metadata, requesterPaymail)
}

// UpsertContactForPaymail add or update contact. When adding a new contact, the system utilizes Paymail's PIKE capability to dispatch an invitation request, asking the counterparty to include the current user in their contacts.
func (wc *WalletClient) UpsertContactForPaymail(ctx context.Context, paymail, fullName string, metadata map[string]any, requesterPaymail string) (*models.Contact, error) {
	payload := map[string]interface{}{
		"fullName":    fullName,
		FieldMetadata: metadata,
	}

	if requesterPaymail != "" {
		payload["requesterPaymail"] = requesterPaymail
	}

	jsonStr, err := json.Marshal(payload)
	if err != nil {
		return nil, WrapError(err)
	}

	var result models.Contact
	if err := wc.doHTTPRequest(
		ctx, http.MethodPut, "/contact/"+paymail, jsonStr, wc.xPriv, wc.signRequest, &result,
	); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSharedConfig gets the shared config
func (wc *WalletClient) GetSharedConfig(ctx context.Context) (*models.SharedConfig, error) {
	var model *models.SharedConfig

	key := wc.xPriv
	if wc.adminXPriv != nil {
		key = wc.adminXPriv
	}
	if key == nil {
		return nil, WrapError(errors.New("neither xPriv nor adminXPriv is provided"))
	}
	if err := wc.doHTTPRequest(
		ctx, http.MethodGet, "/shared-config", nil, key, true, &model,
	); err != nil {
		return nil, err
	}

	return model, nil
}

// AdminNewXpub will register an xPub
func (wc *WalletClient) AdminNewXpub(ctx context.Context, rawXPub string, metadata map[string]any) error {
	// Adding a xpub needs to be signed by an admin key
	if wc.adminXPriv == nil {
		return WrapError(ErrAdminKey)
	}

	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: metadata,
		FieldXpubKey:  rawXPub,
	})
	if err != nil {
		return WrapError(err)
	}

	var xPubData models.Xpub

	return wc.doHTTPRequest(
		ctx, http.MethodPost, "/admin/xpub", jsonStr, wc.adminXPriv, true, &xPubData,
	)
}

// AdminGetStatus get whether admin key is valid
func (wc *WalletClient) AdminGetStatus(ctx context.Context) (bool, error) {
	var status bool
	if err := wc.doHTTPRequest(
		ctx, http.MethodGet, "/admin/status", nil, wc.adminXPriv, true, &status,
	); err != nil {
		return false, err
	}

	return status, nil
}

// AdminGetStats get admin stats
func (wc *WalletClient) AdminGetStats(ctx context.Context) (*models.AdminStats, error) {
	var stats *models.AdminStats
	if err := wc.doHTTPRequest(
		ctx, http.MethodGet, "/admin/stats", nil, wc.adminXPriv, true, &stats,
	); err != nil {
		return nil, err
	}

	return stats, nil
}

// AdminGetAccessKeys get all access keys filtered by conditions
func (wc *WalletClient) AdminGetAccessKeys(
	ctx context.Context,
	conditions *filter.AdminAccessKeyFilter,
	metadata map[string]any,
	queryParams *filter.QueryParams,
) ([]*models.AccessKey, error) {
	return Search[filter.AdminAccessKeyFilter, []*models.AccessKey](
		ctx, http.MethodPost,
		"/admin/access-keys/search",
		wc.adminXPriv,
		conditions,
		metadata,
		queryParams,
		wc.doHTTPRequest,
	)
}

// AdminGetAccessKeysCount get a count of all the access keys filtered by conditions
func (wc *WalletClient) AdminGetAccessKeysCount(
	ctx context.Context,
	conditions *filter.AdminAccessKeyFilter,
	metadata map[string]any,
) (int64, error) {
	return Count[filter.AdminAccessKeyFilter](
		ctx, http.MethodPost,
		"/admin/access-keys/count",
		wc.adminXPriv,
		conditions,
		metadata,
		wc.doHTTPRequest,
	)
}

// AdminGetBlockHeaders get all block headers filtered by conditions
func (wc *WalletClient) AdminGetBlockHeaders(
	ctx context.Context,
	conditions map[string]interface{},
	metadata map[string]any,
	queryParams *filter.QueryParams,
) ([]*models.BlockHeader, error) {
	var models []*models.BlockHeader
	if err := wc.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/block-headers/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetBlockHeadersCount get a count of all the block headers filtered by conditions
func (wc *WalletClient) AdminGetBlockHeadersCount(
	ctx context.Context,
	conditions map[string]interface{},
	metadata map[string]any,
) (int64, error) {
	return wc.adminCount(ctx, conditions, metadata, "/admin/block-headers/count")
}

// AdminGetDestinations get all block destinations filtered by conditions
func (wc *WalletClient) AdminGetDestinations(ctx context.Context, conditions *filter.DestinationFilter,
	metadata map[string]any, queryParams *filter.QueryParams,
) ([]*models.Destination, error) {
	return Search[filter.DestinationFilter, []*models.Destination](
		ctx, http.MethodPost,
		"/admin/destinations/search",
		wc.adminXPriv,
		conditions,
		metadata,
		queryParams,
		wc.doHTTPRequest,
	)
}

// AdminGetDestinationsCount get a count of all the destinations filtered by conditions
func (wc *WalletClient) AdminGetDestinationsCount(ctx context.Context, conditions *filter.DestinationFilter, metadata map[string]any) (int64, error) {
	return Count(
		ctx,
		http.MethodPost,
		"/admin/destinations/count",
		wc.adminXPriv,
		conditions,
		metadata,
		wc.doHTTPRequest,
	)
}

// AdminGetPaymail get a paymail by address
func (wc *WalletClient) AdminGetPaymail(ctx context.Context, address string) (*models.PaymailAddress, error) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldAddress: address,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var model *models.PaymailAddress
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, "/admin/paymail/get", jsonStr, wc.adminXPriv, true, &model,
	); err != nil {
		return nil, err
	}

	return model, nil
}

// AdminGetPaymails get all block paymails filtered by conditions
func (wc *WalletClient) AdminGetPaymails(
	ctx context.Context,
	conditions *filter.AdminPaymailFilter,
	metadata map[string]any,
	queryParams *filter.QueryParams,
) ([]*models.PaymailAddress, error) {
	return Search[filter.AdminPaymailFilter, []*models.PaymailAddress](
		ctx, http.MethodPost,
		"/admin/paymails/search",
		wc.adminXPriv,
		conditions,
		metadata,
		queryParams,
		wc.doHTTPRequest,
	)
}

// AdminGetPaymailsCount get a count of all the paymails filtered by conditions
func (wc *WalletClient) AdminGetPaymailsCount(ctx context.Context, conditions *filter.AdminPaymailFilter, metadata map[string]any) (int64, error) {
	return Count(
		ctx, http.MethodPost,
		"/admin/paymails/count",
		wc.adminXPriv,
		conditions,
		metadata,
		wc.doHTTPRequest,
	)
}

// AdminCreatePaymail create a new paymail for a xpub
func (wc *WalletClient) AdminCreatePaymail(ctx context.Context, rawXPub string, address string, publicName string, avatar string) (*models.PaymailAddress, error) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldXpubKey:    rawXPub,
		FieldAddress:    address,
		FieldPublicName: publicName,
		FieldAvatar:     avatar,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var model *models.PaymailAddress
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, "/admin/paymail/create", jsonStr, wc.adminXPriv, true, &model,
	); err != nil {
		return nil, err
	}

	return model, nil
}

// AdminDeletePaymail delete a paymail address from the database
func (wc *WalletClient) AdminDeletePaymail(ctx context.Context, address string) error {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldAddress: address,
	})
	if err != nil {
		return WrapError(err)
	}

	if err := wc.doHTTPRequest(
		ctx, http.MethodDelete, "/admin/paymail/delete", jsonStr, wc.adminXPriv, true, nil,
	); err != nil {
		return err
	}

	return nil
}

// AdminGetTransactions get all block transactions filtered by conditions
func (wc *WalletClient) AdminGetTransactions(
	ctx context.Context,
	conditions *filter.TransactionFilter,
	metadata map[string]any,
	queryParams *filter.QueryParams,
) ([]*models.Transaction, error) {
	return Search[filter.TransactionFilter, []*models.Transaction](
		ctx, http.MethodPost,
		"/admin/transactions/search",
		wc.adminXPriv,
		conditions,
		metadata,
		queryParams,
		wc.doHTTPRequest,
	)
}

// AdminGetTransactionsCount get a count of all the transactions filtered by conditions
func (wc *WalletClient) AdminGetTransactionsCount(
	ctx context.Context,
	conditions *filter.TransactionFilter,
	metadata map[string]any,
) (int64, error) {
	return Count[filter.TransactionFilter](
		ctx, http.MethodPost,
		"/admin/transactions/count",
		wc.adminXPriv,
		conditions,
		metadata,
		wc.doHTTPRequest,
	)
}

// AdminGetUtxos get all block utxos filtered by conditions
func (wc *WalletClient) AdminGetUtxos(
	ctx context.Context,
	conditions *filter.AdminUtxoFilter,
	metadata map[string]any,
	queryParams *filter.QueryParams,
) ([]*models.Utxo, error) {
	return Search[filter.AdminUtxoFilter, []*models.Utxo](
		ctx, http.MethodPost,
		"/admin/utxos/search",
		wc.adminXPriv,
		conditions,
		metadata,
		queryParams,
		wc.doHTTPRequest,
	)
}

// AdminGetUtxosCount get a count of all the utxos filtered by conditions
func (wc *WalletClient) AdminGetUtxosCount(
	ctx context.Context,
	conditions *filter.AdminUtxoFilter,
	metadata map[string]any,
) (int64, error) {
	return Count[filter.AdminUtxoFilter](
		ctx, http.MethodPost,
		"/admin/utxos/count",
		wc.adminXPriv,
		conditions,
		metadata,
		wc.doHTTPRequest,
	)
}

// AdminGetXPubs get all block xpubs filtered by conditions
func (wc *WalletClient) AdminGetXPubs(ctx context.Context, conditions *filter.XpubFilter,
	metadata map[string]any, queryParams *filter.QueryParams,
) ([]*models.Xpub, error) {
	return Search[filter.XpubFilter, []*models.Xpub](
		ctx, http.MethodPost,
		"/admin/xpubs/search",
		wc.adminXPriv,
		conditions,
		metadata,
		queryParams,
		wc.doHTTPRequest,
	)
}

// AdminGetXPubsCount get a count of all the xpubs filtered by conditions
func (wc *WalletClient) AdminGetXPubsCount(
	ctx context.Context,
	conditions *filter.XpubFilter,
	metadata map[string]any,
) (int64, error) {
	return Count[filter.XpubFilter](
		ctx, http.MethodPost,
		"/admin/xpubs/count",
		wc.adminXPriv,
		conditions,
		metadata,
		wc.doHTTPRequest,
	)
}

func (wc *WalletClient) adminGetModels(
	ctx context.Context,
	conditions map[string]interface{},
	metadata map[string]any,
	queryParams *filter.QueryParams,
	path string,
	models interface{},
) error {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldConditions:  conditions,
		FieldMetadata:    metadata,
		FieldQueryParams: queryParams,
	})
	if err != nil {
		return WrapError(err)
	}

	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, path, jsonStr, wc.adminXPriv, true, &models,
	); err != nil {
		return err
	}

	return nil
}

func (wc *WalletClient) adminCount(ctx context.Context, conditions map[string]interface{}, metadata map[string]any, path string) (int64, error) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldConditions: conditions,
		FieldMetadata:   metadata,
	})
	if err != nil {
		return 0, WrapError(err)
	}

	var count int64
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, path, jsonStr, wc.adminXPriv, true, &count,
	); err != nil {
		return 0, err
	}

	return count, nil
}

// AdminRecordTransaction will record a transaction as an admin
func (wc *WalletClient) AdminRecordTransaction(ctx context.Context, hex string) (*models.Transaction, error) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldHex: hex,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var transaction models.Transaction
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, "/admin/transactions/record", jsonStr, wc.adminXPriv, wc.signRequest, &transaction,
	); err != nil {
		return nil, err
	}

	return &transaction, nil
}

// AdminGetContacts executes an HTTP POST request to search for contacts based on specified conditions, metadata, and query parameters.
func (wc *WalletClient) AdminGetContacts(ctx context.Context, conditions *filter.ContactFilter, metadata map[string]any, queryParams *filter.QueryParams) (*models.SearchContactsResponse, error) {
	return Search[filter.ContactFilter, *models.SearchContactsResponse](
		ctx, http.MethodPost,
		"/admin/contact/search",
		wc.adminXPriv,
		conditions,
		metadata,
		queryParams,
		wc.doHTTPRequest,
	)
}

// AdminUpdateContact executes an HTTP PATCH request to update a specific contact's full name using their ID.
func (wc *WalletClient) AdminUpdateContact(ctx context.Context, id, fullName string, metadata map[string]any) (*models.Contact, error) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		"fullName":    fullName,
		FieldMetadata: metadata,
	})
	if err != nil {
		return nil, WrapError(err)
	}
	var contact models.Contact
	err = wc.doHTTPRequest(ctx, http.MethodPatch, fmt.Sprintf("/admin/contact/%s", id), jsonStr, wc.adminXPriv, true, &contact)
	return &contact, WrapError(err)
}

// AdminDeleteContact executes an HTTP DELETE request to remove a contact using their ID.
func (wc *WalletClient) AdminDeleteContact(ctx context.Context, id string) error {
	err := wc.doHTTPRequest(ctx, http.MethodDelete, fmt.Sprintf("/admin/contact/%s", id), nil, wc.adminXPriv, true, nil)
	return WrapError(err)
}

// AdminAcceptContact executes an HTTP PATCH request to mark a contact as accepted using their ID.
func (wc *WalletClient) AdminAcceptContact(ctx context.Context, id string) (*models.Contact, error) {
	var contact models.Contact
	err := wc.doHTTPRequest(ctx, http.MethodPatch, fmt.Sprintf("/admin/contact/accepted/%s", id), nil, wc.adminXPriv, true, &contact)
	return &contact, WrapError(err)
}

// AdminRejectContact executes an HTTP PATCH request to mark a contact as rejected using their ID.
func (wc *WalletClient) AdminRejectContact(ctx context.Context, id string) (*models.Contact, error) {
	var contact models.Contact
	err := wc.doHTTPRequest(ctx, http.MethodPatch, fmt.Sprintf("/admin/contact/rejected/%s", id), nil, wc.adminXPriv, true, &contact)
	return &contact, WrapError(err)
}

// FinalizeTransaction will finalize the transaction
func (wc *WalletClient) FinalizeTransaction(draft *models.DraftTransaction) (string, error) {
	res, err := GetSignedHex(draft, wc.xPriv)
	if err != nil {
		return "", WrapError(err)
	}

	return res, nil
}

// SendToRecipients send to recipients
func (wc *WalletClient) SendToRecipients(ctx context.Context, recipients []*Recipients, metadata map[string]any) (*models.Transaction, error) {
	draft, err := wc.DraftToRecipients(ctx, recipients, metadata)
	if err != nil {
		return nil, err
	}

	var hex string
	if hex, err = wc.FinalizeTransaction(draft); err != nil {
		return nil, err
	}

	return wc.RecordTransaction(ctx, hex, draft.ID, metadata)
}

// AdminSubscribeWebhook subscribes to a webhook to receive notifications from spv-wallet
func (wc *WalletClient) AdminSubscribeWebhook(ctx context.Context, webhookURL, tokenHeader, tokenValue string) error {
	requestModel := models.SubscribeRequestBody{
		URL:         webhookURL,
		TokenHeader: tokenHeader,
		TokenValue:  tokenValue,
	}
	rawJSON, err := json.Marshal(requestModel)
	if err != nil {
		return WrapError(err)
	}
	err = wc.doHTTPRequest(ctx, http.MethodPost, "/admin/webhooks/subscriptions", rawJSON, wc.adminXPriv, true, nil)
	return WrapError(err)
}

// AdminUnsubscribeWebhook unsubscribes from a webhook
func (wc *WalletClient) AdminUnsubscribeWebhook(ctx context.Context, webhookURL string) error {
	requestModel := models.UnsubscribeRequestBody{
		URL: webhookURL,
	}
	rawJSON, err := json.Marshal(requestModel)
	if err != nil {
		return WrapError(err)
	}
	err = wc.doHTTPRequest(ctx, http.MethodDelete, "/admin/webhooks/subscriptions", rawJSON, wc.adminXPriv, true, nil)
	return err
}
