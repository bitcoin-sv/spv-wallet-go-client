package contacts

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const route = "api/v1/contacts"

type API struct {
	addr       string
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
		querybuilders.WithFilterQueryBuilder(&contactFilterQueryBuilder{
			contactFilter: query.ContactFilter,
			modelFilterBuilder: querybuilders.ModelFilterBuilder{
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
		Get(a.addr)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) ContactWithPaymail(ctx context.Context, paymail string) (*response.Contact, error) {
	var result response.Contact

	URL := a.addr + "/" + paymail
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		Get(URL)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) UpsertContact(ctx context.Context, r commands.UpsertContact) (*response.Contact, error) {
	var result response.CreateContactResponse

	URL := a.addr + "/" + r.Paymail
	_, err := a.httpClient.
		R().
		SetBody(r).
		SetContext(ctx).
		SetResult(&result).
		Put(URL)
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

func (a *API) ConfirmContact(ctx context.Context, paymail string) error {
	URL := a.addr + "/" + paymail + "/confirmation"
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		Post(URL)
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func (a *API) UnconfirmContact(ctx context.Context, paymail string) error {
	URL := a.addr + "/" + paymail + "/confirmation"
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
