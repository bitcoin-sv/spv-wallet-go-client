package cryptoutil_test

import (
	"math"
	"testing"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	compat "github.com/bitcoin-sv/go-sdk/compat/bip32"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/cryptoutil"
	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	tests := map[string]struct {
		expectedHash string
		expectedErr  error
		input        string
	}{
		"input: empty": {
			expectedHash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			input:        "",
		},
		"input: 1234567": {
			input:        "1234567",
			expectedHash: "8bb0cf6eb9b17d0f7d22b456f121257dc1254e1f01665370476383ea776df414",
		},
		"input: xpub": {
			input:        "xpub661MyMwAqRbcFrBJbKwBGCB7d3fr2SaAuXGM95BA62X41m6eW2ehRQGW4xLi9wkEXUGnQZYxVVj4PxXnyrLk7jdqvBAs1Qq9gf6ykMvjR7J",
			expectedHash: "1a0b10d4eda0636aae1709e7e7080485a4d99af3ca2962c6e677cf5b53d8ab8c",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Run(name, func(t *testing.T) {
				got := cryptoutil.Hash(tc.input)
				require.Equal(t, tc.expectedHash, got)
			})
		})
	}
}

func TestDeriveChildKeyFromHex(t *testing.T) {
	const (
		input         = "8bb0cf6eb9b17d0f7d22b456f121257dc1254e1f01665370476383ea776df414"
		XPriv         = "xprv9s21ZrQH143K3N6qVJQAu4EP51qMcyrKYJLkLgmYXgz58xmVxVLSsbx2DfJUtjcnXK8NdvkHMKfmmg5AJT2nqqRWUrjSHX29qEJwBgBPkJQ"
		XPub          = "xpub661MyMwAqRbcFrBJbKwBGCB7d3fr2SaAuXGM95BA62X41m6eW2ehRQGW4xLi9wkEXUGnQZYxVVj4PxXnyrLk7jdqvBAs1Qq9gf6ykMvjR7J"
		expectedXPriv = "xprvA8mj2ZL1w6Nqpi6D2amJLo4Gxy24tW9uv82nQKmamT2rkg5DgjzJZRFnW33e7QJwn65uUWSuN6YQyWrujNjZdVShPRnpNUSRVTru4cxaqfd"
		expectedXPub  = "xpub6Mm5S4rumTw93CAg8cJJhw11WzrZHxsmHLxPCiBCKnZqdUQNEHJZ7DaGMKucRzXPHtoS2ZqsVSRjxVbibEvwmR2wXkZDd8RrTftmm42cRsf"
	)

	generateHDKey := func(s string) *bip32.ExtendedKey {
		k, err := compat.GenerateHDKeyFromString(s)
		if err != nil {
			t.Fatal(err)
		}
		return k
	}

	t.Run("child extended key from  xpriv", func(t *testing.T) {
		key := generateHDKey(XPriv)
		got, err := cryptoutil.DeriveChildKeyFromHex(key, input)

		require.NoError(t, err)
		require.Equal(t, expectedXPriv, got.String())
	})

	t.Run("child extended key from xpub", func(t *testing.T) {
		key := generateHDKey(XPub)
		got, err := cryptoutil.DeriveChildKeyFromHex(key, input)

		require.NoError(t, err)
		require.Equal(t, expectedXPub, got.String())
	})

	t.Run("extended public key from extended private key", func(t *testing.T) {
		key := generateHDKey(XPriv)
		child, err := cryptoutil.DeriveChildKeyFromHex(key, input)
		require.NoError(t, err)

		got, err := child.Neuter()
		require.NoError(t, err)
		require.Equal(t, expectedXPub, got.String())
	})
}

func TestRandomHex(t *testing.T) {
	tests := map[string]struct {
		input       int
		expectedLen int
	}{
		"input: zero": {
			input:       0,
			expectedLen: 0,
		},
		"input: 100_000": {
			input:       100_000,
			expectedLen: 200_000,
		},
		"input: 16": {
			input:       16,
			expectedLen: 32,
		},
		"input: 32": {
			input:       32,
			expectedLen: 64,
		},

		"input: 8": {
			input:       8,
			expectedLen: 16,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := cryptoutil.RandomHex(tc.input)
			require.NoError(t, err)
			require.Len(t, got, tc.expectedLen)
		})
	}
}

func TestInt64ToUint32(t *testing.T) {
	tests := map[string]struct {
		input          int64
		expectedErr    error
		expectedUint32 uint32
	}{
		"input: negative value": {
			input:          -1,
			expectedErr:    errors.ErrNegativeValueNotAllowed,
			expectedUint32: 0,
		},
		"input: max value exceeded": {
			input:          math.MaxUint32 + 1,
			expectedErr:    errors.ErrMaxUint32LimitExceeded,
			expectedUint32: 0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := cryptoutil.Int64ToUint32(tc.input)
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedUint32, got)
		})
	}
}

func TestParseChildNumsFromHex(t *testing.T) {
	tests := map[string]struct {
		hex            string
		expectedErr    error
		expectedResult []uint32
	}{
		"input: empty hex": {
			hex:            "",
			expectedErr:    nil,
			expectedResult: nil,
		},
		"input: invalid hex": {
			hex:            "test",
			expectedErr:    errors.ErrHexHashPartIntParse,
			expectedResult: nil,
		},
		"input: short hex ababab": {
			hex:            "ababab",
			expectedErr:    nil,
			expectedResult: []uint32{11250603},
		},
		"input: medium hex 8bb0cf6eb9b17d0f7d22b456f121257dc1254e1f01665370476383ea776df414": {
			hex:         "8bb0cf6eb9b17d0f7d22b456f121257dc1254e1f01665370476383ea776df414",
			expectedErr: nil,
			expectedResult: []uint32{
				196136815,  // 8bb0cf6e = 2343620462 - 2147483647
				967933200,  // b9b17d0f = 3115416847 - 2147483647
				2099426390, // 7d22b456
				1897997694, // f121257d = 4045481341 - 2147483647
				1092963872, // c1254e1f = 3240447519 - 2147483647
				23483248,   // 01665370
				1197704170, // 476383ea
				2003694612, // 776df414
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := cryptoutil.ParseChildNumsFromHex(tc.hex)
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResult, got)
		})
	}
}
