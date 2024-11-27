package invitations

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/go-resty/resty/v2"
)

const route = "api/v1/invitations"

type API struct {
	url        *url.URL
	httpClient *resty.Client
}

func (a *API) AcceptInvitation(ctx context.Context, paymail string) error {
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		Post(a.url.JoinPath(paymail, "contacts").String())
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func (a *API) RejectInvitation(ctx context.Context, paymail string) error {
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		Delete(a.url.JoinPath(paymail).String())
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
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
		API:    "User Invitations API",
		Err:    err,
	}
}
