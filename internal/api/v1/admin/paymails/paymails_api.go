package paymails

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const (
	route = "api/v1/admin/paymails"
	api   = "Admin User Paymails API"
)

type API struct {
	httpClient *resty.Client
	url        *url.URL
}

func (a *API) DeletePaymail(ctx context.Context, address string) error {
	type body struct {
		Address string `json:"address"`
	}

	_, err := a.httpClient.
		R().
		SetContext(ctx).
		SetBody(body{Address: address}).
		Delete(a.url.JoinPath(address).String())
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func (a *API) CreatePaymail(ctx context.Context, cmd *commands.CreatePaymail) (*response.PaymailAddress, error) {
	var result response.PaymailAddress
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(cmd).
		Post(a.url.String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) Paymail(ctx context.Context, ID string) (*response.PaymailAddress, error) {
	var result response.PaymailAddress
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		Get(a.url.JoinPath(ID).String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) Paymails(ctx context.Context, opts ...queries.QueryOption[filter.AdminPaymailFilter]) (*queries.PaymailsPage, error) {
	query := queries.NewQuery(opts...)
	queryBuilder := querybuilders.NewQueryBuilder(
		querybuilders.WithMetadataFilter(query.Metadata),
		querybuilders.WithPageFilter(query.PageFilter),
		querybuilders.WithFilterQueryBuilder(&adminPaymailFilterBuilder{paymailFilter: query.Filter}),
	)
	params, err := queryBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build paymail address query params: %w", err)
	}

	var result queries.PaymailsPage
	_, err = a.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		SetQueryParams(params.ParseToMap()).
		Get(a.url.String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
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
