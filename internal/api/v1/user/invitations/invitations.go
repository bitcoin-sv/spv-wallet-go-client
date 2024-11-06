package invitations

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

const route = "api/v1/invitations"

type API struct {
	addr string
	cli  *resty.Client
}

func (a *API) AcceptInvitation(ctx context.Context, paymail string) error {
	URL := a.addr + "/" + paymail + "/contacts"
	_, err := a.cli.
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
	_, err := a.cli.
		R().
		SetContext(ctx).
		Delete(URL)
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func NewAPI(addr string, cli *resty.Client) *API {
	return &API{
		addr: addr + "/" + route,
		cli:  cli,
	}
}
