package queryparams

import (
	"fmt"
	"reflect"
	"strings"

	goclienterr "github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type QueryParser[F queries.QueryFilters] struct {
	query *queries.Query[F]
}

func (q *QueryParser[F]) contaisModelFilter(t reflect.Type) bool {
	// Check if the type directly matches ModelFilter
	if t == reflect.TypeOf(filter.ModelFilter{}) {
		return true
	}
	// If the input is a struct, check its fields for ModelFilter
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Type == reflect.TypeOf(filter.ModelFilter{}) {
				return true
			}
		}
	}
	return false
}

func (q *QueryParser[F]) parse(val any, totalParams *URLValues) {
	t := reflect.TypeOf(val)
	v := reflect.ValueOf(val)

	for i := 0; i < reflect.TypeOf(val).NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if q.contaisModelFilter(field.Type) {
			q.parse(value.Interface(), totalParams)
			continue
		}

		if value.Kind() == reflect.Ptr && value.IsNil() {
			continue
		}

		tags := strings.Split(field.Tag.Get("json"), ",")
		if len(tags) == 0 {
			continue
		}
		tag := tags[0]

		switch field.Type {
		case reflect.PointerTo(reflect.TypeOf(filter.TimeRange{})):
			totalParams.AddPair(tag, value.Interface().(*filter.TimeRange))

		case reflect.PointerTo(reflect.TypeOf(false)):
			totalParams.AddPair(tag, value.Elem().Bool())

		case reflect.PointerTo(reflect.TypeOf("")):
			totalParams.AddPair(tag, value.Elem().String())

		case reflect.PointerTo(reflect.TypeOf(uint64(0))), reflect.PointerTo(reflect.TypeOf(uint32(0))):
			totalParams.AddPair(tag, value.Elem().Uint())

		default:
			totalParams.AddPair(tag, fmt.Sprintf("%v", value.Interface()))
		}
	}
}

func (q *QueryParser[F]) Parse() (*URLValues, error) {
	totalParams := NewURLValues()
	// Parse page filter if present
	if q.query.PageFilter != (filter.Page{}) {
		q.parse(q.query.PageFilter, totalParams)
	}

	// Parse metadata
	metadata := &MetadataParser{Metadata: q.query.Metadata, MaxDepth: DefaultMaxDepth}
	params, err := metadata.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}
	totalParams.Append(params)

	// Parse main query filter
	q.parse(q.query.Filter, totalParams)
	return totalParams, nil
}

func NewQueryParser[F queries.QueryFilters](query *queries.Query[F]) (*QueryParser[F], error) {
	if query == nil {
		return nil, goclienterr.ErrQueryParserFailed
	}

	return &QueryParser[F]{query: query}, nil
}
