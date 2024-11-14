package transactions

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const route = "api/v1/transactions"

type API struct {
	addr       string
	httpClient *resty.Client
}

func (a *API) DraftTransaction(ctx context.Context, r *commands.DraftTransaction) (*response.DraftTransaction, error) {
	var result response.DraftTransaction

	URL := a.addr + "/drafts"
	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(r).
		Post(URL)
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
		Post(a.addr)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) UpdateTransactionMetadata(ctx context.Context, r *commands.UpdateTransactionMetadata) (*response.Transaction, error) {
	var result response.Transaction

	URL := a.addr + "/" + r.ID
	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(r).
		Patch(URL)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) Transaction(ctx context.Context, ID string) (*response.Transaction, error) {
	var result response.Transaction

	URL := a.addr + "/" + ID
	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Get(URL)
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

	builderOpts := []querybuilders.QueryBuilderOption{
		querybuilders.WithMetadataFilter(query.Metadata),
		querybuilders.WithPageFilterQueryBuilder(query.Page),
		querybuilders.WithFilterQueryBuilder(&transactionFilterBuilder{
			TransactionFilter:  query.Filter,
			ModelFilterBuilder: querybuilders.ModelFilterBuilder{ModelFilter: query.Filter.ModelFilter},
		}),
	}
	builder := querybuilders.NewQueryBuilder(builderOpts...)
	params, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create transactions query params: %w", err)
	}

	var result response.PageModel[response.Transaction]
	_, err = a.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		SetQueryParams(params.ParseToMap()).
		Get(a.addr)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func NewAPI(addr string, cli *resty.Client) *API {
	return &API{
		addr:       addr + "/" + route,
		httpClient: cli,
	}
}
