package configs

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const (
	route = "api/v1/configs"
	api   = "Shared Config API"
)

type API struct {
	url        *url.URL
	httpClient *resty.Client
}

func (a *API) SharedConfig(ctx context.Context) (*response.SharedConfig, error) {
	var result response.SharedConfig
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		Get(a.url.JoinPath("shared").String())
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
