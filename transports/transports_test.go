package transports

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWithSignRequest will test the method WithSignRequest()
func TestWithSignRequest(t *testing.T) {

	t.Run("get opts", func(t *testing.T) {
		opt := WithSignRequest(false)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("sign request false", func(t *testing.T) {
		opts := []ClientOps{
			WithSignRequest(false),
			WithHTTP(""),
		}
		c, err := NewTransport(opts...)
		require.NoError(t, err)
		require.NotNil(t, c)

		assert.Equal(t, false, c.IsSignRequest())
	})

	t.Run("sign request true", func(t *testing.T) {
		opts := []ClientOps{
			WithSignRequest(true),
			WithHTTP(""),
		}
		c, err := NewTransport(opts...)
		require.NoError(t, err)
		require.NotNil(t, c)

		assert.Equal(t, true, c.IsSignRequest())
	})
}
