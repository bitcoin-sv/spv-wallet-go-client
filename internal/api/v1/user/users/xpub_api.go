package users

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const route = "api/v1/users/current"

type XPubAPI struct {
	addr       string
	httpClient *resty.Client
}

func (x *XPubAPI) XPub(ctx context.Context) (*response.Xpub, error) {
	var result response.Xpub
	_, err := x.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		Get(x.addr)
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
		Patch(x.addr)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func NewXPubAPI(addr string, httpClient *resty.Client) *XPubAPI {
	return &XPubAPI{
		addr:       addr + "/" + route,
		httpClient: httpClient,
	}
}
