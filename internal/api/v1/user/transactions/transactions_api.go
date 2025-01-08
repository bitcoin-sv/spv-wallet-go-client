package transactions

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
	route = "api/v1/transactions"
	api   = "User Transactions API"
)

type TransactionSigner interface {
	TransactionSignedHex(dt *response.DraftTransaction) (string, error)
}

type API struct {
	url               *url.URL
	httpClient        *resty.Client
	transactionSigner TransactionSigner
}

func (a *API) FinalizeTransaction(draft *response.DraftTransaction) (string, error) {
	hex, err := a.transactionSigner.TransactionSignedHex(draft)
	if err != nil {
		return "", fmt.Errorf("failed to finalize transaction: %w", err)
	}

	return hex, nil
}

func (a *API) DraftToRecipients(ctx context.Context, r *commands.SendToRecipients) (*response.DraftTransaction, error) {
	outputs := make([]*response.TransactionOutput, 0)

	for _, recipient := range r.Recipients {
		outputs = append(outputs, &response.TransactionOutput{
			To:       recipient.To,
			Satoshis: recipient.Satoshis,
			OpReturn: recipient.OpReturn,
		})
	}

	return a.DraftTransaction(ctx, &commands.DraftTransaction{
		Config: response.TransactionConfig{
			Outputs: outputs,
		},
		Metadata: r.Metadata,
	})
}

func (a *API) SendToRecipients(ctx context.Context, r *commands.SendToRecipients) (*response.Transaction, error) {
	draft, err := a.DraftToRecipients(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("failed to send draft to recipients: %w", err)
	}

	var hex string
	if hex, err = a.FinalizeTransaction(draft); err != nil {
		return nil, fmt.Errorf("failed to finalize transaction: %w", err)
	}

	return a.RecordTransaction(ctx, &commands.RecordTransaction{
		Metadata:    r.Metadata,
		Hex:         hex,
		ReferenceID: draft.ID,
	})
}

func (a *API) DraftTransaction(ctx context.Context, r *commands.DraftTransaction) (*response.DraftTransaction, error) {
	var result response.DraftTransaction

	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(r).
		Post(a.url.JoinPath("drafts").String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) RecordTransaction(ctx context.Context, r *commands.RecordTransaction) (*response.Transaction, error) {
	var result response.Transaction

	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(r).
		Post(a.url.String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) UpdateTransactionMetadata(ctx context.Context, r *commands.UpdateTransactionMetadata) (*response.Transaction, error) {
	var result response.Transaction

	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(r).
		Patch(a.url.JoinPath(r.ID).String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) Transaction(ctx context.Context, ID string) (*response.Transaction, error) {
	var result response.Transaction

	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Get(a.url.JoinPath(ID).String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) Transactions(ctx context.Context, opts ...queries.QueryOption[filter.TransactionFilter]) (*queries.TransactionPage, error) {
	query := queries.NewQuery(opts...)
	parser, err := queryparams.NewQueryParser(query)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize query parser: %w", err)
	}

	params, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to create transactions query params: %w", err)
	}

	var result response.PageModel[response.Transaction]
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

func NewAPIWithXPriv(URL *url.URL, httpClient *resty.Client, xPriv string) (*API, error) {
	transactionSigner, err := NewXPrivTransactionSigner(xPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactionSigner: %w", err)
	}

	return &API{
			url:               URL.JoinPath(route),
			httpClient:        httpClient,
			transactionSigner: transactionSigner},
		nil
}

func NewAPI(URL *url.URL, httpClient *resty.Client) (*API, error) {
	return &API{
		url:               URL.JoinPath(route),
		httpClient:        httpClient,
		transactionSigner: &noopTransactionSigner{},
	}, nil
}
