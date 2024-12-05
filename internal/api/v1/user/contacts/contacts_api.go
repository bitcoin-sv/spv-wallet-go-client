package contacts

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const (
	route = "api/v1/contacts"
	api   = "User Contacts API"
)

type API struct {
	url        *url.URL
	httpClient *resty.Client
}

func (a *API) Contacts(ctx context.Context, opts ...queries.ContactQueryOption) (*queries.UserContactsPage, error) {
	var query queries.ContactQuery
	for _, o := range opts {
		o(&query)
	}

	queryBuilder := querybuilders.NewQueryBuilder(
		querybuilders.WithMetadataFilter(query.Metadata),
		querybuilders.WithPageFilter(query.PageFilter),
		querybuilders.WithFilterQueryBuilder(&ContactFilterQueryBuilder{
			ContactFilter: query.ContactFilter,
			ModelFilterBuilder: querybuilders.ModelFilterBuilder{
				ModelFilter: query.ContactFilter.ModelFilter,
			},
		}),
	)
	params, err := queryBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build user contacts query params: %w", err)
	}

	var result queries.UserContactsPage
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

func (a *API) ContactWithPaymail(ctx context.Context, paymail string) (*response.Contact, error) {
	var result response.Contact
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		Get(a.url.JoinPath(paymail).String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) UpsertContact(ctx context.Context, cmd commands.UpsertContact) (*response.Contact, error) {
	var result response.CreateContactResponse
	_, err := a.httpClient.
		R().
		SetBody(cmd).
		SetContext(ctx).
		SetResult(&result).
		Put(a.url.JoinPath(cmd.Paymail).String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &response.Contact{
		Model:    result.Contact.Model,
		ID:       result.Contact.ID,
		FullName: result.Contact.FullName,
		Paymail:  result.Contact.Paymail,
		PubKey:   result.Contact.PubKey,
		Status:   result.Contact.Status,
	}, nil
}

func (a *API) RemoveContact(ctx context.Context, paymail string) error {
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		Delete(a.url.JoinPath(paymail).String())
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func (a *API) ConfirmContact(ctx context.Context, paymail string) error {
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		Post(a.url.JoinPath(paymail, "confirmation").String())
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func (a *API) UnconfirmContact(ctx context.Context, paymail string) error {
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		Delete(a.url.JoinPath(paymail, "confirmation").String())
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
		API:    api,
		Err:    err,
	}
}
