package transactions

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/queryparams"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const (
	route = "api/v1/admin/transactions"
	api   = "Admin Transactions API"
)

type API struct {
	httpClient *resty.Client
	url        *url.URL
}

func (a *API) Transaction(ctx context.Context, ID string) (*response.Transaction, error) {
	var result response.Transaction
	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Get(a.url.JoinPath(ID).String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) Transactions(ctx context.Context, opts ...queries.QueryOption[filter.AdminTransactionFilter]) (*queries.TransactionPage, error) {
	query := queries.NewQuery(opts...)
	parser, err := queryparams.NewQueryParser(query)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize query parser: %w", err)
	}

	params, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to create transactions query params: %w", err)
	}

	var result response.PageModel[response.Transaction]
	_, err = a.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		SetQueryParams(params.ParseToMap()).
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
