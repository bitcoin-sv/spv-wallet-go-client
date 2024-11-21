package querybuilders

import (
	"errors"
	"net/url"

	goclienterr "github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type QueryBuilderOption func(*QueryBuilder)

func WithMetadataFilter(m Metadata) QueryBuilderOption {
	return func(qb *QueryBuilder) {
		if m != nil {
			qb.builders = append(qb.builders, &MetadataFilterBuilder{MaxDepth: DefaultMaxDepth, Metadata: m})
		}
	}
}

func WithModelFilter(m filter.ModelFilter) QueryBuilderOption {
	var zero filter.ModelFilter
	return func(qb *QueryBuilder) {
		if m != zero {
			qb.builders = append(qb.builders, &ModelFilterBuilder{ModelFilter: m})
		}
	}
}

func WithPageFilter(p filter.Page) QueryBuilderOption {
	var zero filter.Page
	return func(qb *QueryBuilder) {
		if p != zero {
			qb.builders = append(qb.builders, &PageFilterBuilder{Page: p})
		}
	}
}

func WithFilterQueryBuilder(b FilterQueryBuilder) QueryBuilderOption {
	return func(qb *QueryBuilder) {
		if b != nil {
			qb.builders = append(qb.builders, b)
		}
	}
}

type FilterQueryBuilder interface {
	Build() (url.Values, error)
}

type QueryBuilder struct {
	builders []FilterQueryBuilder
}

func (q *QueryBuilder) Build() (*ExtendedURLValues, error) {
	params := NewExtendedURLValues()
	for _, builder := range q.builders {
		values, err := builder.Build()
		if err != nil {
			return nil, errors.Join(err, goclienterr.ErrFilterQueryBuilder)
		}

		if len(values) > 0 {
			params.Append(values)
		}
	}

	return params, nil
}

func NewQueryBuilder(opts ...QueryBuilderOption) *QueryBuilder {
	var qb QueryBuilder
	for _, o := range opts {
		o(&qb)
	}

	return &qb
}
