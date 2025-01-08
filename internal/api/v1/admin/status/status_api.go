package status

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
)

const (
	route = "v1/admin/status"
	api   = "Admin Status API"
)

type API struct {
	httpClient *resty.Client
	url        *url.URL
}

func (a *API) Status(ctx context.Context) (bool, error) {
	res, err := a.httpClient.
		R().
		SetContext(ctx).
		Get(a.url.String())
	if err != nil {
		if res.StatusCode() == http.StatusUnauthorized {
			return false, nil
		}
		return false, fmt.Errorf("HTTP response failure: %w", err)
	}

	return true, nil
}

func NewAPI(url *url.URL, httpClient *resty.Client) *API {
	return &API{url: url.JoinPath(route), httpClient: httpClient}
}
