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

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/apierrors"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"

	"github.com/bitcoin-sv/spv-wallet-go-client/transports"
	"github.com/bitcoin-sv/spv-wallet-go-client/utils"
)

type transportHTTP struct {
	accessKey   *bec.PrivateKey
	adminXPriv  *bip32.ExtendedKey
	httpClient  *http.Client
	server      string
	signRequest bool
	xPriv       *bip32.ExtendedKey
	xPub        *bip32.ExtendedKey
}

// SetSignRequest turn the signing of the http request on or off
func (wc *WalletClient) SetSignRequest(signRequest bool) {
	wc.signRequest = &signRequest
}

// IsSignRequest return whether to sign all requests
func (wc *WalletClient) IsSignRequest() bool {
	return *wc.signRequest
}

// SetAdminKey set the admin key
func (wc *WalletClient) SetAdminKey(adminKey *bip32.ExtendedKey) {
	wc.adminXPriv = adminKey
}

// GetXPub will get the xpub of the current xpub
func (wc *WalletClient) GetXPub(ctx context.Context) (*models.Xpub, ResponseError) {
	var xPub models.Xpub
	if err := wc.doHTTPRequest(
		ctx, http.MethodGet, "/xpub", nil, wc.xPriv, true, &xPub,
	); err != nil {
		return nil, err
	}

	return &xPub, nil
}

// UpdateXPubMetadata update the metadata of the logged in xpub
func (wc *WalletClient) UpdateXPubMetadata(ctx context.Context, metadata *models.Metadata) (*models.Xpub, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: processMetadata(metadata),
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
func (wc *WalletClient) GetAccessKey(ctx context.Context, id string) (*models.AccessKey, ResponseError) {
	var accessKey models.AccessKey
	if err := wc.doHTTPRequest(
		ctx, http.MethodGet, "/access-key?"+FieldID+"="+id, nil, wc.xPriv, true, &accessKey,
	); err != nil {
		return nil, err
	}

	return &accessKey, nil
}

// GetAccessKeys will get all access keys matching the metadata filter
func (wc *WalletClient) GetAccessKeys(ctx context.Context, metadataConditions *models.Metadata) ([]*models.AccessKey, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: processMetadata(metadataConditions),
	})
	if err != nil {
		return nil, WrapError(err)
	}
	var accessKey []*models.AccessKey
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, "/access-key/search", jsonStr, wc.xPriv, true, &accessKey,
	); err != nil {
		return nil, err
	}

	return accessKey, nil
}

// GetAccessKeysCount will get the count of access keys
func (wc *WalletClient) GetAccessKeysCount(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata) (int64, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldConditions: conditions,
		FieldMetadata:   processMetadata(metadata),
	})
	if err != nil {
		return 0, WrapError(err)
	}

	var count int64
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, "/access-key/count", jsonStr, wc.xPriv, true, &count,
	); err != nil {
		return 0, err
	}

	return count, nil
}

// RevokeAccessKey will revoke an access key by id
func (wc *WalletClient) RevokeAccessKey(ctx context.Context, id string) (*models.AccessKey, ResponseError) {
	var accessKey models.AccessKey
	if err := wc.doHTTPRequest(
		ctx, http.MethodDelete, "/access-key?"+FieldID+"="+id, nil, wc.xPriv, true, &accessKey,
	); err != nil {
		return nil, err
	}

	return &accessKey, nil
}

// CreateAccessKey will create new access key
func (wc *WalletClient) CreateAccessKey(ctx context.Context, metadata *models.Metadata) (*models.AccessKey, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: processMetadata(metadata),
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
func (wc *WalletClient) GetDestinationByID(ctx context.Context, id string) (*models.Destination, ResponseError) {
	var destination models.Destination
	if err := wc.doHTTPRequest(
		ctx, http.MethodGet, "/destination?"+FieldID+"="+id, nil, wc.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}

	return &destination, nil
}

// GetDestinationByAddress will get a destination by address
func (wc *WalletClient) GetDestinationByAddress(ctx context.Context, address string) (*models.Destination, ResponseError) {
	var destination models.Destination
	if err := wc.doHTTPRequest(
		ctx, http.MethodGet, "/destination?"+FieldAddress+"="+address, nil, wc.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}

	return &destination, nil
}

// GetDestinationByLockingScript will get a destination by locking script
func (wc *WalletClient) GetDestinationByLockingScript(ctx context.Context, lockingScript string) (*models.Destination, ResponseError) {
	var destination models.Destination
	if err := wc.doHTTPRequest(
		ctx, http.MethodGet, "/destination?"+FieldLockingScript+"="+lockingScript, nil, wc.xPriv, true, &destination,
	); err != nil {
		return nil, err
	}

	return &destination, nil
}

// GetDestinations will get all destinations matching the metadata filter
func (wc *WalletClient) GetDestinations(ctx context.Context, metadataConditions *models.Metadata) ([]*models.Destination, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: processMetadata(metadataConditions),
	})
	if err != nil {
		return nil, WrapError(err)
	}
	var destinations []*models.Destination
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, "/destination/search", jsonStr, wc.xPriv, true, &destinations,
	); err != nil {
		return nil, err
	}

	return destinations, nil
}

