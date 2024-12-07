package stats

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
)

const (
	route = "/v1/admin/stats"
	api   = "Admin Stats API"
)

type API struct {
	url        *url.URL
	httpClient *resty.Client
}

func (a *API) Stats(ctx context.Context) (*models.AdminStats, error) {
	var result models.AdminStats
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		Get(a.url.String())
	if err != nil {
		return nil, fmt.Errorf("HTTP request failure: %w", err)
	}

	return &result, nil
}

func NewAPI(url *url.URL, httpClient *resty.Client) *API {
	return &API{url: url.JoinPath(route), httpClient: httpClient}
}

func HTTPErrorFormatter(action string, err error) *errutil.HTTPErrorFormatter {
	return &errutil.HTTPErrorFormatter{
		Action: action,
		API:    api,
		Err:    err,
	}
}
