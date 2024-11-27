package querybuilders_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders/querybuilderstest"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestExtendedURLValues_AddPair(t *testing.T) {
	// given:
	to := querybuilderstest.ParseTime(t, "2024-10-07T14:03:26.736816Z")
	from := querybuilderstest.ParseTime(t, "2024-10-07T14:03:26.736816Z")
	expectedValues := url.Values{
		"key1":       []string{"str"},
		"key2":       []string{"1"},
		"key3":       []string{"str_ptr"},
		"key4":       []string{"64"},
		"key5":       []string{"32"},
		"key6":       []string{"false"},
		"key7[from]": []string{from.Format(time.RFC3339)},
		"key7[to]":   []string{to.Format(time.RFC3339)},
	}

	// when:
	params := querybuilders.NewExtendedURLValues()
	params.AddPair("key1", "str")
	params.AddPair("key2", 1)
	params.AddPair("key3", querybuilderstest.Ptr("str_ptr"))
	params.AddPair("key4", querybuilderstest.Ptr(uint64(64)))
	params.AddPair("key5", querybuilderstest.Ptr(uint32(32)))
	params.AddPair("key6", querybuilderstest.Ptr(bool(false)))
	params.AddPair("key7", &filter.TimeRange{
		From: &from,
		To:   &to,
	})

	// then:
	require.EqualValues(t, expectedValues, params.Values)
}

func TestExtendedURLValues_ParseToMap(t *testing.T) {
	// given:
	to := querybuilderstest.ParseTime(t, "2024-10-07T14:03:26.736816Z")
	from := querybuilderstest.ParseTime(t, "2024-10-07T14:03:26.736816Z")
	expectedValues := map[string]string{
		"key1":       "str",
		"key2":       "1",
		"key3":       "str_ptr",
		"key4":       "64",
		"key5":       "32",
		"key6":       "false",
		"key7[from]": from.Format(time.RFC3339),
		"key7[to]":   to.Format(time.RFC3339),
	}

	params := querybuilders.NewExtendedURLValues()
	params.AddPair("key1", "str")
	params.AddPair("key2", 1)
	params.AddPair("key3", querybuilderstest.Ptr("str_ptr"))
	params.AddPair("key4", querybuilderstest.Ptr(uint64(64)))
	params.AddPair("key5", querybuilderstest.Ptr(uint32(32)))
	params.AddPair("key6", querybuilderstest.Ptr(bool(false)))
	params.AddPair("key7", &filter.TimeRange{
		From: &from,
		To:   &to,
	})

	// when:
	got := params.ParseToMap()

	// then:
	require.EqualValues(t, expectedValues, got)
}
