package totp

import (
	"testing"
	"time"

	"github.com/libsv/go-bk/bip32"
	"github.com/stretchr/testify/require"
)

func TestTotpService(t *testing.T) {
	const a_xprivStr = "xprv9s21ZrQH143K4BNj7w6PX1kTzUDMrKGc9VUEHhguqdbdyPVdP6t6NJNNMxuksetwDCFiauEvwNtNyt5xkXPS6eDtvf7e1GAcDERdoeApfGc"
	const b_xprivStr = "xprv9s21ZrQH143K2yLE69dKjrRmHE5atymAuVzqWj7gZcBDdhSZQcGZiFmMBwpXDC36bvMV398HXwVK4UUCwB5oddU5dU93QKD7JFyyizYwLDp"

	t.Run("Passcode has the exact number of digits", func(t *testing.T) {
		// given
		const givenDigits = 4

		a_xpriv, _ := bip32.NewKeyFromString(a_xprivStr)
		b_xpriv, _ := bip32.NewKeyFromString(b_xprivStr)
		b_xpub, _ := b_xpriv.Neuter()

		sut := &Service{
			Digits: givenDigits,
		}

		// when
		pc, err := sut.GenarateTotp(a_xpriv, b_xpub)

		// then
		require.NoError(t, err)
		require.Len(t, pc, givenDigits)
	})

	t.Run("Passcode is valid", func(t *testing.T) {
		// given
		a_xpriv, _ := bip32.NewKeyFromString(a_xprivStr)
		a_xpub, _ := a_xpriv.Neuter()
		b_xpriv, _ := bip32.NewKeyFromString(b_xprivStr)
		b_xpub, _ := b_xpriv.Neuter()

		sut := &Service{
			Digits: 2,
		}
		a_passcode, err := sut.GenarateTotp(a_xpriv, b_xpub)
		require.NoError(t, err)

		// when
		isValid, err := sut.ValidateTotp(b_xpriv, a_xpub, a_passcode)
		require.NoError(t, err)

		// then
		require.NoError(t, err)
		require.True(t, isValid)
	})

	t.Run("Passcode is invalid after given seconds", func(t *testing.T) {
		// given
		const givenSeconds = 3

		a_xpriv, _ := bip32.NewKeyFromString(a_xprivStr)
		a_xpub, _ := a_xpriv.Neuter()
		b_xpriv, _ := bip32.NewKeyFromString(b_xprivStr)
		b_xpub, _ := b_xpriv.Neuter()

		sut := &Service{
			Digits: 2,
			Period: givenSeconds,
		}
		a_passcode, err := sut.GenarateTotp(a_xpriv, b_xpub)
		require.NoError(t, err)

		// when
		isValid, err := sut.ValidateTotp(b_xpriv, a_xpub, a_passcode)
		require.NoError(t, err)
		require.True(t, isValid)

		time.Sleep(givenSeconds * time.Second)

		isValid, err = sut.ValidateTotp(b_xpriv, a_xpub, a_passcode)

		// then
		require.NoError(t, err)
		require.False(t, isValid)
	})

}
