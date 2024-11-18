package utxos

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/go-resty/resty/v2"
)

const route = "api/v1/utxos"

type API struct {
	addr       string
	httpClient *resty.Client
}

func (a *API) UTXOs(ctx context.Context, opts ...queries.UtxoQueryOption) (*queries.UtxosPage, error) {
	var query queries.UtxoQuery
	for _, o := range opts {
		o(&query)
	}

	queryBuilder := querybuilders.NewQueryBuilder(
		querybuilders.WithMetadataFilter(query.Metadata),
		querybuilders.WithPageFilter(query.PageFilter),
		querybuilders.WithFilterQueryBuilder(&utxoFilterQueryBuilder{
			utxoFilter:         query.UtxoFilter,
			modelFilterBuilder: querybuilders.ModelFilterBuilder{ModelFilter: query.UtxoFilter.ModelFilter},
		}),
	)
	params, err := queryBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build utxo query params: %w", err)
	}

	var result queries.UtxosPage
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

func NewAPI(addr string, httpClient *resty.Client) *API {
	return &API{
		addr:       addr + "/" + route,
		httpClient: httpClient,
	}
}
