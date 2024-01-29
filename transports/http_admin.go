package transports

import (
	"context"
	"encoding/json"
	"net/http"

	buxmodels "github.com/BuxOrg/bux-models"
)

// NewXpub will register an xPub
func (h *TransportHTTP) NewXpub(ctx context.Context, rawXPub string, metadata *buxmodels.Metadata) ResponseError {
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

	var xPubData buxmodels.Xpub

	return h.doHTTPRequest(
		ctx, http.MethodPost, "/xpub", jsonStr, h.adminXPriv, true, &xPubData,
	)
}

// RegisterXpub alias for NewXpub
func (h *TransportHTTP) RegisterXpub(ctx context.Context, rawXPub string, metadata *buxmodels.Metadata) ResponseError {
	return h.NewXpub(ctx, rawXPub, metadata)
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
func (h *TransportHTTP) AdminGetStats(ctx context.Context) (*buxmodels.AdminStats, ResponseError) {
	var stats *buxmodels.AdminStats
	if err := h.doHTTPRequest(
		ctx, http.MethodGet, "/admin/stats", nil, h.xPriv, true, &stats,
	); err != nil {
		return nil, err
	}

	return stats, nil
}

// AdminGetAccessKeys get all access keys filtered by conditions
func (h *TransportHTTP) AdminGetAccessKeys(ctx context.Context, conditions map[string]interface{},
	metadata *buxmodels.Metadata, queryParams *QueryParams,
) ([]*buxmodels.AccessKey, ResponseError) {
	var models []*buxmodels.AccessKey
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/access-keys/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetAccessKeysCount get a count of all the access keys filtered by conditions
func (h *TransportHTTP) AdminGetAccessKeysCount(ctx context.Context, conditions map[string]interface{},
	metadata *buxmodels.Metadata,
) (int64, ResponseError) {
	return h.adminCount(ctx, conditions, metadata, "/admin/access-keys/count")
}

// AdminGetBlockHeaders get all block headers filtered by conditions
func (h *TransportHTTP) AdminGetBlockHeaders(ctx context.Context, conditions map[string]interface{},
	metadata *buxmodels.Metadata, queryParams *QueryParams,
) ([]*buxmodels.BlockHeader, ResponseError) {
	var models []*buxmodels.BlockHeader
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/block-headers/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetBlockHeadersCount get a count of all the block headers filtered by conditions
func (h *TransportHTTP) AdminGetBlockHeadersCount(ctx context.Context, conditions map[string]interface{},
	metadata *buxmodels.Metadata,
) (int64, ResponseError) {
	return h.adminCount(ctx, conditions, metadata, "/admin/block-headers/count")
}

// AdminGetDestinations get all block destinations filtered by conditions
func (h *TransportHTTP) AdminGetDestinations(ctx context.Context, conditions map[string]interface{},
	metadata *buxmodels.Metadata, queryParams *QueryParams,
) ([]*buxmodels.Destination, ResponseError) {
	var models []*buxmodels.Destination
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/destinations/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetDestinationsCount get a count of all the destinations filtered by conditions
func (h *TransportHTTP) AdminGetDestinationsCount(ctx context.Context, conditions map[string]interface{},
	metadata *buxmodels.Metadata,
) (int64, ResponseError) {
	return h.adminCount(ctx, conditions, metadata, "/admin/destinations/count")
}

// AdminGetPaymail get a paymail by address
func (h *TransportHTTP) AdminGetPaymail(ctx context.Context, address string) (*buxmodels.PaymailAddress, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldAddress: address,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var model *buxmodels.PaymailAddress
	if err := h.doHTTPRequest(
		ctx, http.MethodPost, "/admin/paymail/get", jsonStr, h.xPriv, true, &model,
	); err != nil {
		return nil, err
	}

	return model, nil
}

// AdminGetPaymails get all block paymails filtered by conditions
func (h *TransportHTTP) AdminGetPaymails(ctx context.Context, conditions map[string]interface{},
	metadata *buxmodels.Metadata, queryParams *QueryParams,
) ([]*buxmodels.PaymailAddress, ResponseError) {
	var models []*buxmodels.PaymailAddress
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/paymails/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetPaymailsCount get a count of all the paymails filtered by conditions
func (h *TransportHTTP) AdminGetPaymailsCount(ctx context.Context, conditions map[string]interface{},
	metadata *buxmodels.Metadata,
) (int64, ResponseError) {
	return h.adminCount(ctx, conditions, metadata, "/admin/paymails/count")
}

// AdminCreatePaymail create a new paymail for a xpub
func (h *TransportHTTP) AdminCreatePaymail(ctx context.Context, xPubID string, address string, publicName string, avatar string) (*buxmodels.PaymailAddress, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldXpubID:     xPubID,
		FieldAddress:    address,
		FieldPublicName: publicName,
		FieldAvatar:     avatar,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var model *buxmodels.PaymailAddress
	if err := h.doHTTPRequest(
		ctx, http.MethodPost, "/admin/paymail/create", jsonStr, h.xPriv, true, &model,
	); err != nil {
		return nil, err
	}

	return model, nil
}

// AdminDeletePaymail delete a paymail address from the database
func (h *TransportHTTP) AdminDeletePaymail(ctx context.Context, address string) (*buxmodels.PaymailAddress, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldAddress: address,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var model *buxmodels.PaymailAddress
	if err := h.doHTTPRequest(
		ctx, http.MethodDelete, "/admin/paymail/delete", jsonStr, h.xPriv, true, &model,
	); err != nil {
		return nil, err
	}

	return model, nil
}

// AdminGetTransactions get all block transactions filtered by conditions
func (h *TransportHTTP) AdminGetTransactions(ctx context.Context, conditions map[string]interface{},
	metadata *buxmodels.Metadata, queryParams *QueryParams,
) ([]*buxmodels.Transaction, ResponseError) {
	var models []*buxmodels.Transaction
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/transactions/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetTransactionsCount get a count of all the transactions filtered by conditions
func (h *TransportHTTP) AdminGetTransactionsCount(ctx context.Context, conditions map[string]interface{},
	metadata *buxmodels.Metadata,
) (int64, ResponseError) {
	return h.adminCount(ctx, conditions, metadata, "/admin/transactions/count")
}

// AdminGetUtxos get all block utxos filtered by conditions
func (h *TransportHTTP) AdminGetUtxos(ctx context.Context, conditions map[string]interface{},
	metadata *buxmodels.Metadata, queryParams *QueryParams,
) ([]*buxmodels.Utxo, ResponseError) {
	var models []*buxmodels.Utxo
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/utxos/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetUtxosCount get a count of all the utxos filtered by conditions
func (h *TransportHTTP) AdminGetUtxosCount(ctx context.Context, conditions map[string]interface{},
	metadata *buxmodels.Metadata,
) (int64, ResponseError) {
	return h.adminCount(ctx, conditions, metadata, "/admin/utxos/count")
}

// AdminGetXPubs get all block xpubs filtered by conditions
func (h *TransportHTTP) AdminGetXPubs(ctx context.Context, conditions map[string]interface{},
	metadata *buxmodels.Metadata, queryParams *QueryParams,
) ([]*buxmodels.Xpub, ResponseError) {
	var models []*buxmodels.Xpub
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/xpubs/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetXPubsCount get a count of all the xpubs filtered by conditions
func (h *TransportHTTP) AdminGetXPubsCount(ctx context.Context, conditions map[string]interface{},
	metadata *buxmodels.Metadata,
) (int64, ResponseError) {
	return h.adminCount(ctx, conditions, metadata, "/admin/xpubs/count")
}

func (h *TransportHTTP) adminGetModels(ctx context.Context, conditions map[string]interface{},
	metadata *buxmodels.Metadata, queryParams *QueryParams, path string, models interface{},
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
		ctx, http.MethodGet, path, jsonStr, h.xPriv, true, &models,
	); err != nil {
		return err
	}

	return nil
}

func (h *TransportHTTP) adminCount(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata, path string) (int64, ResponseError) {
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
func (h *TransportHTTP) AdminRecordTransaction(ctx context.Context, hex string) (*buxmodels.Transaction, ResponseError) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldHex: hex,
	})
	if err != nil {
		return nil, WrapError(err)
	}

	var transaction buxmodels.Transaction
	if err := h.doHTTPRequest(
		ctx, http.MethodPost, "/admin/transactions/record", jsonStr, h.xPriv, h.signRequest, &transaction,
	); err != nil {
		return nil, err
	}

	return &transaction, nil
}
