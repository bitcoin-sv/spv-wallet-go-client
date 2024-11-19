package transactions

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const route = "api/v1/transactions"

type API struct {
	url        *url.URL
	httpClient *resty.Client
}

func (a *API) DraftTransaction(ctx context.Context, r *commands.DraftTransaction) (*response.DraftTransaction, error) {
	var result response.DraftTransaction

	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(r).
		Post(a.url.JoinPath("drafts").String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) RecordTransaction(ctx context.Context, r *commands.RecordTransaction) (*response.Transaction, error) {
	var result response.Transaction

	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(r).
		Post(a.url.String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) UpdateTransactionMetadata(ctx context.Context, r *commands.UpdateTransactionMetadata) (*response.Transaction, error) {
	var result response.Transaction

	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(r).
		Patch(a.url.JoinPath(r.ID).String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
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

func (a *API) Transactions(ctx context.Context, transactionsOpts ...queries.TransactionsQueryOption) (*queries.TransactionPage, error) {
	var query queries.TransactionsQuery
	for _, o := range transactionsOpts {
		o(&query)
	}

	queryBuilder := querybuilders.NewQueryBuilder(querybuilders.WithMetadataFilter(query.Metadata),
		querybuilders.WithPageFilter(query.Page),
		querybuilders.WithFilterQueryBuilder(&transactionFilterBuilder{
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

func NewAPI(URL *url.URL, httpClient *resty.Client) *API {
	return &API{
		url:        URL.JoinPath(route),
		httpClient: httpClient,
	}
}
