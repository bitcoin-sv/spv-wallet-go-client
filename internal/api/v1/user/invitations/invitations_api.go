package invitations

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

const route = "api/v1/invitations"

type API struct {
	addr       string
	httpClient *resty.Client
}

func (a *API) AcceptInvitation(ctx context.Context, paymail string) error {
	URL := a.addr + "/" + paymail + "/contacts"
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		Post(URL)
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func (a *API) RejectInvitation(ctx context.Context, paymail string) error {
	URL := a.addr + "/" + paymail
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		Delete(URL)
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func NewAPI(addr string, httpClient *resty.Client) *API {
	return &API{
		addr:       addr + "/" + route,
		httpClient: httpClient,
	}
}
