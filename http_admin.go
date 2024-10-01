package walletclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

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

// AdminGetWebhooks gets all webhooks
func (wc *WalletClient) AdminGetWebhooks(ctx context.Context) ([]*models.Webhook, error) {
	var webhooks []*models.Webhook
	err := wc.doHTTPRequest(ctx, http.MethodGet, "/admin/webhooks/subscriptions", nil, wc.adminXPriv, true, &webhooks)
	if err != nil {
		return nil, WrapError(err)
	}
	return webhooks, nil
}