// GetDestinationsCount will get the count of destinations matching the metadata filter
func (wc *WalletClient) GetDestinationsCount(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata) (int64, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldConditions: conditions,
		FieldMetadata:   processMetadata(metadata),
	})
	if err != nil {
		return 0, WrapError(err)
	}

	var count int64
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, "/destination/count", jsonStr, wc.xPriv, true, &count,
	); err != nil {
		return 0, err
	}

	return count, nil
}

// NewDestination will create a new destination and return it
func (wc *WalletClient) NewDestination(ctx context.Context, metadata *models.Metadata) (*models.Destination, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: processMetadata(metadata),
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
func (wc *WalletClient) UpdateDestinationMetadataByID(ctx context.Context, id string,
	metadata *models.Metadata,
) (*models.Destination, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldID:       id,
		FieldMetadata: processMetadata(metadata),
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
func (wc *WalletClient) UpdateDestinationMetadataByAddress(ctx context.Context, address string,
	metadata *models.Metadata,
) (*models.Destination, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldAddress:  address,
		FieldMetadata: processMetadata(metadata),
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
func (wc *WalletClient) UpdateDestinationMetadataByLockingScript(ctx context.Context, lockingScript string,
	metadata *models.Metadata,
) (*models.Destination, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldLockingScript: lockingScript,
		FieldMetadata:      processMetadata(metadata),
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

// GetTransactions will get transactions by conditions
func (wc *WalletClient) GetTransactions(ctx context.Context, conditions map[string]interface{},
	metadataConditions *models.Metadata, queryParams *QueryParams,
) ([]*models.Transaction, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldConditions:  conditions,
		FieldMetadata:    processMetadata(metadataConditions),
		FieldQueryParams: queryParams,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var transactions []*models.Transaction
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, "/transaction/search", jsonStr, wc.xPriv, *wc.signRequest, &transactions,
	); err != nil {
		return nil, err
	}

	return transactions, nil
}

// GetTransactionsCount get number of user transactions
func (wc *WalletClient) GetTransactionsCount(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata,
) (int64, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldConditions: conditions,
		FieldMetadata:   processMetadata(metadata),
	})
	if err != nil {
		return 0, WrapError(err)
	}

	var count int64
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, "/transaction/count", jsonStr, wc.xPriv, *wc.signRequest, &count,
	); err != nil {
		return 0, err
	}

	return count, nil
}

// DraftToRecipients is a draft transaction to a slice of recipients
func (wc *WalletClient) DraftToRecipients(ctx context.Context, recipients []*Recipients,
	metadata *models.Metadata,
) (*models.DraftTransaction, ResponseError) {
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
		FieldMetadata: processMetadata(metadata),
	})
}

// DraftTransaction is a draft transaction
func (wc *WalletClient) DraftTransaction(ctx context.Context, transactionConfig *models.TransactionConfig,
	metadata *models.Metadata,
) (*models.DraftTransaction, ResponseError) {
	return wc.createDraftTransaction(ctx, map[string]interface{}{
		FieldConfig:   transactionConfig,
		FieldMetadata: processMetadata(metadata),
	})
}

// createDraftTransaction will create a draft transaction
func (wc *WalletClient) createDraftTransaction(ctx context.Context,
	jsonData map[string]interface{},
) (*models.DraftTransaction, ResponseError) {
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
		return nil, WrapError(apierrors.ErrDraftNotFound)
	}

	return draftTransaction, nil
}

// RecordTransaction will record a transaction
func (wc *WalletClient) RecordTransaction(ctx context.Context, hex, referenceID string,
	metadata *models.Metadata,
) (*models.Transaction, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldHex:         hex,
		FieldReferenceID: referenceID,
		FieldMetadata:    processMetadata(metadata),
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var transaction models.Transaction
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, "/transaction/record", jsonStr, wc.xPriv, *wc.signRequest, &transaction,
	); err != nil {
		return nil, err
	}

	return &transaction, nil
}

