package merkleroots

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
)

const route = "api/v1/merkleroots"

type API struct {
	addr       string
	httpClient *resty.Client
}

func (a *API) MerkleRoots(ctx context.Context, opts ...queries.MerkleRootsQueryOption) ([]*models.MerkleRoot, error) {
	var result models.ExclusiveStartKeyPage[[]*models.MerkleRoot]

	params := CreateQueryParams(opts...)
	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		SetQueryParams(params).
		Get(a.addr)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return result.Content, nil
}

func CreateQueryParams(opts ...queries.MerkleRootsQueryOption) map[string]string {
	var q queries.MerkleRootsQuery
	for _, o := range opts {
		o(&q)
	}

	params := make(map[string]string)
	if q.BatchSize > 0 {
		params["batchSize"] = strconv.FormatInt(int64(q.BatchSize), 10)
	}
	if strings.TrimSpace(q.LastEvaluatedKey) != "" {
		params["lastEvaluatedKey"] = q.LastEvaluatedKey
	}

	return params
}

func NewAPI(addr string, cli *resty.Client) *API {
	return &API{
		addr:       addr + "/" + route,
		httpClient: cli,
	}
}
