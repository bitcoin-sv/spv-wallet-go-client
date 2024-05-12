package walletclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/models"
)

// AdminNewXpub will register an xPub
func (wc *WalletClient) AdminNewXpub(ctx context.Context, rawXPub string, metadata *models.Metadata) ResponseError {
	// Adding a xpub needs to be signed by an admin key
	if wc.adminXPriv == nil {
		return WrapError(ErrAdminKey)
	}

	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: processMetadata(metadata),
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
func (wc *WalletClient) AdminGetStatus(ctx context.Context) (bool, ResponseError) {
	var status bool
	if err := wc.doHTTPRequest(
		ctx, http.MethodGet, "/admin/status", nil, wc.xPriv, true, &status,
	); err != nil {
		return false, err
	}

	return status, nil
}

// AdminGetStats get admin stats
func (wc *WalletClient) AdminGetStats(ctx context.Context) (*models.AdminStats, ResponseError) {
	var stats *models.AdminStats
	if err := wc.doHTTPRequest(
		ctx, http.MethodGet, "/admin/stats", nil, wc.xPriv, true, &stats,
	); err != nil {
		return nil, err
	}

	return stats, nil
}

// AdminGetAccessKeys get all access keys filtered by conditions
func (wc *WalletClient) AdminGetAccessKeys(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata, queryParams *QueryParams,
) ([]*models.AccessKey, ResponseError) {
	var models []*models.AccessKey
	if err := wc.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/access-keys/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetAccessKeysCount get a count of all the access keys filtered by conditions
func (wc *WalletClient) AdminGetAccessKeysCount(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata,
) (int64, ResponseError) {
	return wc.adminCount(ctx, conditions, metadata, "/admin/access-keys/count")
}

// AdminGetBlockHeaders get all block headers filtered by conditions
func (wc *WalletClient) AdminGetBlockHeaders(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata, queryParams *QueryParams,
) ([]*models.BlockHeader, ResponseError) {
	var models []*models.BlockHeader
	if err := wc.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/block-headers/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetBlockHeadersCount get a count of all the block headers filtered by conditions
func (wc *WalletClient) AdminGetBlockHeadersCount(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata,
) (int64, ResponseError) {
	return wc.adminCount(ctx, conditions, metadata, "/admin/block-headers/count")
}

// AdminGetDestinations get all block destinations filtered by conditions
func (wc *WalletClient) AdminGetDestinations(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata, queryParams *QueryParams,
) ([]*models.Destination, ResponseError) {
	var models []*models.Destination
	if err := wc.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/destinations/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetDestinationsCount get a count of all the destinations filtered by conditions
func (wc *WalletClient) AdminGetDestinationsCount(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata,
) (int64, ResponseError) {
	return wc.adminCount(ctx, conditions, metadata, "/admin/destinations/count")
}

// AdminGetPaymail get a paymail by address
func (wc *WalletClient) AdminGetPaymail(ctx context.Context, address string) (*models.PaymailAddress, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldAddress: address,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var model *models.PaymailAddress
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, "/admin/paymail/get", jsonStr, wc.xPriv, true, &model,
	); err != nil {
		return nil, err
	}

	return model, nil
}

// AdminGetPaymails get all block paymails filtered by conditions
func (wc *WalletClient) AdminGetPaymails(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata, queryParams *QueryParams,
) ([]*models.PaymailAddress, ResponseError) {
	var models []*models.PaymailAddress
	if err := wc.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/paymails/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetPaymailsCount get a count of all the paymails filtered by conditions
func (wc *WalletClient) AdminGetPaymailsCount(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata,
) (int64, ResponseError) {
	return wc.adminCount(ctx, conditions, metadata, "/admin/paymails/count")
}

// AdminCreatePaymail create a new paymail for a xpub
func (wc *WalletClient) AdminCreatePaymail(ctx context.Context, rawXPub string, address string, publicName string, avatar string) (*models.PaymailAddress, ResponseError) {
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
		ctx, http.MethodPost, "/admin/paymail/create", jsonStr, wc.xPriv, true, &model,
	); err != nil {
		return nil, err
	}

	return model, nil
}

// AdminDeletePaymail delete a paymail address from the database
func (wc *WalletClient) AdminDeletePaymail(ctx context.Context, address string) ResponseError {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldAddress: address,
	})
	if err != nil {
		return WrapError(err)
	}

	if err := wc.doHTTPRequest(
		ctx, http.MethodDelete, "/admin/paymail/delete", jsonStr, wc.xPriv, true, nil,
	); err != nil {
		return err
	}

	return nil
}

// AdminGetTransactions get all block transactions filtered by conditions
func (wc *WalletClient) AdminGetTransactions(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata, queryParams *QueryParams,
) ([]*models.Transaction, ResponseError) {
	var models []*models.Transaction
	if err := wc.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/transactions/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetTransactionsCount get a count of all the transactions filtered by conditions
func (wc *WalletClient) AdminGetTransactionsCount(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata,
) (int64, ResponseError) {
	return wc.adminCount(ctx, conditions, metadata, "/admin/transactions/count")
}

// AdminGetUtxos get all block utxos filtered by conditions
func (wc *WalletClient) AdminGetUtxos(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata, queryParams *QueryParams,
) ([]*models.Utxo, ResponseError) {
	var models []*models.Utxo
	if err := wc.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/utxos/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetUtxosCount get a count of all the utxos filtered by conditions
func (wc *WalletClient) AdminGetUtxosCount(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata,
) (int64, ResponseError) {
	return wc.adminCount(ctx, conditions, metadata, "/admin/utxos/count")
}

// AdminGetXPubs get all block xpubs filtered by conditions
func (wc *WalletClient) AdminGetXPubs(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata, queryParams *QueryParams,
) ([]*models.Xpub, ResponseError) {
	var models []*models.Xpub
	if err := wc.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/xpubs/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetXPubsCount get a count of all the xpubs filtered by conditions
func (wc *WalletClient) AdminGetXPubsCount(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata,
) (int64, ResponseError) {
	return wc.adminCount(ctx, conditions, metadata, "/admin/xpubs/count")
}

func (wc *WalletClient) adminGetModels(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata, queryParams *QueryParams, path string, models interface{},
) ResponseError {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldConditions:  conditions,
		FieldMetadata:    processMetadata(metadata),
		FieldQueryParams: queryParams,
	})
	if err != nil {
		return WrapError(err)
	}

	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, path, jsonStr, wc.xPriv, true, &models,
	); err != nil {
		return err
	}

	return nil
}

func (wc *WalletClient) adminCount(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata, path string) (int64, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldConditions: conditions,
		FieldMetadata:   processMetadata(metadata),
	})
	if err != nil {
		return 0, WrapError(err)
	}

	var count int64
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, path, jsonStr, wc.xPriv, true, &count,
	); err != nil {
		return 0, err
	}

	return count, nil
}

