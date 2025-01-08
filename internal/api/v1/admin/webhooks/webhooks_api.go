package webhooks

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/go-resty/resty/v2"
)

const (
	route = "/api/v1/admin/webhooks/subscriptions"
	api   = "Admin Webhooks API"
)

type API struct {
	httpClient *resty.Client
	url        *url.URL
}

func (a *API) SubscribeWebhook(ctx context.Context, cmd *commands.CreateWebhookSubscription) error {
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		SetBody(cmd).
		Post(a.url.String())
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func (a *API) UnsubscribeWebhook(ctx context.Context, cmd *commands.CancelWebhookSubscription) error {
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		SetBody(cmd).
		Delete(a.url.String())
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func NewAPI(url *url.URL, httpClient *resty.Client) *API {
	return &API{url: url.JoinPath(route), httpClient: httpClient}
}
