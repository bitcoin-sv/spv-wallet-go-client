package contacts

import (
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type contactFilterQueryBuilder struct {
	modelFilterBuilder querybuilders.ModelFilterBuilder
	contactFilter      filter.ContactFilter
}

func (c *contactFilterQueryBuilder) Build() (url.Values, error) {
	modelFilterBuilder, err := c.modelFilterBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build model filter query params: %w", err)
	}

	params := querybuilders.NewExtendedURLValues()
	if len(modelFilterBuilder) > 0 {
		params.Append(modelFilterBuilder)
	}

	params.AddPair("id", c.contactFilter.ID)
	params.AddPair("fullName", c.contactFilter.FullName)
	params.AddPair("paymail", c.contactFilter.Paymail)
	params.AddPair("pubKey", c.contactFilter.PubKey)
	params.AddPair("status", c.contactFilter.Status)
	return params.Values, nil
}
