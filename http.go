package walletclient

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet-go-client/utils"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

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
		ctx, http.MethodGet, "/users/current", nil, wc.xPriv, true, &xPub,
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
		ctx, http.MethodPatch, "/users/current", jsonStr, wc.xPriv, true, &xPub,
	); err != nil {
		return nil, err
	}

	return &xPub, nil
}

// GetAccessKey will get an access key by id
func (wc *WalletClient) GetAccessKey(ctx context.Context, id string) (*models.AccessKey, error) {
	var accessKey models.AccessKey
	if err := wc.doHTTPRequest(
		ctx, http.MethodGet, "/users/current/keys/"+id, nil, wc.xPriv, true, &accessKey,
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
		"/users/current",
		wc.xPriv,
		conditions,
		metadata,
		queryParams,
		wc.doHTTPRequest,
	)
}

// RevokeAccessKey will revoke an access key by id
func (wc *WalletClient) RevokeAccessKey(ctx context.Context, id string) (*models.AccessKey, error) {
	var accessKey models.AccessKey
	if err := wc.doHTTPRequest(
		ctx, http.MethodDelete, "/users/current/keys/"+id, nil, wc.xPriv, true, &accessKey,
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
		ctx, http.MethodPost, "/users/current/keys", jsonStr, wc.xPriv, true, &accessKey,
	); err != nil {
		return nil, err
	}

	return &accessKey, nil
}

// GetTransaction will get a transaction by ID
func (wc *WalletClient) GetTransaction(ctx context.Context, txID string) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := wc.doHTTPRequest(ctx, http.MethodGet, "/transactions/"+FieldID+"="+txID, nil, wc.xPriv, wc.signRequest, &transaction); err != nil {
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
		"/transactions",
		wc.xPriv,
		conditions,
		metadata,
		queryParams,
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
		ctx, http.MethodPost, "/transactions", jsonStr, wc.xPriv, true, &draftTransaction,
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
		ctx, http.MethodPost, "/transactions", jsonStr, wc.xPriv, wc.signRequest, &transaction,
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
		ctx, http.MethodPatch, "/transactions", jsonStr, wc.xPriv, wc.signRequest, &transaction,
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

// GetUtxos will get a list of utxos filtered by conditions and metadata
func (wc *WalletClient) GetUtxos(ctx context.Context, conditions *filter.UtxoFilter, metadata map[string]any, queryParams *filter.QueryParams) ([]*models.Utxo, error) {
	return Search[filter.UtxoFilter, []*models.Utxo](
		ctx, http.MethodPost,
		"/utxos",
		wc.xPriv,
		conditions,
		metadata,
		queryParams,
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

	var privateKey *ec.PrivateKey
	if privateKey, err = ec.PrivateKeyFromHex(
		privateKeyHex,
	); err != nil {
		return
	}
	publicKey := privateKey.PubKey()

	// Get the AccessKey
	payload = new(models.AuthPayload)
	payload.AccessKey = hex.EncodeToString(publicKey.SerializeCompressed())

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
		xPub, err := bip32.GetExtendedPublicKey(xPriv)
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
		return ErrMissingAccessKey
	}
	return SetSignatureFromAccessKey(&req.Header, hex.EncodeToString(wc.accessKey.Serialize()), string(rawJSON))
}

// AcceptContact will accept the contact associated with the paymail
func (wc *WalletClient) AcceptContact(ctx context.Context, paymail string) error {
	if err := wc.doHTTPRequest(
		ctx, http.MethodPatch, "/contacts/"+paymail, nil, wc.xPriv, wc.signRequest, nil,
	); err != nil {
		return err
	}

	return nil
}

// RejectContact will reject the contact associated with the paymail
func (wc *WalletClient) RejectContact(ctx context.Context, paymail string) error {
	if err := wc.doHTTPRequest(
		ctx, http.MethodPatch, "/contacts/"+paymail, nil, wc.xPriv, wc.signRequest, nil,
	); err != nil {
		return err
	}

	return nil
}

// ConfirmContact will confirm the contact associated with the paymail
func (wc *WalletClient) ConfirmContact(ctx context.Context, contact *models.Contact, passcode, requesterPaymail string, period, digits uint) error {
	isTotpValid, err := wc.ValidateTotpForContact(contact, passcode, requesterPaymail, period, digits)
	if err != nil {
		return WrapError(ErrTotpInvalid)
	}

	if !isTotpValid {
		return WrapError(ErrTotpInvalid)
	}

	if err := wc.doHTTPRequest(
		ctx, http.MethodPatch, "/contacts/"+contact.Paymail, nil, wc.xPriv, wc.signRequest, nil,
	); err != nil {
		return err
	}

	return nil
}

// GetContacts will get contacts by conditions
func (wc *WalletClient) GetContacts(ctx context.Context, conditions *filter.ContactFilter, metadata map[string]any, queryParams *filter.QueryParams) (*models.SearchContactsResponse, error) {
	return Search[filter.ContactFilter, *models.SearchContactsResponse](
		ctx, http.MethodPost,
		"/contacts/",
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
		ctx, http.MethodPut, "/contacts/"+paymail, jsonStr, wc.xPriv, wc.signRequest, &result,
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
		return nil, WrapError(ErrMissingKey)
	}
	if err := wc.doHTTPRequest(
		ctx, http.MethodGet, "/configs/shared", nil, key, true, &model,
	); err != nil {
		return nil, err
	}

	return model, nil
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
