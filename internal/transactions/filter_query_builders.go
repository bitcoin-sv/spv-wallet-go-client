package transactions

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type Metadata map[string]any

const DefaultMaxDepth = 100

type MetadataFilterQueryBuilder struct {
	MaxDepth int
	Metadata Metadata
}

func (m *MetadataFilterQueryBuilder) Build() (url.Values, error) {
	var sb strings.Builder
	for k, v := range m.Metadata {
		var path strings.Builder
		pref := fmt.Sprintf("metadata[%s]", k)
		path.WriteString(pref)
		if err := m.dfs(0, &path, v, &sb, pref); err != nil {
			return nil, err
		}
		path.Reset()
	}

	params := make(url.Values)
	ss := strings.Split(TrimLastAmpersand(sb.String()), "&")
	for _, s := range ss {
		if len(s) == 0 {
			continue
		}
		p := strings.Split(s, "=")
		params.Add(p[0], p[1])
	}
	return params, nil
}

func (m *MetadataFilterQueryBuilder) dfs(depth int, path *strings.Builder, val any, ans *strings.Builder, pref string) error {
	if depth > m.MaxDepth {
		return fmt.Errorf("%w - max depth: %d", ErrMetadataFilterMaxDepthExceeded, m.MaxDepth)
	}

	t := reflect.TypeOf(val)
	switch t.Kind() {
	case reflect.Map:
		if err := m.mapDfs(depth+1, val, path, ans, pref); err != nil {
			return err
		}
	case reflect.Slice:
		if err := m.sliceDfs(depth+1, val, path, ans, pref); err != nil {
			return err
		}
	default:
		path.WriteString(fmt.Sprintf("=%v&", val))
		ans.WriteString(path.String())
	}
	return nil
}

func (m *MetadataFilterQueryBuilder) mapDfs(depth int, val any, path *strings.Builder, ans *strings.Builder, pref string) error {
	rval := reflect.ValueOf(val)
	for _, k := range rval.MapKeys() {
		mpv := rval.MapIndex(k)
		path.WriteString(fmt.Sprintf("[%v]", k.Interface()))
		if err := m.dfs(depth+1, path, mpv.Interface(), ans, pref); err != nil {
			return err
		}
		// Reset path after processing each map entry (backtracking).
		str := path.String()
		trim := str[:strings.LastIndex(str, "[")]
		path.Reset()
		path.WriteString(trim)
	}
	return nil
}

func (m *MetadataFilterQueryBuilder) sliceDfs(depth int, val any, path *strings.Builder, ans *strings.Builder, pref string) error {
	slice := reflect.ValueOf(val)
	for i := 0; i < slice.Len(); i++ {
		path.WriteString("[]")
		slv := slice.Index(i)
		if err := m.dfs(depth+1, path, slv.Interface(), ans, pref); err != nil {
			return err
		}
		// Reset path after processing each slice element (backtracking).
		str := path.String()
		trim := str[:strings.LastIndex(str, "[]")]
		path.Reset()
		path.WriteString(trim)
	}
	return nil
}

var ErrMetadataFilterMaxDepthExceeded = errors.New("maximum depth of nesting in metadata map exceeded")

func TrimLastAmpersand(s string) string {
	if len(s) > 0 && s[len(s)-1] == '&' {
		return s[:len(s)-1]
	}
	return s
}

type ModelFilterQueryBuilder struct {
	ModelFilter filter.ModelFilter
}

func (m *ModelFilterQueryBuilder) Build() (url.Values, error) {
	params := NewExtendedURLValues()
	params.AddPair("includeDeleted", m.ModelFilter.IncludeDeleted)
	params.AddPair("createdRange", m.ModelFilter.CreatedRange)
	params.AddPair("updatedRange", m.ModelFilter.UpdatedRange)
	return params.Values, nil
}

type QueryParamsFilterQueryBuilder struct {
	QueryParamsFilter filter.QueryParams
}

func (q *QueryParamsFilterQueryBuilder) Build() (url.Values, error) {
	params := NewExtendedURLValues()
	params.AddPair("page", q.QueryParamsFilter.Page)
	params.AddPair("size", q.QueryParamsFilter.PageSize)
	params.AddPair("sortBy", q.QueryParamsFilter.OrderByField)
	params.AddPair("sort", q.QueryParamsFilter.SortDirection)
	return params.Values, nil
}

type TransactionFilterQueryBuilder struct {
	TransactionFilter       filter.TransactionFilter
	ModelFilterQueryBuilder ModelFilterQueryBuilder
}

func (t *TransactionFilterQueryBuilder) Build() (url.Values, error) {
	mfv, err := t.ModelFilterQueryBuilder.Build()
	if err != nil {
		return nil, err
	}

	params := NewExtendedURLValues()
	if len(mfv) > 0 {
		params.Append(mfv)
	}

	params.AddPair("id", t.TransactionFilter.Id)
	params.AddPair("hex", t.TransactionFilter.Hex)
	params.AddPair("blockHash", t.TransactionFilter.BlockHash)
	params.AddPair("blockHeight", t.TransactionFilter.BlockHeight)
	params.AddPair("fee", t.TransactionFilter.Fee)
	params.AddPair("numberOfInputs", t.TransactionFilter.NumberOfInputs)
	params.AddPair("numberOfOutputs", t.TransactionFilter.NumberOfOutputs)
	params.AddPair("draftId", t.TransactionFilter.DraftID)
	params.AddPair("totalValue", t.TransactionFilter.TotalValue)
	params.AddPair("status", t.TransactionFilter.Status)
	return params.Values, nil
}

type ExtendedURLValues struct {
	url.Values
}

func (e *ExtendedURLValues) AddPair(key string, val any) {
	if val == nil || len(key) == 0 {
		return
	}
	write := func(v any) { e.Add(key, fmt.Sprintf("%v", v)) }
	writeRange := func(v filter.TimeRange) {
		if v.From != nil && !v.From.IsZero() {
			e.Add(fmt.Sprintf("%s[from]", key), v.From.Format(time.RFC3339))
		}
		if v.To != nil && !v.To.IsZero() {
			e.Add(fmt.Sprintf("%s[to]", key), v.To.Format(time.RFC3339))
		}
	}

	switch v := val.(type) {
	case int:
		if v > 0 {
			write(v)
		}
	case string:
		if len(v) > 0 {
			write(v)
		}
	case *string:
		if v != nil && len(*v) > 0 {
			write(*v)
		}
	case *uint64:
		if v != nil && *v > 0 {
			write(*v)
		}
	case *uint32:
		if v != nil && *v > 0 {
			write(*v)
		}
	case *bool:
		if v != nil {
			write(*v)
		}
	case *filter.TimeRange:
		if v != nil {
			writeRange(*v)
		}
	}
}

func (e *ExtendedURLValues) Append(vv ...url.Values) {
	for _, v := range vv {
		for k, iv := range v {
			e.Values[k] = append(e.Values[k], iv...)
		}
	}
}

func NewExtendedURLValues() *ExtendedURLValues {
	e := ExtendedURLValues{
		make(url.Values),
	}
	return &e
}
