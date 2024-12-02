package querybuilderstest

import (
	"testing"
	"time"
)

func ParseTime(t *testing.T, s string) time.Time {
	ts, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t.Fatalf("test helper - time parse: %s", err)
	}
	return ts
}

func Ptr[T any](value T) *T {
	return &value
}
