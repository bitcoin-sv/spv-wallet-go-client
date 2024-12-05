package invitations

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/go-resty/resty/v2"
)

const (
	route = "api/v1/admin/invitations"
	api   = "Admin Invitations API"
)

type API struct {
	httpClient *resty.Client
	url        *url.URL
}

func (a *API) AcceptInvitation(ctx context.Context, ID string) error {
	URL := a.url.JoinPath(ID).String()
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		Post(URL)
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func (a *API) RejectInvitation(ctx context.Context, ID string) error {
	URL := a.url.JoinPath(ID).String()
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		Delete(URL)
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
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
