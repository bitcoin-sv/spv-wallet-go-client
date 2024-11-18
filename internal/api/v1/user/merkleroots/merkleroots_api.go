package merkleroots

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/go-resty/resty/v2"
)

const route = "api/v1/merkleroots"

type API struct {
	addr       string
	httpClient *resty.Client
}

func (a *API) MerkleRoots(ctx context.Context, merkleRootOpts ...queries.MerkleRootsQueryOption) (*queries.MerkleRootPage, error) {
	var query queries.MerkleRootsQuery
	for _, o := range merkleRootOpts {
		o(&query)
	}

	queryBuilder := querybuilders.NewQueryBuilder(querybuilders.WithFilterQueryBuilder(&merkleRootsFilterQueryBuilder{query: query}))
	params, err := queryBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build merkle roots query params: %w", err)
	}

	var result queries.MerkleRootPage
	_, err = a.httpClient.R().
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
