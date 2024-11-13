package querybuilders_test

import (
	"net/url"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/stretchr/testify/require"
)

func TestMetadataFilterBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		metadata       querybuilders.Metadata
		expectedParams url.Values
		expectedErr    error
		depth          int
	}{
		"metadata: empty map": {
			depth:          querybuilders.DefaultMaxDepth,
			expectedParams: make(url.Values),
		},
		"metadata: map entry [key]=nil": {
			depth:          querybuilders.DefaultMaxDepth,
			expectedParams: make(url.Values),
			metadata: querybuilders.Metadata{
				"key": nil,
			},
		},
		"metadata: map entries [key1]=value1, [key2]=nil": {
			depth: querybuilders.DefaultMaxDepth,
			expectedParams: url.Values{
				"metadata[key1][]": []string{"value1"},
			},
			metadata: querybuilders.Metadata{
				"key1": []string{"value1"},
				"key2": nil,
			},
		},
		"metadata: map entry [key]=value1": {
			depth: querybuilders.DefaultMaxDepth,
			expectedParams: url.Values{
				"metadata[key]": []string{"value1"},
			},
			metadata: querybuilders.Metadata{
				"key": "value1",
			},
		},
		"metadata: map entries [key1]=value1, [key2]=1024": {
			depth: querybuilders.DefaultMaxDepth,
			expectedParams: url.Values{
				"metadata[key1]": []string{"value1"},
				"metadata[key2]": []string{"1024"},
			},
			metadata: querybuilders.Metadata{
				"key1": "value1",
				"key2": 1024,
			},
		},
		"metadata: two keys nested in one": {
			depth: querybuilders.DefaultMaxDepth,
			expectedParams: url.Values{
				"metadata[key1][key2]": []string{"value1"},
				"metadata[key1][key3]": []string{"1024"},
			},
			metadata: querybuilders.Metadata{
				"key1": querybuilders.Metadata{
					"key2": "value1",
					"key3": 1024,
				},
			},
		},
		"metadata: map entries [hey=123&522]=value1, [key2]=value=123": {
			depth: querybuilders.DefaultMaxDepth,
			expectedParams: url.Values{
				"metadata[hey=123&522]": []string{"value1"},
				"metadata[key2]":        []string{"value=123"},
			},
			metadata: querybuilders.Metadata{
				"hey=123&522": "value1",
				"key2":        "value=123",
			},
		},
		"metadata: map entries [key1]=value1, [key2]=[]{value2,value3,value4}": {
			depth: querybuilders.DefaultMaxDepth,
			expectedParams: url.Values{
				"metadata[key1]":   []string{"value1"},
				"metadata[key2][]": []string{"value2", "value3", "value4"},
			},
			metadata: querybuilders.Metadata{
				"key1": "value1",
				"key2": []string{"value2", "value3", "value4"},
			},
		},
		"metadata: map entries [key1]=value1, [key2]=[]{value2, value3, value4}, [key3]=value5, [key4]=[]{value6,value7,value8}": {
			depth: querybuilders.DefaultMaxDepth,
			expectedParams: url.Values{
				"metadata[key1]":   []string{"value1"},
				"metadata[key2][]": []string{"value2", "value3", "value4"},
				"metadata[key3]":   []string{"value5"},
				"metadata[key4][]": []string{"value6", "value7", "value8"},
			},
			metadata: querybuilders.Metadata{
				"key1": "value1",
				"key2": []string{"value2", "value3", "value4"},
				"key3": "value5",
				"key4": []string{"value6", "value7", "value8"},
			},
		},
		"metadata: map entries [key1]=value1, [key2]=[value1,value2,value3,value4], [key3][key3_nested]=value5, [key4][key4_nested]=[6, 7]": {
			depth: querybuilders.DefaultMaxDepth,
			expectedParams: url.Values{
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
			},
			metadata: querybuilders.Metadata{
				"key1": "value1",
				"key2": []string{"value2", "value3", "value4"},
				"key3": querybuilders.Metadata{
					"key3_nested": "value5",
				},
				"key4": querybuilders.Metadata{
					"key4_nested": []int{6, 7},
				},
			},
		},
		"metadata: 11 map entries, complex nesting, max depth set to 100": {
			depth: querybuilders.DefaultMaxDepth,
			expectedParams: url.Values{
				"metadata[key1][key2][key3][key1]":                                             []string{"abc"},
				"metadata[key1][key2][key3][key2][key1]":                                       []string{"9"},
				"metadata[key1][key2][key3][key3][key1][key2][key1][]":                         []string{"1", "2", "3", "4"},
				"metadata[key1][key2][key3][key3][key1][key2][key2]":                           []string{"10"},
				"metadata[key1][key2][key3][key3][key1][key2][key3]":                           []string{"abc"},
				"metadata[key1][key2][key3][key3][key1][key2][key4][key1][key1][key1]":         []string{"2"},
				"metadata[key1][key2][key3][key3][key1][key2][key4][key1][key1][key2]":         []string{"cde"},
				"metadata[key1][key2][key3][key3][key1][key2][key4][key1][key1][key3][key1][]": []string{"5", "6", "7", "8"},
				"metadata[key1][key2][key3][key3][key1][key2][key4][key1][key1][key3][key2][]": []string{"a", "b", "c"},
			},
			metadata: querybuilders.Metadata{
				"key1": querybuilders.Metadata{
					"key2": querybuilders.Metadata{
						"key3": querybuilders.Metadata{
							"key1": "abc",
							"key2": querybuilders.Metadata{
								"key1": 9,
							},
							"key3": querybuilders.Metadata{
								"key1": querybuilders.Metadata{
									"key2": querybuilders.Metadata{
										"key1": []int{1, 2, 3, 4},
										"key2": 10,
										"key3": "abc",
										"key4": querybuilders.Metadata{
											"key1": querybuilders.Metadata{
												"key1": querybuilders.Metadata{
													"key1": 2,
													"key2": "cde",
													"key3": querybuilders.Metadata{
														"key1": []int{5, 6, 7, 8},
														"key2": []string{"a", "b", "c"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"metadata: map entries depth exceeded - map entries: 4, max depth: 3": {
			metadata: querybuilders.Metadata{
				"key1": querybuilders.Metadata{
					"key2": querybuilders.Metadata{
						"key3": querybuilders.Metadata{
							"key4": "value1",
						},
					},
				},
			},
			depth:       3,
			expectedErr: querybuilders.ErrMetadataFilterMaxDepthExceeded,
		},
		"metadata: unsupported map in array": {
			metadata: querybuilders.Metadata{
				"key1": querybuilders.Metadata{
					"key2": []any{
						querybuilders.Metadata{
							"key3": "value1",
						},
					},
				},
			},
			depth:       querybuilders.DefaultMaxDepth,
			expectedErr: querybuilders.ErrMetadataWrongTypeInArray,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			builder := querybuilders.MetadataFilterBuilder{
				MaxDepth: tc.depth,
				Metadata: tc.metadata,
			}

			// then:
			got, err := builder.Build()
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}
