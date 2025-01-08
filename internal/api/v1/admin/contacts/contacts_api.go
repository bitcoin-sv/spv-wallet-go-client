package contacts

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/queryparams"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const (
	route = "api/v1/admin/contacts"
	api   = "Admin Contacts API"
)

type API struct {
	httpClient *resty.Client
	url        *url.URL
}

func (a *API) CreateContact(ctx context.Context, cmd *commands.CreateContact) (*response.Contact, error) {
	var result response.Contact
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		SetBody(cmd).
		SetResult(&result).
		Post(a.url.JoinPath(cmd.Paymail).String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) Contacts(ctx context.Context, opts ...queries.QueryOption[filter.AdminContactFilter]) (*queries.ContactsPage, error) {
	query := queries.NewQuery(opts...)
	parser, err := queryparams.NewQueryParser(query)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize query parser: %w", err)
	}

	params, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to build admin contacts query params: %w", err)
	}

	var result queries.ContactsPage
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

func (a *API) UpdateContact(ctx context.Context, cmd *commands.UpdateContact) (*response.Contact, error) {
	var result response.Contact
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(cmd).
		Put(a.url.JoinPath(cmd.ID).String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) ConfirmContacts(ctx context.Context, cmd *commands.ConfirmContacts) error {
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		SetBody(cmd).
		Post(a.url.JoinPath("confirmations").String())
	if err != nil {
		return fmt.Errorf("HTTP response failure :%w", err)
	}
	return nil
}

func (a *API) DeleteContact(ctx context.Context, ID string) error {
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		Delete(a.url.JoinPath(ID).String())
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func NewAPI(url *url.URL, httpClient *resty.Client) *API {
	return &API{url: url.JoinPath(route), httpClient: httpClient}
}