// UpdateTransactionMetadata update the metadata of a transaction
func (wc *WalletClient) UpdateTransactionMetadata(ctx context.Context, txID string,
	metadata *models.Metadata,
) (*models.Transaction, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldID:       txID,
		FieldMetadata: processMetadata(metadata),
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var transaction models.Transaction
	if err := wc.doHTTPRequest(
		ctx, http.MethodPatch, "/transaction", jsonStr, wc.xPriv, *wc.signRequest, &transaction,
	); err != nil {
		return nil, err
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
	header.Set(models.AuthAccessKey, authData.AccessKey)

	return setSignatureHeaders(header, authData)
}

// GetUtxo will get a utxo by transaction ID
func (wc *WalletClient) GetUtxo(ctx context.Context, txID string, outputIndex uint32) (*models.Utxo, ResponseError) {
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
func (wc *WalletClient) GetUtxos(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata, queryParams *QueryParams) ([]*models.Utxo, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldConditions:  conditions,
		FieldMetadata:    processMetadata(metadata),
		FieldQueryParams: queryParams,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var utxos []*models.Utxo
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, "/utxo/search", jsonStr, wc.xPriv, *wc.signRequest, &utxos,
	); err != nil {
		return nil, err
	}

	return utxos, nil
}

// GetUtxosCount will get the count of utxos filtered by conditions and metadata
func (wc *WalletClient) GetUtxosCount(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata) (int64, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldConditions: conditions,
		FieldMetadata:   processMetadata(metadata),
	})
	if err != nil {
		return 0, WrapError(err)
	}

	var count int64
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, "/utxo/count", jsonStr, wc.xPriv, *wc.signRequest, &count,
	); err != nil {
		return 0, err
	}

	return count, nil
}

// createSignatureAccessKey will create a signature for the given access key & body contents
func createSignatureAccessKey(privateKeyHex, bodyString string) (payload *models.AuthPayload, err error) {
	// No key?
	if privateKeyHex == "" {
		err = apierrors.ErrMissingAccessKey
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
) ResponseError {
	req, err := http.NewRequestWithContext(ctx, method, *wc.server+path, bytes.NewBuffer(rawJSON))
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

func (wc *WalletClient) authenticateWithXpriv(sign bool, req *http.Request, xPriv *bip32.ExtendedKey, rawJSON []byte) ResponseError {
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

func (wc *WalletClient) authenticateWithAccessKey(req *http.Request, rawJSON []byte) ResponseError {
	return SetSignatureFromAccessKey(&req.Header, hex.EncodeToString(wc.accessKey.Serialise()), string(rawJSON))
}

// AcceptContact will accept the contact associated with the paymail
func (wc *WalletClient) AcceptContact(ctx context.Context, paymail string) ResponseError {
	if err := wc.doHTTPRequest(
		ctx, http.MethodPatch, "/contact/accepted/"+paymail, nil, wc.xPriv, *wc.signRequest, nil,
	); err != nil {
		return err
	}

	return nil
}

// RejectContact will reject the contact associated with the paymail
func (wc *WalletClient) RejectContact(ctx context.Context, paymail string) ResponseError {
	if err := wc.doHTTPRequest(
		ctx, http.MethodPatch, "/contact/rejected/"+paymail, nil, wc.xPriv, *wc.signRequest, nil,
	); err != nil {
		return err
	}

	return nil
}

// ConfirmContact will confirm the contact associated with the paymail
func (wc *WalletClient) ConfirmContact(ctx context.Context, contact *models.Contact, passcode, requesterPaymail string, period, digits uint) ResponseError {
	isTotpValid, err := wc.ValidateTotpForContact(contact, passcode, requesterPaymail, period, digits)
	if err != nil {
		return WrapError(fmt.Errorf("totp validation failed: %w", err))
	}

	if !isTotpValid {
		return WrapError(errors.New("totp is invalid"))
	}

	if err := wc.doHTTPRequest(
		ctx, http.MethodPatch, "/contact/confirmed/"+contact.Paymail, nil, wc.xPriv, *wc.signRequest, nil,
	); err != nil {
		return err
	}

	return nil
}

// GetContacts will get contacts by conditions
func (wc *WalletClient) GetContacts(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata, queryParams *QueryParams) ([]*models.Contact, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldConditions:  conditions,
		FieldMetadata:    processMetadata(metadata),
		FieldQueryParams: queryParams,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var result []*models.Contact
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, "/contact/search", jsonStr, wc.xPriv, *wc.signRequest, &result,
	); err != nil {
		return nil, err
	}

	return result, nil
}

// UpsertContact add or update contact. When adding a new contact, the system utilizes Paymail's PIKE capability to dispatch an invitation request, asking the counterparty to include the current user in their contacts.
func (wc *WalletClient) UpsertContact(ctx context.Context, paymail, fullName string, metadata *models.Metadata) (*models.Contact, ResponseError) {
	return wc.UpsertContactForPaymail(ctx, paymail, fullName, metadata, "")
}

func (wc *WalletClient) UpsertContactForPaymail(ctx context.Context, paymail, fullName string, metadata *models.Metadata, requesterPaymail string) (*models.Contact, transports.ResponseError) {
	payload := map[string]interface{}{
		"fullName":    fullName,
		FieldMetadata: processMetadata(metadata),
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
		ctx, http.MethodPut, "/contact/"+paymail, jsonStr, wc.xPriv, *wc.signRequest, &result,
	); err != nil {
		return nil, err
	}

	return &result, nil
}
