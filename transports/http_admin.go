package transports

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/models"
)

// AdminNewXpub will register an xPub
func (h *TransportHTTP) AdminNewXpub(ctx context.Context, rawXPub string, metadata *models.Metadata) ResponseError {
	// Adding a xpub needs to be signed by an admin key
	if h.adminXPriv == nil {
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

	return h.doHTTPRequest(
		ctx, http.MethodPost, "/admin/xpub", jsonStr, h.adminXPriv, true, &xPubData,
	)
}

// AdminGetStatus get whether admin key is valid
func (h *TransportHTTP) AdminGetStatus(ctx context.Context) (bool, ResponseError) {
	var status bool
	if err := h.doHTTPRequest(
		ctx, http.MethodGet, "/admin/status", nil, h.xPriv, true, &status,
	); err != nil {
		return false, err
	}

	return status, nil
}

// AdminGetStats get admin stats
func (h *TransportHTTP) AdminGetStats(ctx context.Context) (*models.AdminStats, ResponseError) {
	var stats *models.AdminStats
	if err := h.doHTTPRequest(
		ctx, http.MethodGet, "/admin/stats", nil, h.xPriv, true, &stats,
	); err != nil {
		return nil, err
	}

	return stats, nil
}

// AdminGetAccessKeys get all access keys filtered by conditions
func (h *TransportHTTP) AdminGetAccessKeys(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata, queryParams *QueryParams,
) ([]*models.AccessKey, ResponseError) {
	var models []*models.AccessKey
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/access-keys/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetAccessKeysCount get a count of all the access keys filtered by conditions
func (h *TransportHTTP) AdminGetAccessKeysCount(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata,
) (int64, ResponseError) {
	return h.adminCount(ctx, conditions, metadata, "/admin/access-keys/count")
}

// AdminGetBlockHeaders get all block headers filtered by conditions
func (h *TransportHTTP) AdminGetBlockHeaders(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata, queryParams *QueryParams,
) ([]*models.BlockHeader, ResponseError) {
	var models []*models.BlockHeader
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/block-headers/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetBlockHeadersCount get a count of all the block headers filtered by conditions
func (h *TransportHTTP) AdminGetBlockHeadersCount(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata,
) (int64, ResponseError) {
	return h.adminCount(ctx, conditions, metadata, "/admin/block-headers/count")
}

// AdminGetDestinations get all block destinations filtered by conditions
func (h *TransportHTTP) AdminGetDestinations(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata, queryParams *QueryParams,
) ([]*models.Destination, ResponseError) {
	var models []*models.Destination
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/destinations/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetDestinationsCount get a count of all the destinations filtered by conditions
func (h *TransportHTTP) AdminGetDestinationsCount(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata,
) (int64, ResponseError) {
	return h.adminCount(ctx, conditions, metadata, "/admin/destinations/count")
}

// AdminGetPaymail get a paymail by address
func (h *TransportHTTP) AdminGetPaymail(ctx context.Context, address string) (*models.PaymailAddress, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldAddress: address,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var model *models.PaymailAddress
	if err := h.doHTTPRequest(
		ctx, http.MethodPost, "/admin/paymail/get", jsonStr, h.xPriv, true, &model,
	); err != nil {
		return nil, err
	}

	return model, nil
}

// AdminGetPaymails get all block paymails filtered by conditions
func (h *TransportHTTP) AdminGetPaymails(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata, queryParams *QueryParams,
) ([]*models.PaymailAddress, ResponseError) {
	var models []*models.PaymailAddress
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/paymails/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetPaymailsCount get a count of all the paymails filtered by conditions
func (h *TransportHTTP) AdminGetPaymailsCount(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata,
) (int64, ResponseError) {
	return h.adminCount(ctx, conditions, metadata, "/admin/paymails/count")
}

// AdminCreatePaymail create a new paymail for a xpub
func (h *TransportHTTP) AdminCreatePaymail(ctx context.Context, rawXPub string, address string, publicName string, avatar string) (*models.PaymailAddress, ResponseError) {
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
	if err := h.doHTTPRequest(
		ctx, http.MethodPost, "/admin/paymail/create", jsonStr, h.xPriv, true, &model,
	); err != nil {
		return nil, err
	}

	return model, nil
}

// AdminDeletePaymail delete a paymail address from the database
func (h *TransportHTTP) AdminDeletePaymail(ctx context.Context, address string) ResponseError {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldAddress: address,
	})
	if err != nil {
		return WrapError(err)
	}

	if err := h.doHTTPRequest(
		ctx, http.MethodDelete, "/admin/paymail/delete", jsonStr, h.xPriv, true, nil,
	); err != nil {
		return err
	}

	return nil
}

// AdminGetTransactions get all block transactions filtered by conditions
func (h *TransportHTTP) AdminGetTransactions(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata, queryParams *QueryParams,
) ([]*models.Transaction, ResponseError) {
	var models []*models.Transaction
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/transactions/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetTransactionsCount get a count of all the transactions filtered by conditions
func (h *TransportHTTP) AdminGetTransactionsCount(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata,
) (int64, ResponseError) {
	return h.adminCount(ctx, conditions, metadata, "/admin/transactions/count")
}

// AdminGetUtxos get all block utxos filtered by conditions
func (h *TransportHTTP) AdminGetUtxos(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata, queryParams *QueryParams,
) ([]*models.Utxo, ResponseError) {
	var models []*models.Utxo
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/utxos/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetUtxosCount get a count of all the utxos filtered by conditions
func (h *TransportHTTP) AdminGetUtxosCount(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata,
) (int64, ResponseError) {
	return h.adminCount(ctx, conditions, metadata, "/admin/utxos/count")
}

// AdminGetXPubs get all block xpubs filtered by conditions
func (h *TransportHTTP) AdminGetXPubs(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata, queryParams *QueryParams,
) ([]*models.Xpub, ResponseError) {
	var models []*models.Xpub
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/xpubs/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetXPubsCount get a count of all the xpubs filtered by conditions
func (h *TransportHTTP) AdminGetXPubsCount(ctx context.Context, conditions map[string]interface{},
	metadata *models.Metadata,
) (int64, ResponseError) {
	return h.adminCount(ctx, conditions, metadata, "/admin/xpubs/count")
}

func (h *TransportHTTP) adminGetModels(ctx context.Context, conditions map[string]interface{},
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

	if err := h.doHTTPRequest(
		ctx, http.MethodPost, path, jsonStr, h.xPriv, true, &models,
	); err != nil {
		return err
	}

	return nil
}

func (h *TransportHTTP) adminCount(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata, path string) (int64, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldConditions: conditions,
		FieldMetadata:   processMetadata(metadata),
	})
	if err != nil {
		return 0, WrapError(err)
	}

	var count int64
	if err := h.doHTTPRequest(
		ctx, http.MethodPost, path, jsonStr, h.xPriv, true, &count,
	); err != nil {
		return 0, err
	}

	return count, nil
}

// AdminRecordTransaction will record a transaction as an admin
func (h *TransportHTTP) AdminRecordTransaction(ctx context.Context, hex string) (*models.Transaction, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldHex: hex,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var transaction models.Transaction
	if err := h.doHTTPRequest(
		ctx, http.MethodPost, "/admin/transactions/record", jsonStr, h.xPriv, h.signRequest, &transaction,
	); err != nil {
		return nil, err
	}

	return &transaction, nil
}

// AdminGetSharedConfig gets the shared config
func (h *TransportHTTP) AdminGetSharedConfig(ctx context.Context) (*models.SharedConfig, ResponseError) {
	var model *models.SharedConfig
	if err := h.doHTTPRequest(
		ctx, http.MethodGet, "/admin/shared-config", nil, h.xPriv, true, &model,
	); err != nil {
		return nil, err
	}

	return model, nil
}

// AdminGetContacts executes an HTTP POST request to search for contacts based on specified conditions, metadata, and query parameters.
func (h *TransportHTTP) AdminGetContacts(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata, queryParams *QueryParams) ([]*models.Contact, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldConditions:  conditions,
		FieldMetadata:    processMetadata(metadata),
		FieldQueryParams: queryParams,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var contacts []*models.Contact
	err = h.doHTTPRequest(ctx, http.MethodPost, "/admin/contact/search", jsonStr, h.adminXPriv, true, &contacts)
	return contacts, WrapError(err)
}

// AdminUpdateContact executes an HTTP PATCH request to update a specific contact's full name using their ID.
func (h *TransportHTTP) AdminUpdateContact(ctx context.Context, id, fullName string, metadata *models.Metadata) (*models.Contact, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		"fullName":    fullName,
		FieldMetadata: processMetadata(metadata),
	})
	if err != nil {
		return nil, WrapError(err)
	}
	var contact models.Contact
	err = h.doHTTPRequest(ctx, http.MethodPatch, fmt.Sprintf("/admin/contact/%s", id), jsonStr, h.adminXPriv, true, &contact)
	return &contact, WrapError(err)
}

// AdminDeleteContact executes an HTTP DELETE request to remove a contact using their ID.
func (h *TransportHTTP) AdminDeleteContact(ctx context.Context, id string) ResponseError {
	err := h.doHTTPRequest(ctx, http.MethodDelete, fmt.Sprintf("/admin/contact/%s", id), nil, h.adminXPriv, true, nil)
	return WrapError(err)
}

// AdminAcceptContact executes an HTTP PATCH request to mark a contact as accepted using their ID.
func (h *TransportHTTP) AdminAcceptContact(ctx context.Context, id string) (*models.Contact, ResponseError) {
	var contact models.Contact
	err := h.doHTTPRequest(ctx, http.MethodPatch, fmt.Sprintf("/admin/contact/accepted/%s", id), nil, h.adminXPriv, true, &contact)
	return &contact, WrapError(err)
}

// AdminRejectContact executes an HTTP PATCH request to mark a contact as rejected using their ID.
func (h *TransportHTTP) AdminRejectContact(ctx context.Context, id string) (*models.Contact, ResponseError) {
	var contact models.Contact
	err := h.doHTTPRequest(ctx, http.MethodPatch, fmt.Sprintf("/admin/contact/rejected/%s", id), nil, h.adminXPriv, true, &contact)
	return &contact, WrapError(err)
}
