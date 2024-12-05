package transactions

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/transactions"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const route = "api/v1/admin/transactions"

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

func (a *API) Transactions(ctx context.Context, opts ...queries.TransactionsQueryOption) (*queries.TransactionPage, error) {
	var query queries.TransactionsQuery
	for _, o := range opts {
		o(&query)
	}

	queryBuilder := querybuilders.NewQueryBuilder(querybuilders.WithMetadataFilter(query.Metadata),
		querybuilders.WithPageFilter(query.Page),
		querybuilders.WithFilterQueryBuilder(&transactions.TransactionFilterBuilder{
			TransactionFilter:  query.Filter,
			ModelFilterBuilder: querybuilders.ModelFilterBuilder{ModelFilter: query.Filter.ModelFilter},
		}),
	)
	params, err := queryBuilder.Build()
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
		API:    "Admin Transactions API",
		Err:    err,
	}
}
