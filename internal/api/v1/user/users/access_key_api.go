package users

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

type AccessKeyAPI struct {
	addr       string
	httpClient *resty.Client
}

func (a *AccessKeyAPI) GenerateAccessKey(ctx context.Context, cmd *commands.GenerateAccessKey) (*response.AccessKey, error) {
	var result response.AccessKey

	URL := a.addr + "/keys"
	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(cmd).
		Post(URL)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *AccessKeyAPI) AccessKey(ctx context.Context, ID string) (*response.AccessKey, error) {
	var result response.AccessKey

	URL := a.addr + "/keys/" + ID
	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Get(URL)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *AccessKeyAPI) AccessKeys(ctx context.Context, opts ...queries.AccessKeyQueryOption) (*queries.AccessKeyPage, error) {
	var query queries.AccessKeyQuery
	for _, o := range opts {
		o(&query)
	}

	queryBuilder := querybuilders.NewQueryBuilder(
		querybuilders.WithMetadataFilter(query.Metadata),
		querybuilders.WithPageFilter(query.PageFilter),
		querybuilders.WithFilterQueryBuilder(&accessKeyFilterQueryBuilder{
			accessKeyFilter:    query.AccessKeyFilter,
			modelFilterBuilder: querybuilders.ModelFilterBuilder{ModelFilter: query.AccessKeyFilter.ModelFilter},
		}),
	)
	params, err := queryBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build access keys query params: %w", err)
	}

	var result response.PageModel[response.AccessKey]
	URL := a.addr + "/keys"
	_, err = a.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		SetQueryParams(params.ParseToMap()).
		Get(URL)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *AccessKeyAPI) RevokeAccessKey(ctx context.Context, ID string) error {
	URL := a.addr + "/keys/" + ID
	_, err := a.httpClient.R().
		SetContext(ctx).
		Delete(URL)
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func NewAccessKeyAPI(addr string, httpClient *resty.Client) *AccessKeyAPI {
	return &AccessKeyAPI{
		addr:       addr + "/" + route,
		httpClient: httpClient,
	}
}
