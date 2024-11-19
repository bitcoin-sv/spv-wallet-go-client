package merkleroots

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/go-resty/resty/v2"
)

const route = "api/v1/merkleroots"

type API struct {
	url        *url.URL
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
		Get(a.url.String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func NewAPI(url *url.URL, httpClient *resty.Client) *API {
	return &API{
		url:        url.JoinPath(route),
		httpClient: httpClient,
	}
}
