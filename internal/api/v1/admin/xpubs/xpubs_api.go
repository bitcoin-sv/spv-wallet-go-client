package xpubs

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/queryparams"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const (
	route = "api/v1/admin/users"
	api   = "Admin XPub API"
)

type API struct {
	url        *url.URL
	httpClient *resty.Client
}

func (a *API) CreateXPub(ctx context.Context, cmd *commands.CreateUserXpub) (*response.Xpub, error) {
	var result response.Xpub
	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(cmd).
		Post(a.url.String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) XPubs(ctx context.Context, opts ...queries.QueryOption[filter.XpubFilter]) (*queries.XPubPage, error) {
	query := queries.NewQuery(opts...)
	parser, err := queryparams.NewQueryParser(query)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize query parser: %w", err)
	}

	params, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to build user xpubs query params: %w", err)
	}

	var result queries.XPubPage
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
