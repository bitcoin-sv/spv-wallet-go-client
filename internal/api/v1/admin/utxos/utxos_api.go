package utxos

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/queryparams"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/go-resty/resty/v2"
)

const (
	route = "/api/v1/admin/utxos"
	api   = "Admin UTXOs API"
)

type API struct {
	httpClient *resty.Client
	url        *url.URL
}

func (a *API) UTXOs(ctx context.Context, opts ...queries.QueryOption[filter.AdminUtxoFilter]) (*queries.UtxosPage, error) {
	query := queries.NewQuery(opts...)
	parser, err := queryparams.NewQueryParser(query)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize query parser: %w", err)
	}

	params, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to build utxos query params: %w", err)
	}

	var result queries.UtxosPage
	_, err = a.httpClient.
		R().
		SetContext(ctx).
		SetQueryParams(params.ParseToMap()).
		SetResult(&result).
		Get(a.url.String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func NewAPI(url *url.URL, httpClient *resty.Client) *API {
	return &API{url: url.JoinPath(route), httpClient: httpClient}
}

func HTTPErrorFormatter(action string, err error) *errutil.HTTPErrorFormatter {
	return &errutil.HTTPErrorFormatter{
		Action: action,
		API:    api,
		Err:    err,
	}
}
