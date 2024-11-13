package configs

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const route = "api/v1/configs"

type API struct {
	addr       string
	httpClient *resty.Client
}

func (a *API) SharedConfig(ctx context.Context) (*response.SharedConfig, error) {
	var result response.SharedConfig
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		Get(a.addr + "/shared")
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
