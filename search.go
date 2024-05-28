package walletclient

import (
	"context"
	"encoding/json"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/libsv/go-bk/bip32"
)

// SearchRequester is a function that sends a request to the server and returns the response.
type SearchRequester func(ctx context.Context, method string, path string, rawJSON []byte, xPriv *bip32.ExtendedKey, sign bool, responseJSON interface{}) ResponseError

// Search prepares and sends a search request to the server.
func Search[TFilter any, TItem any](
	ctx context.Context,
	method string,
	path string,
	xPriv *bip32.ExtendedKey,
	f TFilter,
	metadata map[string]interface{},
	queryParams *filter.QueryParams,
	requester SearchRequester,
) ([]*TItem, ResponseError) {
	jsonStr, err := json.Marshal(filter.SearchModel[TFilter]{
		ConditionsModel: filter.ConditionsModel[TFilter]{
			Conditions: f,
			Metadata:   toMapPtr(metadata),
		},
		QueryParams: queryParams,
	})
	if err != nil {
		return nil, WrapError(err)
	}
	var items []*TItem
	if err := requester(ctx, method, path, jsonStr, xPriv, true, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// Count prepares and sends a count request to the server.
func Count[TFilter any](
	ctx context.Context,
	method string,
	path string,
	xPriv *bip32.ExtendedKey,
	f TFilter,
	metadata map[string]interface{},
	requester SearchRequester,
) (int64, ResponseError) {
	jsonStr, err := json.Marshal(filter.ConditionsModel[TFilter]{
		Conditions: f,
		Metadata:   toMapPtr(metadata),
	})
	if err != nil {
		return 0, WrapError(err)
	}
	var count int64
	if err := requester(ctx, method, path, jsonStr, xPriv, true, &count); err != nil {
		return 0, err
	}

	return count, nil
}

func toMapPtr(m map[string]interface{}) *map[string]interface{} {
	if m == nil {
		return nil
	}
	return &m
}
