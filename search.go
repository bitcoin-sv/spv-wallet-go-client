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
func Search[TFilter any, TResp any](
	ctx context.Context,
	method string,
	path string,
	xPriv *bip32.ExtendedKey,
	f *TFilter,
	metadata map[string]any,
	queryParams *filter.QueryParams,
	requester SearchRequester,
) (TResp, ResponseError) {
	jsonStr, err := json.Marshal(filter.SearchModel[TFilter]{
		ConditionsModel: filter.ConditionsModel[TFilter]{
			Conditions: f,
			Metadata:   toMapPtr(metadata),
		},
		QueryParams: queryParams,
	})
	var resp TResp // before initialization, this var is empty slice or nil so it can be returned in case of error
	if err != nil {
		return resp, WrapError(err)
	}

	if err := requester(ctx, method, path, jsonStr, xPriv, true, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

// Count prepares and sends a count request to the server.
func Count[TFilter any](
	ctx context.Context,
	method string,
	path string,
	xPriv *bip32.ExtendedKey,
	f *TFilter,
	metadata map[string]any,
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

// Optional returns a pointer to provided value, it's necessary to define "optional" fields in filters
func Optional[T any](val T) *T {
	return &val
}

func toMapPtr(m map[string]any) *map[string]any {
	if m == nil {
		return nil
	}
	return &m
}
