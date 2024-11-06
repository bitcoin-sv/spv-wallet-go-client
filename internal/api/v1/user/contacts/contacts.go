package contacts

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

// TODO: 1. Contacts should accept the optional query parameters similar to the transactions.
// TODO: 2. Remove all helpers from the test file, all necessary funcs should be included in the testfixtures pkg.

const route = "api/v1/contacts"

type API struct {
	addr string
	cli  *resty.Client
}

type UpsertContactRequest struct {
	FullName         string         `json:"fullName"`
	Metadata         map[string]any `json:"metadata"`
	RequesterPaymail string         `json:"requesterPaymail"`
}

func (a *API) Contacts(ctx context.Context) ([]*response.Contact, error) {
	var result response.PageModel[response.Contact]
	_, err := a.cli.
		R().
		SetContext(ctx).
		SetResult(&result).
		Get(a.addr)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return result.Content, nil
}

func (a *API) ContactWithPaymail(ctx context.Context, paymail string) (*response.Contact, error) {
	var result response.Contact

	URL := a.addr + "/" + paymail
	_, err := a.cli.
		R().
		SetContext(ctx).
		SetResult(&result).
		Get(URL)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) UpsertContact(ctx context.Context, r UpsertContactRequest) (*response.Contact, error) {
	var result response.CreateContactResponse

	URL := a.addr + "/" + r.RequesterPaymail
	_, err := a.cli.
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
	_, err := a.cli.
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
	_, err := a.cli.
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
