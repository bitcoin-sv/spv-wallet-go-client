package contacts

import (
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type ContactFilterQueryBuilder struct {
	ModelFilterBuilder querybuilders.ModelFilterBuilder
	ContactFilter      filter.ContactFilter
}

func (c *ContactFilterQueryBuilder) Build() (url.Values, error) {
	modelFilterBuilder, err := c.ModelFilterBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build model filter query params: %w", err)
	}

	params := querybuilders.NewExtendedURLValues()
	if len(modelFilterBuilder) > 0 {
		params.Append(modelFilterBuilder)
	}

	params.AddPair("id", c.ContactFilter.ID)
	params.AddPair("fullName", c.ContactFilter.FullName)
	params.AddPair("paymail", c.ContactFilter.Paymail)
	params.AddPair("pubKey", c.ContactFilter.PubKey)
	params.AddPair("status", c.ContactFilter.Status)
	return params.Values, nil
}
