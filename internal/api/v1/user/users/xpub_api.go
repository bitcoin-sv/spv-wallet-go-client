package users

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const route = "api/v1/users/current"

type XPubAPI struct {
	url        *url.URL
	httpClient *resty.Client
}

func (x *XPubAPI) XPub(ctx context.Context) (*response.Xpub, error) {
	var result response.Xpub
	_, err := x.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		Get(x.url.String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (x *XPubAPI) UpdateXPubMetadata(ctx context.Context, cmd *commands.UpdateXPubMetadata) (*response.Xpub, error) {
	var result response.Xpub
	_, err := x.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(cmd).
		Patch(x.url.String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func NewXPubAPI(url *url.URL, httpClient *resty.Client) *XPubAPI {
	return &XPubAPI{
		url:        url.JoinPath(route),
		httpClient: httpClient,
	}

}
func XPubsHTTPErrorFormatter(action string, err error) *errutil.HTTPErrorFormatter {
	return &errutil.HTTPErrorFormatter{
		Action: action,
		API:    "User XPubs API",
		Err:    err,
	}
}
