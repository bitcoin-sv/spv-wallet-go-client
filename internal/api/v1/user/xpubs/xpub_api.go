package xpubs

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const (
	route = "api/v1/users/current"
	api   = "User XPubs API"
)

type API struct {
	url        *url.URL
	httpClient *resty.Client
}

func (a *API) XPub(ctx context.Context) (*response.Xpub, error) {
	var result response.Xpub
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		Get(a.url.String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) UpdateXPubMetadata(ctx context.Context, cmd *commands.UpdateXPubMetadata) (*response.Xpub, error) {
	var result response.Xpub
	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(cmd).
		Patch(a.url.String())
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
func HTTPErrorFormatter(action string, err error) *errutil.HTTPErrorFormatter {
	return &errutil.HTTPErrorFormatter{
		Action: action,
		API:    api,
		Err:    err,
	}
}
