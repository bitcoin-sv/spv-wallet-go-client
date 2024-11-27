package querybuilders

import (
	"fmt"
	"net/url"
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

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

func (e *ExtendedURLValues) ParseToMap() map[string]string {
	m := make(map[string]string)
	for k, v := range e.Values {
		m[k] = v[0]
	}

	return m
}

func (e *ExtendedURLValues) Append(vv ...url.Values) {
	for _, v := range vv {
		for k, iv := range v {
			e.Values[k] = append(e.Values[k], iv...)
		}
	}
}

func NewExtendedURLValues() *ExtendedURLValues {
	return &ExtendedURLValues{
		make(url.Values),
	}
}
