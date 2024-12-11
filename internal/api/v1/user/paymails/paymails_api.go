package paymails

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/go-resty/resty/v2"
)

const (
	route = "api/v1/paymails"
	api   = "User Paymails API"
)

type API struct {
	url        *url.URL
	httpClient *resty.Client
}

func (a *API) Paymails(ctx context.Context, opts ...queries.PaymailQueryOption[filter.PaymailFilter]) (*queries.PaymailAddressPage, error) {
	var query queries.PaymailQuery[filter.PaymailFilter]
	for _, o := range opts {
		o(&query)
	}

	queryBuilder := querybuilders.NewQueryBuilder(
		querybuilders.WithMetadataFilter(query.Metadata),
		querybuilders.WithPageFilter(query.PageFilter),
		querybuilders.WithFilterQueryBuilder(&PaymailFilterBuilder{
			PaymailFilter:      query.PaymailFilter,
			ModelFilterBuilder: querybuilders.ModelFilterBuilder{ModelFilter: query.PaymailFilter.ModelFilter},
		}),
	)
	params, err := queryBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build paymail address query params: %w", err)
	}

	var result queries.PaymailAddressPage
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
