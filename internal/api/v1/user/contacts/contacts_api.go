package contacts

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/queryparams"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
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

func (a *API) Contacts(ctx context.Context, opts ...queries.QueryOption[filter.ContactFilter]) (*queries.ContactsPage, error) {
	query := queries.NewQuery(opts...)
	parser, err := queryparams.NewQueryParser(query)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize query parser: %w", err)
	}

	params, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to build user contacts query params: %w", err)
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
		Put(a.url.JoinPath(cmd.ContactPaymail).String())
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
