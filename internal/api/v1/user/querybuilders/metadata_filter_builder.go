package querybuilders

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
)

type Metadata map[string]any

const DefaultMaxDepth = 100

type metadataPath string

func (m metadataPath) NestPath(key any) metadataPath {
	return metadataPath(fmt.Sprintf("%s[%v]", m, key))
}

func (m metadataPath) AddToURL(urlValues url.Values, value any) {
	urlValues.Add(string(m), fmt.Sprintf("%v", value))
}

func (m metadataPath) AddArrayToURL(urlValues url.Values, values []any) {
	key := string(m) + "[]"
	for _, value := range values {
		urlValues.Add(key, fmt.Sprintf("%v", value))
	}
}

func newMetadataPath(key string) metadataPath {
	return metadataPath(fmt.Sprintf("metadata[%s]", key))
}

type MetadataFilterBuilder struct {
	MaxDepth int
	Metadata Metadata
}

func (m *MetadataFilterBuilder) Build() (url.Values, error) {
	params := make(url.Values)
	for k, v := range m.Metadata {
		path := newMetadataPath(k)
		if err := m.generateQueryParams(0, path, v, params); err != nil {
			return nil, err
		}
	}

	return params, nil
}

func (m *MetadataFilterBuilder) generateQueryParams(depth int, path metadataPath, val any, params url.Values) error {
	if depth > m.MaxDepth {
		return fmt.Errorf("%w - max depth: %d", ErrMetadataFilterMaxDepthExceeded, m.MaxDepth)
	}

	if val == nil {
		return nil
	}

	switch reflect.TypeOf(val).Kind() {
	case reflect.Map:
		return m.processMapQueryParams(depth+1, val, path, params)
	case reflect.Slice:
		return m.processSliceQueryParams(val, path, params)
	default:
		path.AddToURL(params, val)
		return nil
	}
}

func (m *MetadataFilterBuilder) processMapQueryParams(depth int, val any, param metadataPath, params url.Values) error {
	rval := reflect.ValueOf(val)
	for _, k := range rval.MapKeys() {
		nested := param.NestPath(k.Interface())
		if err := m.generateQueryParams(depth+1, nested, rval.MapIndex(k).Interface(), params); err != nil {
			return err
		}
	}

	return nil
}

func (m *MetadataFilterBuilder) processSliceQueryParams(val any, path metadataPath, params url.Values) error {
	slice := reflect.ValueOf(val)
	arr := make([]any, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		item := slice.Index(i)

		// safe check - only primitive types are allowed in arrays
		// note: kind := item.Kind() is not enough, because it returns interface instead of actual underlying type
		kind := reflect.TypeOf(item.Interface()).Kind()
		if kind == reflect.Map || kind == reflect.Slice {
			return ErrMetadataWrongTypeInArray
		}

		arr[i] = item.Interface()
	}
	path.AddArrayToURL(params, arr)

	return nil
}

var (
	ErrMetadataFilterMaxDepthExceeded = errors.New("maximum depth of nesting in metadata map exceeded")
	ErrMetadataWrongTypeInArray       = errors.New("wrong type in array")
)