// AdminRecordTransaction will record a transaction as an admin
func (wc *WalletClient) AdminRecordTransaction(ctx context.Context, hex string) (*models.Transaction, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldHex: hex,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var transaction models.Transaction
	if err := wc.doHTTPRequest(
		ctx, http.MethodPost, "/admin/transactions/record", jsonStr, wc.xPriv, *wc.signRequest, &transaction,
	); err != nil {
		return nil, err
	}

	return &transaction, nil
}

// AdminGetSharedConfig gets the shared config
func (wc *WalletClient) AdminGetSharedConfig(ctx context.Context) (*models.SharedConfig, ResponseError) {
	var model *models.SharedConfig
	if err := wc.doHTTPRequest(
		ctx, http.MethodGet, "/admin/shared-config", nil, wc.xPriv, true, &model,
	); err != nil {
		return nil, err
	}

	return model, nil
}

// AdminGetContacts executes an HTTP POST request to search for contacts based on specified conditions, metadata, and query parameters.
func (wc *WalletClient) AdminGetContacts(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata, queryParams *QueryParams) ([]*models.Contact, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldConditions:  conditions,
		FieldMetadata:    processMetadata(metadata),
		FieldQueryParams: queryParams,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var contacts []*models.Contact
	err = wc.doHTTPRequest(ctx, http.MethodPost, "/admin/contact/search", jsonStr, wc.adminXPriv, true, &contacts)
	return contacts, WrapError(err)
}

// AdminUpdateContact executes an HTTP PATCH request to update a specific contact's full name using their ID.
func (wc *WalletClient) AdminUpdateContact(ctx context.Context, id, fullName string, metadata *models.Metadata) (*models.Contact, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		"fullName":    fullName,
		FieldMetadata: processMetadata(metadata),
	})
	if err != nil {
		return nil, WrapError(err)
	}
	var contact models.Contact
	err = wc.doHTTPRequest(ctx, http.MethodPatch, fmt.Sprintf("/admin/contact/%s", id), jsonStr, wc.adminXPriv, true, &contact)
	return &contact, WrapError(err)
}

// AdminDeleteContact executes an HTTP DELETE request to remove a contact using their ID.
func (wc *WalletClient) AdminDeleteContact(ctx context.Context, id string) ResponseError {
	err := wc.doHTTPRequest(ctx, http.MethodDelete, fmt.Sprintf("/admin/contact/%s", id), nil, wc.adminXPriv, true, nil)
	return WrapError(err)
}

// AdminAcceptContact executes an HTTP PATCH request to mark a contact as accepted using their ID.
func (wc *WalletClient) AdminAcceptContact(ctx context.Context, id string) (*models.Contact, ResponseError) {
	var contact models.Contact
	err := wc.doHTTPRequest(ctx, http.MethodPatch, fmt.Sprintf("/admin/contact/accepted/%s", id), nil, wc.adminXPriv, true, &contact)
	return &contact, WrapError(err)
}

// AdminRejectContact executes an HTTP PATCH request to mark a contact as rejected using their ID.
func (wc *WalletClient) AdminRejectContact(ctx context.Context, id string) (*models.Contact, ResponseError) {
	var contact models.Contact
	err := wc.doHTTPRequest(ctx, http.MethodPatch, fmt.Sprintf("/admin/contact/rejected/%s", id), nil, wc.adminXPriv, true, &contact)
	return &contact, WrapError(err)
}
